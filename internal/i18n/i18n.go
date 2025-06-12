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
		"weak_password":           "Password does not meet security requirements",
		"unexpected_failure":      "Unexpected failure, please check server logs for more information",
		"unknown_error":           "Unknown error occurred",
		"validation_failed":       "Validation failed",
		"bad_json":                "Could not parse request body as JSON",
		"internal_server_error":   "Internal server error",
		"unauthorized":            "Unauthorized",
		"forbidden":               "Forbidden",
		"not_found":               "Not found",
		"too_many_requests":       "Too many requests",
		"conflict":                "Conflict",
		"unprocessable_entity":    "Unprocessable entity",
		"bad_request":             "Bad request",
		"duplicate_email":         "A user with this email address has already been registered",
		"duplicate_phone":         "A user with this phone number has already been registered",
		"captcha_failed":          "Captcha verification failed",
		"email_not_confirmed":     "Email not confirmed",
		"phone_not_confirmed":     "Phone not confirmed",
		"mfa_verification_failed": "MFA verification failed",
		"insufficient_aal":        "Insufficient authentication assurance level",
	},
	LanguageChinese: {
		"weak_password":           "密码不符合安全要求",
		"unexpected_failure":      "意外错误，请检查服务器日志以获取更多信息",
		"unknown_error":           "发生未知错误",
		"validation_failed":       "验证失败",
		"bad_json":                "无法解析请求体为JSON格式",
		"internal_server_error":   "内部服务器错误",
		"unauthorized":            "未授权",
		"forbidden":               "禁止访问",
		"not_found":               "未找到",
		"too_many_requests":       "请求过于频繁",
		"conflict":                "冲突",
		"unprocessable_entity":    "无法处理的实体",
		"bad_request":             "错误的请求",
		"duplicate_email":         "该邮箱地址已被注册",
		"duplicate_phone":         "该手机号码已被注册",
		"captcha_failed":          "验证码验证失败",
		"email_not_confirmed":     "邮箱未确认",
		"phone_not_confirmed":     "手机号未确认",
		"mfa_verification_failed": "多因子认证验证失败",
		"insufficient_aal":        "认证保证级别不足",
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
