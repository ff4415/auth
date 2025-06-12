package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetLanguageFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		customHeader   string
		acceptLanguage string
		expected       Language
	}{
		{
			name:           "Query parameter takes priority",
			queryParam:     "zh",
			customHeader:   "en",
			acceptLanguage: "fr",
			expected:       LanguageChinese,
		},
		{
			name:           "Custom header when no query param",
			customHeader:   "zh",
			acceptLanguage: "en",
			expected:       LanguageChinese,
		},
		{
			name:           "Accept-Language when others missing",
			acceptLanguage: "zh-CN,zh;q=0.9,en;q=0.8",
			expected:       LanguageChinese,
		},
		{
			name:     "Default to English when all missing",
			expected: LanguageEnglish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)

			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Set("lang", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			if tt.customHeader != "" {
				req.Header.Set("X-Language", tt.customHeader)
			}

			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}

			result := GetLanguageFromRequest(req)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		name       string
		acceptLang string
		expected   Language
	}{
		{
			name:       "Chinese with quality values",
			acceptLang: "zh-CN,zh;q=0.9,en;q=0.8",
			expected:   LanguageChinese,
		},
		{
			name:       "English with quality values",
			acceptLang: "en-US,en;q=0.9,zh;q=0.8",
			expected:   LanguageEnglish,
		},
		{
			name:       "Chinese higher priority",
			acceptLang: "fr;q=0.5,zh;q=0.9,en;q=0.8",
			expected:   LanguageChinese,
		},
		{
			name:       "No supported languages",
			acceptLang: "fr,de,es",
			expected:   LanguageEnglish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAcceptLanguage(tt.acceptLang)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestGetLanguageFromJWT(t *testing.T) {
	tests := []struct {
		name     string
		claims   map[string]interface{}
		expected Language
	}{
		{
			name: "User metadata language preference",
			claims: map[string]interface{}{
				"user_metadata": map[string]interface{}{
					"user_language": "zh",
				},
			},
			expected: LanguageChinese,
		},
		{
			name: "App metadata fallback",
			claims: map[string]interface{}{
				"app_metadata": map[string]interface{}{
					"user_language": "zh",
				},
			},
			expected: LanguageChinese,
		},
		{
			name: "No language preference",
			claims: map[string]interface{}{
				"user_metadata": map[string]interface{}{},
			},
			expected: LanguageEnglish,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetLanguageFromJWT(tt.claims)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestLanguageMiddleware(t *testing.T) {
	// Create test handler
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := GetLanguageFromContextHTTP(r)
		w.WriteHeader(200)
		w.Write([]byte(string(lang)))
	})

	// Wrap with middleware
	handler := LanguageMiddleware(testHandler)

	t.Run("Middleware sets language in context", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test?lang=zh", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		if w.Body.String() != "zh" {
			t.Errorf("Expected 'zh', got '%s'", w.Body.String())
		}
	})
}

func TestGetMessage(t *testing.T) {
	tests := []struct {
		name     string
		lang     Language
		key      string
		expected string
	}{
		{
			name:     "English message",
			lang:     LanguageEnglish,
			key:      "unauthorized",
			expected: "Unauthorized",
		},
		{
			name:     "Chinese message",
			lang:     LanguageChinese,
			key:      "unauthorized",
			expected: "未授权",
		},
		{
			name:     "Fallback to English",
			lang:     LanguageChinese,
			key:      "nonexistent_key",
			expected: "nonexistent_key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetMessage(tt.lang, tt.key)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
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
			name:            "Database error should be hidden",
			lang:            LanguageEnglish,
			originalMessage: "pq: connection refused",
			expected:        "Internal server error",
		},
		{
			name:            "Weak password error",
			lang:            LanguageChinese,
			originalMessage: "weak password detected",
			expected:        "密码不符合安全要求",
		},
		{
			name:            "Generic error fallback",
			lang:            LanguageEnglish,
			originalMessage: "some unknown error",
			expected:        "Unknown error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetUserFriendlyMessage(tt.lang, tt.errorCode, tt.originalMessage)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
