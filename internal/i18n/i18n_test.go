package i18n

import (
	"net/http"
	"net/url"
	"testing"
)

func TestGetLanguageFromRequest(t *testing.T) {
	tests := []struct {
		name     string
		setupReq func(*http.Request)
		expected Language
	}{
		{
			name: "Chinese from query parameter",
			setupReq: func(r *http.Request) {
				r.URL.RawQuery = "lang=zh"
			},
			expected: LanguageChinese,
		},
		{
			name: "English from query parameter",
			setupReq: func(r *http.Request) {
				r.URL.RawQuery = "lang=en"
			},
			expected: LanguageEnglish,
		},
		{
			name: "Chinese from custom header",
			setupReq: func(r *http.Request) {
				r.Header.Set("X-Language", "zh-CN")
			},
			expected: LanguageChinese,
		},
		{
			name: "English from Accept-Language header",
			setupReq: func(r *http.Request) {
				r.Header.Set("Accept-Language", "en-US,en;q=0.9")
			},
			expected: LanguageEnglish,
		},
		{
			name: "Chinese from Accept-Language header",
			setupReq: func(r *http.Request) {
				r.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
			},
			expected: LanguageChinese,
		},
		{
			name: "Default to English",
			setupReq: func(r *http.Request) {
				// No language preferences
			},
			expected: LanguageEnglish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &http.Request{
				URL:    &url.URL{},
				Header: make(http.Header),
			}
			tt.setupReq(req)

			result := GetLanguageFromRequest(req)
			if result != tt.expected {
				t.Errorf("GetLanguageFromRequest() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name     string
		lang     Language
		key      string
		expected string
	}{
		{
			name:     "English weak password message",
			lang:     LanguageEnglish,
			key:      "weak_password",
			expected: "Password does not meet security requirements",
		},
		{
			name:     "Chinese weak password message",
			lang:     LanguageChinese,
			key:      "weak_password",
			expected: "密码不符合安全要求",
		},
		{
			name:     "Unknown key fallback to English",
			lang:     LanguageChinese,
			key:      "unknown_key",
			expected: "unknown_key",
		},
		{
			name:     "English duplicate email message",
			lang:     LanguageEnglish,
			key:      "duplicate_email",
			expected: "A user with this email address has already been registered",
		},
		{
			name:     "Chinese duplicate email message",
			lang:     LanguageChinese,
			key:      "duplicate_email",
			expected: "该邮箱地址已被注册",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMessage(tt.lang, tt.key)
			if result != tt.expected {
				t.Errorf("GetMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetUserFriendlyMessage(t *testing.T) {
	tests := []struct {
		name            string
		lang            Language
		errorCode       string
		originalMessage string
		expected        string
	}{
		{
			name:            "Hide internal SQL error - English",
			lang:            LanguageEnglish,
			errorCode:       "",
			originalMessage: "pq: duplicate key value violates unique constraint",
			expected:        "Internal server error",
		},
		{
			name:            "Hide internal SQL error - Chinese",
			lang:            LanguageChinese,
			errorCode:       "",
			originalMessage: "database connection failed",
			expected:        "内部服务器错误",
		},
		{
			name:            "Weak password error - English",
			lang:            LanguageEnglish,
			errorCode:       "",
			originalMessage: "Weak password detected",
			expected:        "Password does not meet security requirements",
		},
		{
			name:            "Weak password error - Chinese",
			lang:            LanguageChinese,
			errorCode:       "",
			originalMessage: "Weak password detected",
			expected:        "密码不符合安全要求",
		},
		{
			name:            "Duplicate email error - English",
			lang:            LanguageEnglish,
			errorCode:       "",
			originalMessage: "Duplicate email found in database",
			expected:        "A user with this email address has already been registered",
		},
		{
			name:            "Duplicate email error - Chinese",
			lang:            LanguageChinese,
			errorCode:       "",
			originalMessage: "Duplicate email found in database",
			expected:        "该邮箱地址已被注册",
		},
		{
			name:            "Error code mapping - English",
			lang:            LanguageEnglish,
			errorCode:       "unauthorized",
			originalMessage: "Some internal message",
			expected:        "Unauthorized",
		},
		{
			name:            "Error code mapping - Chinese",
			lang:            LanguageChinese,
			errorCode:       "unauthorized",
			originalMessage: "Some internal message",
			expected:        "未授权",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserFriendlyMessage(tt.lang, tt.errorCode, tt.originalMessage)
			if result != tt.expected {
				t.Errorf("GetUserFriendlyMessage() = %v, want %v", result, tt.expected)
			}
		})
	}
}
