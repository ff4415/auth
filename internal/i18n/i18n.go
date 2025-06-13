package i18n

import (
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// Language represents supported languages
type Language string

const (
	LanguageEnglish Language = "en"
	LanguageChinese Language = "zh"
)

// Messages contains localized error messages
var Messages = map[Language]map[string]string{
	LanguageEnglish: {
		// Basic errors
		"unknown":               "Unknown error occurred",
		"unexpected_failure":    "Unexpected failure, please check server logs for more information",
		"validation_failed":     "Validation failed",
		"bad_json":              "Could not parse request body as JSON",
		"internal_server_error": "Internal server error",
		"unauthorized":          "Unauthorized",
		"forbidden":             "Forbidden",
		"not_found":             "Not found",
		"too_many_requests":     "Too many requests",
		"conflict":              "Conflict",
		"unprocessable_entity":  "Unprocessable entity",
		"bad_request":           "Bad request",

		// User related errors
		"email_exists":                 "A user with this email address has already been registered",
		"phone_exists":                 "A user with this phone number has already been registered",
		"user_not_found":               "User not found",
		"user_banned":                  "User has been banned",
		"user_already_exists":          "User already exists",
		"user_sso_managed":             "User is managed by SSO",
		"email_not_confirmed":          "Email not confirmed",
		"phone_not_confirmed":          "Phone not confirmed",
		"email_address_not_authorized": "Email address not authorized",
		"email_address_invalid":        "Invalid email address",
		"weak_password":                "Password does not meet security requirements",
		"same_password":                "New password must be different from the current password",

		// Session related errors
		"session_not_found":          "Session not found",
		"session_expired":            "Session has expired",
		"refresh_token_not_found":    "Refresh token not found",
		"refresh_token_already_used": "Refresh token has already been used",
		"flow_state_not_found":       "Flow state not found",
		"flow_state_expired":         "Flow state has expired",

		// Authentication related errors
		"bad_jwt":                    "Invalid JWT token",
		"not_admin":                  "Not authorized as admin",
		"no_authorization":           "No authorization provided",
		"invalid_credentials":        "Invalid login credentials",
		"reauthentication_needed":    "Reauthentication required",
		"reauthentication_not_valid": "Reauthentication failed",
		"insufficient_aal":           "Insufficient authentication assurance level",

		// OAuth related errors
		"bad_oauth_state":                       "Invalid OAuth state",
		"bad_oauth_callback":                    "Invalid OAuth callback",
		"oauth_provider_not_supported":          "OAuth provider not supported",
		"unexpected_audience":                   "Unexpected audience in token",
		"single_identity_not_deletable":         "Cannot delete the only identity",
		"email_conflict_identity_not_deletable": "Cannot delete identity with email conflict",
		"identity_already_exists":               "Identity already exists",
		"identity_not_found":                    "Identity not found",

		// Provider related errors
		"email_provider_disabled":     "Email provider is disabled",
		"phone_provider_disabled":     "Phone provider is disabled",
		"provider_disabled":           "Provider is disabled",
		"web3_provider_disabled":      "Web3 provider is disabled",
		"web3_unsupported_chain":      "Unsupported blockchain network",
		"anonymous_provider_disabled": "Anonymous provider is disabled",

		// MFA related errors
		"too_many_enrolled_mfa_factors":   "Too many enrolled MFA factors",
		"mfa_factor_name_conflict":        "MFA factor name already exists",
		"mfa_factor_not_found":            "MFA factor not found",
		"mfa_ip_address_mismatch":         "MFA IP address mismatch",
		"mfa_challenge_expired":           "MFA challenge has expired",
		"mfa_verification_failed":         "MFA verification failed",
		"mfa_verification_rejected":       "MFA verification rejected",
		"mfa_phone_enroll_not_enabled":    "MFA phone enrollment is not enabled",
		"mfa_phone_verify_not_enabled":    "MFA phone verification is not enabled",
		"mfa_totp_enroll_not_enabled":     "MFA TOTP enrollment is not enabled",
		"mfa_totp_verify_not_enabled":     "MFA TOTP verification is not enabled",
		"mfa_webauthn_enroll_not_enabled": "MFA WebAuthn enrollment is not enabled",
		"mfa_webauthn_verify_not_enabled": "MFA WebAuthn verification is not enabled",
		"mfa_verified_factor_exists":      "MFA factor already verified",

		// Verification related errors
		"captcha_failed":    "Captcha verification failed",
		"otp_expired":       "One-time password has expired",
		"otp_disabled":      "One-time password is disabled",
		"bad_code_verifier": "Invalid code verifier",

		// SAML related errors
		"saml_provider_disabled":     "SAML provider is disabled",
		"saml_relay_state_not_found": "SAML relay state not found",
		"saml_relay_state_expired":   "SAML relay state has expired",
		"saml_idp_not_found":         "SAML Identity Provider not found",
		"saml_assertion_no_user_id":  "SAML assertion missing user ID",
		"saml_assertion_no_email":    "SAML assertion missing email",
		"sso_provider_not_found":     "SSO provider not found",
		"saml_metadata_fetch_failed": "Failed to fetch SAML metadata",
		"saml_idp_already_exists":    "SAML Identity Provider already exists",
		"sso_domain_already_exists":  "SSO domain already exists",
		"saml_entity_id_mismatch":    "SAML entity ID mismatch",

		// Rate limit related errors
		"over_request_rate_limit":    "Too many requests, please try again later",
		"over_email_send_rate_limit": "Too many email sending attempts, please try again later",
		"over_sms_send_rate_limit":   "Too many SMS sending attempts, please try again later",

		// Webhook related errors
		"hook_timeout":                      "Webhook request timed out",
		"hook_timeout_after_retry":          "Webhook request timed out after retry",
		"hook_payload_over_size_limit":      "Webhook payload exceeds size limit",
		"hook_payload_invalid_content_type": "Invalid webhook payload content type",

		// Other errors
		"request_timeout":         "Request timed out",
		"manual_linking_disabled": "Manual linking is disabled",
		"sms_send_failed":         "Failed to send SMS",
		"invite_not_found":        "Invite not found",
		"signup_disabled":         "Sign up is disabled",
	},
	LanguageChinese: {
		// Basic errors
		"unknown":               "发生未知错误",
		"unexpected_failure":    "意外错误，请检查服务器日志以获取更多信息",
		"validation_failed":     "验证失败",
		"bad_json":              "无法解析请求体为JSON格式",
		"internal_server_error": "内部服务器错误",
		"unauthorized":          "未授权",
		"forbidden":             "禁止访问",
		"not_found":             "未找到",
		"too_many_requests":     "请求过于频繁",
		"conflict":              "冲突",
		"unprocessable_entity":  "无法处理的实体",
		"bad_request":           "错误的请求",

		// User related errors
		"email_exists":                 "该邮箱地址已被注册",
		"phone_exists":                 "该手机号码已被注册",
		"user_not_found":               "用户不存在",
		"user_banned":                  "用户已被禁用",
		"user_already_exists":          "用户已存在",
		"user_sso_managed":             "用户由SSO管理",
		"email_not_confirmed":          "邮箱未确认",
		"phone_not_confirmed":          "手机号未确认",
		"email_address_not_authorized": "邮箱地址未授权",
		"email_address_invalid":        "无效的邮箱地址",
		"weak_password":                "密码不符合安全要求",
		"same_password":                "新密码必须与当前密码不同",

		// Session related errors
		"session_not_found":          "会话不存在",
		"session_expired":            "会话已过期",
		"refresh_token_not_found":    "刷新令牌不存在",
		"refresh_token_already_used": "刷新令牌已被使用",
		"flow_state_not_found":       "流程状态不存在",
		"flow_state_expired":         "流程状态已过期",

		// Authentication related errors
		"bad_jwt":                    "无效的JWT令牌",
		"not_admin":                  "未授权为管理员",
		"no_authorization":           "未提供授权",
		"invalid_credentials":        "无效的登录凭据",
		"reauthentication_needed":    "需要重新认证",
		"reauthentication_not_valid": "重新认证失败",
		"insufficient_aal":           "认证保证级别不足",

		// OAuth related errors
		"bad_oauth_state":                       "无效的OAuth状态",
		"bad_oauth_callback":                    "无效的OAuth回调",
		"oauth_provider_not_supported":          "不支持的OAuth提供商",
		"unexpected_audience":                   "令牌中的受众不符合预期",
		"single_identity_not_deletable":         "无法删除唯一的身份",
		"email_conflict_identity_not_deletable": "无法删除存在邮箱冲突的身份",
		"identity_already_exists":               "身份已存在",
		"identity_not_found":                    "身份不存在",

		// Provider related errors
		"email_provider_disabled":     "邮箱提供商已禁用",
		"phone_provider_disabled":     "手机号提供商已禁用",
		"provider_disabled":           "提供商已禁用",
		"web3_provider_disabled":      "Web3提供商已禁用",
		"web3_unsupported_chain":      "不支持的区块链网络",
		"anonymous_provider_disabled": "匿名提供商已禁用",

		// MFA related errors
		"too_many_enrolled_mfa_factors":   "已注册过多MFA因素",
		"mfa_factor_name_conflict":        "MFA因素名称已存在",
		"mfa_factor_not_found":            "MFA因素不存在",
		"mfa_ip_address_mismatch":         "MFA IP地址不匹配",
		"mfa_challenge_expired":           "MFA挑战已过期",
		"mfa_verification_failed":         "MFA验证失败",
		"mfa_verification_rejected":       "MFA验证被拒绝",
		"mfa_phone_enroll_not_enabled":    "MFA手机注册未启用",
		"mfa_phone_verify_not_enabled":    "MFA手机验证未启用",
		"mfa_totp_enroll_not_enabled":     "MFA TOTP注册未启用",
		"mfa_totp_verify_not_enabled":     "MFA TOTP验证未启用",
		"mfa_webauthn_enroll_not_enabled": "MFA WebAuthn注册未启用",
		"mfa_webauthn_verify_not_enabled": "MFA WebAuthn验证未启用",
		"mfa_verified_factor_exists":      "MFA因素已验证",

		// Verification related errors
		"captcha_failed":    "验证码验证失败",
		"otp_expired":       "一次性密码已过期",
		"otp_disabled":      "一次性密码已禁用",
		"bad_code_verifier": "无效的代码验证器",

		// SAML related errors
		"saml_provider_disabled":     "SAML提供商已禁用",
		"saml_relay_state_not_found": "SAML中继状态不存在",
		"saml_relay_state_expired":   "SAML中继状态已过期",
		"saml_idp_not_found":         "SAML身份提供商不存在",
		"saml_assertion_no_user_id":  "SAML断言缺少用户ID",
		"saml_assertion_no_email":    "SAML断言缺少邮箱",
		"sso_provider_not_found":     "SSO提供商不存在",
		"saml_metadata_fetch_failed": "获取SAML元数据失败",
		"saml_idp_already_exists":    "SAML身份提供商已存在",
		"sso_domain_already_exists":  "SSO域名已存在",
		"saml_entity_id_mismatch":    "SAML实体ID不匹配",

		// Rate limit related errors
		"over_request_rate_limit":    "请求过于频繁，请稍后再试",
		"over_email_send_rate_limit": "邮件发送过于频繁，请稍后再试",
		"over_sms_send_rate_limit":   "短信发送过于频繁，请稍后再试",

		// Webhook related errors
		"hook_timeout":                      "Webhook请求超时",
		"hook_timeout_after_retry":          "Webhook请求重试后超时",
		"hook_payload_over_size_limit":      "Webhook负载超过大小限制",
		"hook_payload_invalid_content_type": "无效的Webhook负载内容类型",

		// Other errors
		"request_timeout":         "请求超时",
		"manual_linking_disabled": "手动链接已禁用",
		"sms_send_failed":         "发送短信失败",
		"invite_not_found":        "邀请不存在",
		"signup_disabled":         "注册已禁用",
	},
}

// GetLanguageFromRequest extracts language preference from request
func GetLanguageFromRequest(r *http.Request) Language {
	// Check query parameter first
	if lang := r.URL.Query().Get("lang"); lang != "" {
		return normalizeLanguage(lang)
	}

	// Check custom header
	if lang := r.Header.Get("X-Language"); lang != "" {
		return normalizeLanguage(lang)
	}

	// Check Accept-Language header
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		return parseAcceptLanguage(acceptLang)
	}

	// Default to English
	return LanguageEnglish
}

