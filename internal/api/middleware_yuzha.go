package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/security"
)

func (a *API) verifyYuZhaLabCode(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()

	body := &security.YuZhaLabRequest{}
	if err := security.RetrieveYuZhaLabRequestParams(req, body); err != nil {
		return nil, err
	}

	verificationResult, err := security.VerifyYuZhaLabRequest(body)
	if err != nil {
		return nil, apierrors.NewInternalServerError("yuzha lab code verification process failed").WithInternalError(err)
	}

	if !verificationResult.Success {
		return nil, apierrors.NewBadRequestError(apierrors.ErrorCodeMFAVerificationFailed, "yuzha lab code verification failed: request disallowed (%s)", strings.Join(verificationResult.ErrorCodes, ", "))
	}

	return ctx, nil
}
