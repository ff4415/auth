package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/pkg/errors"
	"github.com/supabase/auth/internal/api/apierrors"
	"github.com/supabase/auth/internal/i18n"
	"github.com/supabase/auth/internal/security"
	"github.com/supabase/auth/internal/utilities"
)

// verifyAndSetLanguage 检测并设置用户语言偏好到上下文
// 适配 auth API 的中间件格式
func (a *API) verifyAndSetLanguage(_ http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()

	// 使用i18n包检测用户语言偏好
	lang := i18n.GetLanguageFromRequest(req)

	// 将语言偏好存储到上下文中
	ctx = context.WithValue(ctx, i18n.UserLanguageKey, lang)

	return ctx, nil
}

func (a *API) verifyYuZhaLabCode(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()

	// 根据请求路径判断应该使用哪种参数类型
	var requestBody interface{}
	var err error

	if strings.Contains(req.URL.Path, "/recover") {
		body := &RecoverParams{}
		if err := retrieveRequestParams(req, body); err != nil {
			return nil, err
		}
		requestBody = body
	} else {
		body := &SignupParams{}
		if err := retrieveRequestParams(req, body); err != nil {
			return nil, err
		}
		requestBody = body
	}

	verificationResult, err := VerifyYuZhaLabRequest(requestBody)
	if err != nil {
		return nil, apierrors.NewBadRequestError(apierrors.ErrorCodeCaptchaFailed, "yuzha lab code verification process failed: %s", err.Error())
	}

	if !verificationResult.Success {
		return nil, apierrors.NewBadRequestError(apierrors.ErrorCodeMFAVerificationFailed, "yuzha lab code verification failed: request disallowed (%s)", strings.Join(verificationResult.ErrorCodes, ", "))
	}

	return ctx, nil
}

func VerifyYuZhaLabRequest(requestBody interface{}) (security.VerificationResponse, error) {
	var email, phone, codeResponse string

	switch body := requestBody.(type) {
	case *SignupParams:
		if body.Data == nil {
			return security.VerificationResponse{}, errors.New("request data is nil")
		}

		codeValue, exists := body.Data["verify_code"]
		if !exists {
			return security.VerificationResponse{}, errors.New("code not found in request data")
		}

		code, ok := codeValue.(string)
		if !ok {
			return security.VerificationResponse{}, errors.New("code is not a string")
		}

		codeResponse = strings.TrimSpace(code)
		email = strings.TrimSpace(body.Email)
		phone = strings.TrimSpace(body.Phone)

	case *RecoverParams:
		if body.Data == nil {
			return security.VerificationResponse{}, errors.New("request data is nil")
		}

		codeValue, exists := body.Data["verify_code"]
		if !exists {
			return security.VerificationResponse{}, errors.New("code not found in request data")
		}

		code, ok := codeValue.(string)
		if !ok {
			return security.VerificationResponse{}, errors.New("code is not a string")
		}

		codeResponse = strings.TrimSpace(code)
		email = strings.TrimSpace(body.Email)
		// RecoverParams 不支持 phone

	default:
		return security.VerificationResponse{}, errors.New("unsupported request body type")
	}

	if codeResponse == "" {
		return security.VerificationResponse{}, errors.New("no code found in request")
	}

	if email == "" && phone == "" {
		return security.VerificationResponse{}, apierrors.NewBadRequestError(apierrors.ErrorCodeValidationFailed, "email or phone is required")
	}

	if email != "" && phone != "" {
		return security.VerificationResponse{}, apierrors.NewBadRequestError(apierrors.ErrorCodeValidationFailed, "email and phone cannot both be provided")
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
		return security.VerificationResponse{}, errors.Wrap(err, "failed to marshal JSON")
	}

	// 创建POST请求
	url := "http://127.0.0.1:8889/v1/verify/verifycode"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return security.VerificationResponse{}, errors.Wrap(err, "couldn't initialize request object for captcha check")
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	res, err := security.Client.Do(req)
	if err != nil {
		return security.VerificationResponse{}, errors.Wrap(err, "failed to verify captcha response")
	}
	defer utilities.SafeClose(res.Body)

	var verificationResponse security.VerificationResponse

	if err := json.NewDecoder(res.Body).Decode(&verificationResponse); err != nil {
		return security.VerificationResponse{}, errors.Wrap(err, "failed to decode captcha response: not JSON")
	}

	return verificationResponse, nil
}
