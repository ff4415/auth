package api

import (
	"context"
	"net/http"

	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/i18n"
	"github.com/supabase/auth/internal/observability"
	"github.com/supabase/auth/internal/utilities"
)

func HandleResponseError_YuZhaGai(err error, w http.ResponseWriter, r *http.Request) {
	log := observability.GetLogEntry(r).Entry
	errorID := utilities.GetRequestID(r.Context())

	// Get user's language preference from context (set by middleware)
	userLang := i18n.GetLanguageFromContextHTTP(r)

	apiVersion, averr := DetermineClosestAPIVersion(r.Header.Get(APIVersionHeaderName))
	if averr != nil {
		log.WithError(averr).Warn("Invalid version passed to " + APIVersionHeaderName + " header, defaulting to initial version")
	} else if apiVersion != APIVersionInitial {
		// Echo back the determined API version from the request
		w.Header().Set(APIVersionHeaderName, FormatAPIVersion(apiVersion))
	}

	switch e := err.(type) {
	case *WeakPasswordError:
		// Log the original error for debugging
		log.Info("Weak password error: ", e.Error())

		// Get localized message
		localizedMessage := i18n.GetMessage(userLang, "weak_password")

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			var output struct {
				HTTPErrorResponse20240101
			}

			output.Code = apierrors.ErrorCodeWeakPassword
			output.Message = localizedMessage

			if jsonErr := sendJSON(w, http.StatusUnprocessableEntity, output); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}

		} else {
			var output struct {
				HTTPError
			}

			output.HTTPStatus = http.StatusUnprocessableEntity
			output.ErrorCode = apierrors.ErrorCodeWeakPassword
			output.Message = localizedMessage

			w.Header().Set("x-sb-error-code", output.ErrorCode)

			if jsonErr := sendJSON(w, output.HTTPStatus, output); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		}

	case *HTTPError:
		// Log the original error with full details for debugging
		switch {
		case e.HTTPStatus >= http.StatusInternalServerError:
			e.ErrorID = errorID
			// Log full error details for internal server errors
			log.WithError(e.Cause()).Error(e.Error())
		case e.HTTPStatus == http.StatusTooManyRequests:
			log.WithError(e.Cause()).Warn(e.Error())
		default:
			log.WithError(e.Cause()).Info(e.Error())
		}

		if e.ErrorCode != "" {
			w.Header().Set("x-sb-error-code", e.ErrorCode)
		}

		// Get user-friendly localized message
		localizedMessage := i18n.GetUserFriendlyMessage(userLang, string(e.ErrorCode), e.Message)

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			resp := HTTPErrorResponse20240101{
				Code:    e.ErrorCode,
				Message: localizedMessage,
			}

			if resp.Code == "" {
				if e.HTTPStatus == http.StatusInternalServerError {
					resp.Code = apierrors.ErrorCodeUnexpectedFailure
				} else {
					resp.Code = apierrors.ErrorCodeUnknown
				}
			}

			if jsonErr := sendJSON(w, e.HTTPStatus, resp); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		} else {
			if e.ErrorCode == "" {
				if e.HTTPStatus == http.StatusInternalServerError {
					e.ErrorCode = apierrors.ErrorCodeUnexpectedFailure
				} else {
					e.ErrorCode = apierrors.ErrorCodeUnknown
				}
			}

			// For PostgreSQL errors, still log them but return user-friendly messages
			if pgErr := utilities.NewPostgresError(e.InternalError); pgErr != nil {
				// Log the actual PostgreSQL error for debugging
				log.WithError(e.InternalError).Error("PostgreSQL error occurred")

				// Return a generic error message to the user
				httpError := HTTPError{
					HTTPStatus: http.StatusInternalServerError,
					ErrorCode:  apierrors.ErrorCodeUnexpectedFailure,
					Message:    i18n.GetMessage(userLang, "internal_server_error"),
				}

				if jsonErr := sendJSON(w, httpError.HTTPStatus, httpError); jsonErr != nil && jsonErr != context.DeadlineExceeded {
					log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
				}
				return
			}

			// Create a sanitized error response
			sanitizedError := HTTPError{
				HTTPStatus: e.HTTPStatus,
				ErrorCode:  e.ErrorCode,
				Message:    localizedMessage,
			}

			if jsonErr := sendJSON(w, sanitizedError.HTTPStatus, sanitizedError); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		}

	case *OAuthError:
		// Log the OAuth error for debugging
		log.WithError(e.Cause()).Info("OAuth error occurred")

		// Return user-friendly OAuth error message
		localizedMessage := i18n.GetUserFriendlyMessage(userLang, "", e.Description)

		sanitizedError := OAuthError{
			Err:         e.Err,
			Description: localizedMessage,
		}

		if jsonErr := sendJSON(w, http.StatusBadRequest, sanitizedError); jsonErr != nil && jsonErr != context.DeadlineExceeded {
			log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
		}

	case ErrorCause:
		HandleResponseError_YuZhaGai(e.Cause(), w, r)

	default:
		// Log the full error details for debugging
		log.WithError(e).Errorf("Unhandled server error: %s", e.Error())

		// Return generic error message to user
		localizedMessage := i18n.GetMessage(userLang, "unexpected_failure")

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			resp := HTTPErrorResponse20240101{
				Code:    apierrors.ErrorCodeUnexpectedFailure,
				Message: localizedMessage,
			}

			if jsonErr := sendJSON(w, http.StatusInternalServerError, resp); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		} else {
			httpError := HTTPError{
				HTTPStatus: http.StatusInternalServerError,
				ErrorCode:  apierrors.ErrorCodeUnexpectedFailure,
				Message:    localizedMessage,
			}

			if jsonErr := sendJSON(w, http.StatusInternalServerError, httpError); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		}
	}
}
