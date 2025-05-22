package security

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/utilities"
)

type YuZhaLabRequest struct {
	Security YuZhaLabSecurity `json:"yuzhalab_meta_security"`
}

type YuZhaLabSecurity struct {
	Email string `json:"email"`
	Phone string `json:"phone"`
	Code  string `json:"code"`
}

type YuZhaLabRequestParams interface {
	YuZhaLabRequest |
		struct {
			Email string `json:"email"`
			Phone string `json:"phone"`
		} |
		struct {
			Email string `json:"email"`
		}
}

func RetrieveYuZhaLabRequestParams[A YuZhaLabRequestParams](r *http.Request, params *A) error {
	body, err := utilities.GetBodyBytes(r)
	if err != nil {
		return apierrors.NewInternalServerError("Could not read body into byte slice").WithInternalError(err)
	}
	if err := json.Unmarshal(body, params); err != nil {
		return apierrors.NewBadRequestError(apierrors.ErrorCodeBadJSON, "Could not parse request body as JSON: %v", err)
	}
	return nil
}

func VerifyYuZhaLabRequest(requestBody *YuZhaLabRequest) (VerificationResponse, error) {
	codeResponse := strings.TrimSpace(requestBody.Security.Code)
	email := strings.TrimSpace(requestBody.Security.Email)
	phone := strings.TrimSpace(requestBody.Security.Phone)

	if codeResponse == "" {
		return VerificationResponse{}, errors.New("no captcha response (captcha_token) found in request")
	}
	if email == "" && phone == "" {
		return VerificationResponse{}, apierrors.NewBadRequestError(apierrors.ErrorCodeValidationFailed, "email or phone is required")
	}

	if email != "" && phone != "" {
		return VerificationResponse{}, apierrors.NewBadRequestError(apierrors.ErrorCodeValidationFailed, "email and phone cannot both be provided")
	}

	b := make(map[string]string)

	if email != "" {
		b["email"] = email
	} else if phone != "" {
		b["phone"] = phone
	}
	b["code"] = codeResponse

	// 将请求体转换为JSON
	jsonData, err := json.Marshal(b)
	if err != nil {
		return VerificationResponse{}, errors.Wrap(err, "failed to marshal JSON")
	}

	// 创建POST请求
	url := "http://127.0.0.1:8889/v1/verify/verifycode"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return VerificationResponse{}, errors.Wrap(err, "couldn't initialize request object for captcha check")
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := Client.Do(req)
	if err != nil {
		return VerificationResponse{}, errors.Wrap(err, "failed to verify captcha response")
	}
	defer utilities.SafeClose(res.Body)

	var verificationResponse VerificationResponse

	if err := json.NewDecoder(res.Body).Decode(&verificationResponse); err != nil {
		return VerificationResponse{}, errors.Wrap(err, "failed to decode captcha response: not JSON")
	}

	return verificationResponse, nil
}
