package api

import (
	"net/http"

	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/models"
	"github.com/supabase/auth/internal/storage"
)

// RecoverParams holds the parameters for a password recovery request
type RecoverParams struct {
	Email               string                 `json:"email"`
	Phone               string                 `json:"phone"`
	CodeChallenge       string                 `json:"code_challenge"`
	CodeChallengeMethod string                 `json:"code_challenge_method"`
	Data                map[string]interface{} `json:"data,omitempty"`
}

func (p *RecoverParams) Validate(a *API) error {
	if p.Email == "" && p.Phone == "" {
		return apierrors.NewBadRequestError(apierrors.ErrorCodeValidationFailed, "Password recovery requires an email or phone")
	}

	if p.Email != "" && p.Phone != "" {
		return apierrors.NewBadRequestError(apierrors.ErrorCodeValidationFailed, "Only one of email or phone should be specified")
	}

	var err error
	if p.Email != "" {
		if p.Email, err = a.validateEmail(p.Email); err != nil {
			return err
		}
	}

	if p.Phone != "" {
		if p.Phone, err = validatePhone(p.Phone); err != nil {
			return err
		}
	}

	if err := validatePKCEParams(p.CodeChallengeMethod, p.CodeChallenge); err != nil {
		return err
	}
	return nil
}

// Recover sends a recovery email
func (a *API) Recover(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	db := a.db.WithContext(ctx)
	params := &RecoverParams{}
	if err := retrieveRequestParams(r, params); err != nil {
		return err
	}

	flowType := getFlowFromChallenge(params.CodeChallenge)
	if err := params.Validate(a); err != nil {
		return err
	}

	var user *models.User
	var err error
	aud := a.requestAud(ctx, r)

	if params.Email != "" {
		user, err = models.FindUserByEmailAndAudience(db, params.Email, aud)
	} else if params.Phone != "" {
		user, err = models.FindUserByPhoneAndAudience(db, params.Phone, aud)
	}

	if err != nil {
		if models.IsNotFoundError(err) {
			return sendJSON(w, http.StatusOK, map[string]string{})
		}
		return apierrors.NewInternalServerError("Unable to process request").WithInternalError(err)
	}
	if isPKCEFlow(flowType) {
		if _, err := generateFlowState(db, models.Recovery.String(), models.Recovery, params.CodeChallengeMethod, params.CodeChallenge, &(user.ID)); err != nil {
			return err
		}
	}

	err = db.Transaction(func(tx *storage.Connection) error {
		if terr := models.NewAuditLogEntry(r, tx, user, models.UserRecoveryRequestedAction, "", nil); terr != nil {
			return terr
		}
		if params.Email != "" {
			return a.sendPasswordRecovery(r, tx, user, flowType)
		} else if params.Phone != "" {
			return a.sendPasswordRecoverySMS(r, tx, user, flowType)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return sendJSON(w, http.StatusOK, map[string]string{})
}