// normalizeLanguage converts language codes to supported languages
func normalizeLanguage(lang string) Language {
	lang = strings.ToLower(strings.TrimSpace(lang))

	// Handle Chinese variants
	if strings.HasPrefix(lang, "zh") {
		return LanguageChinese
	}

	// Handle English variants
	if strings.HasPrefix(lang, "en") {
		return LanguageEnglish
	}

	// Check exact matches
	switch lang {
	case "zh", "zh-cn", "zh-hans", "chinese":
		return LanguageChinese
	case "en", "en-us", "en-gb", "english":
		return LanguageEnglish
	}

	// Default to English for unsupported languages
	return LanguageEnglish
}

// GetMessage returns localized message for given key and language
func GetMessage(lang Language, key string) string {
	if messages, exists := Messages[lang]; exists {
		if message, exists := messages[key]; exists {
			return message
		}
	}

	// Fallback to English
	if lang != LanguageEnglish {
		if messages, exists := Messages[LanguageEnglish]; exists {
			if message, exists := messages[key]; exists {
				return message
			}
		}
	}

	// Final fallback to key itself
	return key
}

// GetUserFriendlyMessage returns a user-friendly error message, hiding internal details
func GetUserFriendlyMessage(lang Language, errorCode string, originalMessage string) string {
	// Try to get message by error code first
	if errorCode != "" {
		if message := GetMessage(lang, errorCode); message != errorCode {
			return message
		}
	}

	// Map common error patterns to user-friendly messages
	lowerMsg := strings.ToLower(originalMessage)
	switch {
	case strings.Contains(lowerMsg, "weak password"):
		return GetMessage(lang, "weak_password")
	case strings.Contains(lowerMsg, "duplicate") && strings.Contains(lowerMsg, "email"):
		return GetMessage(lang, "duplicate_email")
	case strings.Contains(lowerMsg, "duplicate") && strings.Contains(lowerMsg, "phone"):
		return GetMessage(lang, "duplicate_phone")
	case strings.Contains(lowerMsg, "captcha"):
		return GetMessage(lang, "captcha_failed")
	case strings.Contains(lowerMsg, "validation"):
		return GetMessage(lang, "validation_failed")
	case strings.Contains(lowerMsg, "json"):
		return GetMessage(lang, "bad_json")
	case strings.Contains(lowerMsg, "unauthorized"):
		return GetMessage(lang, "unauthorized")
	case strings.Contains(lowerMsg, "forbidden"):
		return GetMessage(lang, "forbidden")
	case strings.Contains(lowerMsg, "not found"):
		return GetMessage(lang, "not_found")
	case strings.Contains(lowerMsg, "too many"):
		return GetMessage(lang, "too_many_requests")
	case strings.Contains(lowerMsg, "conflict"):
		return GetMessage(lang, "conflict")
	case strings.Contains(lowerMsg, "internal") ||
		strings.Contains(lowerMsg, "server error") ||
		strings.Contains(lowerMsg, "database") ||
		strings.Contains(lowerMsg, "sql") ||
		strings.Contains(lowerMsg, "pq:"):
		// For internal server errors, always return a generic message
		return GetMessage(lang, "internal_server_error")
	default:
		// For any other error, return a generic message to hide internal details
		return GetMessage(lang, "unknown_error")
	}
}

