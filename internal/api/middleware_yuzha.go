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

func (a *API) verifyYuZhaLabCode(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()

	// body := &security.YuZhaLabRequest{}
	// if err := security.RetrieveYuZhaLabRequestParams(req, body); err != nil {
	// 	return nil, err
	// }

	body := &SignupParams{}
	if err := retrieveRequestParams(req, body); err != nil {
		return nil, err
	}

	verificationResult, err := VerifyYuZhaLabRequest(body)
	if err != nil {
		return nil, apierrors.NewBadRequestError(apierrors.ErrorCodeCaptchaFailed, "yuzha lab code verification process failed: %s", err.Error())
	}

	if !verificationResult.Success {
		return nil, apierrors.NewBadRequestError(apierrors.ErrorCodeMFAVerificationFailed, "yuzha lab code verification failed: request disallowed (%s)", strings.Join(verificationResult.ErrorCodes, ", "))
	}

	return ctx, nil
}

func VerifyYuZhaLabRequest(requestBody *SignupParams) (security.VerificationResponse, error) {
	if requestBody.Data == nil {
		return security.VerificationResponse{}, errors.New("request data is nil")
	}

	codeValue, exists := requestBody.Data["verify_code"]
	if !exists {
		return security.VerificationResponse{}, errors.New("code not found in request data")
	}

	codeResponse, ok := codeValue.(string)
	if !ok {
		return security.VerificationResponse{}, errors.New("code is not a string")
	}

	codeResponse = strings.TrimSpace(codeResponse)
	if codeResponse == "" {
		return security.VerificationResponse{}, errors.New("no code found in request")
	}
	email := strings.TrimSpace(requestBody.Email)
	phone := strings.TrimSpace(requestBody.Phone)

	if codeResponse == "" {
		return security.VerificationResponse{}, errors.New("no captcha response (captcha_token) found in request")
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

// verifyAndSetLanguage 检测并设置用户语言偏好到上下文 (类似verifyYuZhaLabCode的中间件模式)
func (a *API) verifyAndSetLanguage(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()

	// 使用i18n包检测用户语言偏好
	lang := i18n.GetLanguageFromRequest(req)

	// 将语言偏好存储到上下文中
	ctx = context.WithValue(ctx, i18n.UserLanguageKey, lang)

	return ctx, nil
}

/*
使用示例:

在路由中使用此中间件:
```go
// 在需要语言检测的路由中使用
func (a *API) setupRoutes() {
	// 使用verifyAndSetLanguage中间件
	r.With(a.requireAdminCredentials, a.verifyAndSetLanguage).Post("/admin/users", a.adminCreateUser)

	// 或者与其他中间件组合使用
	r.With(a.verifyYuZhaLabCode, a.verifyAndSetLanguage).Post("/signup", a.signup)
}

// 在处理函数中获取语言
func (a *API) someHandler(w http.ResponseWriter, r *http.Request) error {
	// 从上下文获取用户语言偏好
	lang := i18n.GetLanguageFromContext(r.Context())

	// 使用语言进行本地化处理
	message := i18n.GetMessage(lang, "success")

	return sendJSON(w, http.StatusOK, map[string]interface{}{
		"message": message,
	})
}
*/
