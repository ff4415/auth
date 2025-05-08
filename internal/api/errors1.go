package api

import (
	"context"
	"net/http"

	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/observability"
	"github.com/supabase/auth/internal/utilities"
)

func HandleResponseError_YuZhaGai(err error, w http.ResponseWriter, r *http.Request) {
	log := observability.GetLogEntry(r).Entry
	errorID := utilities.GetRequestID(r.Context())
	apiVersion, averr := DetermineClosestAPIVersion(r.Header.Get(APIVersionHeaderName))
	if averr != nil {
		log.WithError(averr).Warn("Invalid version passed to " + APIVersionHeaderName + " header, defaulting to initial version")
	} else if apiVersion != APIVersionInitial {
		// Echo back the determined API version from the request
		w.Header().Set(APIVersionHeaderName, FormatAPIVersion(apiVersion))
	}

	switch e := err.(type) {
	case *WeakPasswordError:
		if apiVersion.Compare(APIVersion20240101) >= 0 {
			var output struct {
				HTTPErrorResponse20240101
				// Payload struct {
				// 	Reasons []string `json:"reasons,omitempty"`
				// } `json:"weak_password,omitempty"`
			}

			output.Code = apierrors.ErrorCodeWeakPassword
			output.Message = e.Message
			// output.Payload.Reasons = e.Reasons

			if jsonErr := sendJSON(w, http.StatusUnprocessableEntity, output); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}

		} else {
			var output struct {
				HTTPError
				// Payload struct {
				// 	Reasons []string `json:"reasons,omitempty"`
				// } `json:"weak_password,omitempty"`
			}

			output.HTTPStatus = http.StatusUnprocessableEntity
			output.ErrorCode = apierrors.ErrorCodeWeakPassword
			output.Message = e.Message
			// output.Payload.Reasons = e.Reasons

			w.Header().Set("x-sb-error-code", output.ErrorCode)

			if jsonErr := sendJSON(w, output.HTTPStatus, output); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		}

	case *HTTPError:
		switch {
		case e.HTTPStatus >= http.StatusInternalServerError:
			e.ErrorID = errorID
			// this will get us the stack trace too
			log.WithError(e.Cause()).Error(e.Error())
		case e.HTTPStatus == http.StatusTooManyRequests:
			log.WithError(e.Cause()).Warn(e.Error())
		default:
			log.WithError(e.Cause()).Info(e.Error())
		}

		if e.ErrorCode != "" {
			w.Header().Set("x-sb-error-code", e.ErrorCode)
		}

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			resp := HTTPErrorResponse20240101{
				Code:    e.ErrorCode,
				Message: e.Message,
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

			// Provide better error messages for certain user-triggered Postgres errors.
			if pgErr := utilities.NewPostgresError(e.InternalError); pgErr != nil {
				if jsonErr := sendJSON(w, pgErr.HttpStatusCode, pgErr); jsonErr != nil && jsonErr != context.DeadlineExceeded {
					log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
				}
				return
			}

			if jsonErr := sendJSON(w, e.HTTPStatus, e); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		}

	case *OAuthError:
		log.WithError(e.Cause()).Info(e.Error())
		if jsonErr := sendJSON(w, http.StatusBadRequest, e); jsonErr != nil && jsonErr != context.DeadlineExceeded {
			log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
		}

	case ErrorCause:
		HandleResponseError(e.Cause(), w, r)

	default:
		log.WithError(e).Errorf("Unhandled server error: %s", e.Error())

		if apiVersion.Compare(APIVersion20240101) >= 0 {
			resp := HTTPErrorResponse20240101{
				Code:    apierrors.ErrorCodeUnexpectedFailure,
				Message: "Unexpected failure, please check server logs for more information",
			}

			if jsonErr := sendJSON(w, http.StatusInternalServerError, resp); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		} else {
			httpError := HTTPError{
				HTTPStatus: http.StatusInternalServerError,
				ErrorCode:  apierrors.ErrorCodeUnexpectedFailure,
				Message:    "Unexpected failure, please check server logs for more information",
			}

			if jsonErr := sendJSON(w, http.StatusInternalServerError, httpError); jsonErr != nil && jsonErr != context.DeadlineExceeded {
				log.WithError(jsonErr).Warn("Failed to send JSON on ResponseWriter")
			}
		}
	}
}