// parseAcceptLanguage parses Accept-Language header with quality values
func parseAcceptLanguage(acceptLang string) Language {
	type langQuality struct {
		lang    string
		quality float64
	}

	var languages []langQuality

	// Split by comma to get individual language preferences
	parts := strings.Split(acceptLang, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// Split by semicolon to separate language from quality
		langParts := strings.Split(part, ";")
		lang := strings.TrimSpace(langParts[0])
		quality := 1.0 // Default quality

		// Parse quality value if present
		if len(langParts) > 1 {
			for _, param := range langParts[1:] {
				param = strings.TrimSpace(param)
				if strings.HasPrefix(param, "q=") {
					if q, err := strconv.ParseFloat(param[2:], 64); err == nil {
						quality = q
					}
				}
			}
		}

		languages = append(languages, langQuality{lang: lang, quality: quality})
	}

	// Sort by quality (highest first)
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].quality > languages[j].quality
	})

	// Return the first supported language
	for _, langQ := range languages {
		normalizedLang := normalizeLanguage(langQ.lang)
		// Check if it's a supported language (not default fallback)
		if normalizedLang == LanguageChinese || normalizedLang == LanguageEnglish {
			if strings.HasPrefix(strings.ToLower(langQ.lang), string(normalizedLang)) {
				return normalizedLang
			}
		}
	}

	// Default to English if no supported language found
	return LanguageEnglish
}
