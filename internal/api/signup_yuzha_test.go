package api

import (
	"net/http/httptest"
	"testing"
)

func TestConfigureDefaultsYuzha(t *testing.T) {
	t.Run("Basic defaults when no data", func(t *testing.T) {
		params := &SignupParams{}
		params.ConfigureDefaultsYuzha()

		if params.Data["user_language"] != "en" {
			t.Errorf("Expected user_language 'en', got %v", params.Data["user_language"])
		}
		if params.Data["country"] != "US" {
			t.Errorf("Expected country 'US', got %v", params.Data["country"])
		}
	})

	t.Run("Preserve existing data", func(t *testing.T) {
		params := &SignupParams{
			Data: map[string]interface{}{
				"user_language": "zh",
				"country":       "CN",
				"custom_field":  "value",
			},
		}
		params.ConfigureDefaultsYuzha()

		if params.Data["user_language"] != "zh" {
			t.Errorf("Expected preserved user_language 'zh', got %v", params.Data["user_language"])
		}
		if params.Data["country"] != "CN" {
			t.Errorf("Expected preserved country 'CN', got %v", params.Data["country"])
		}
		if params.Data["custom_field"] != "value" {
			t.Errorf("Expected preserved custom_field 'value', got %v", params.Data["custom_field"])
		}
	})

	t.Run("Fill missing defaults", func(t *testing.T) {
		params := &SignupParams{
			Data: map[string]interface{}{
				"user_language": "zh",
				// country is missing
			},
		}
		params.ConfigureDefaultsYuzha()

		if params.Data["user_language"] != "zh" {
			t.Errorf("Expected preserved user_language 'zh', got %v", params.Data["user_language"])
		}
		if params.Data["country"] != "US" {
			t.Errorf("Expected default country 'US', got %v", params.Data["country"])
		}
	})
}

func TestConfigureDefaultsYuzhaWithContext(t *testing.T) {
	tests := []struct {
		name             string
		acceptLanguage   string
		expectedLang     string
		expectedCountry  string
		expectedTimezone string
	}{
		{
			name:             "Chinese user from China",
			acceptLanguage:   "zh-CN,zh;q=0.9,en;q=0.8",
			expectedLang:     "zh",
			expectedCountry:  "CN",
			expectedTimezone: "Asia/Shanghai",
		},
		{
			name:             "English user from US",
			acceptLanguage:   "en-US,en;q=0.9",
			expectedLang:     "en",
			expectedCountry:  "US",
			expectedTimezone: "America/New_York",
		},
		{
			name:             "English user from UK",
			acceptLanguage:   "en-GB,en;q=0.9",
			expectedLang:     "en",
			expectedCountry:  "GB",
			expectedTimezone: "America/New_York", // Still en default
		},
		{
			name:             "No Accept-Language header",
			acceptLanguage:   "",
			expectedLang:     "en",
			expectedCountry:  "US",
			expectedTimezone: "America/New_York",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/auth/signup", nil)
			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}

			params := &SignupParams{}
			params.ConfigureDefaultsYuzhaWithContext(req)

			// Verify language
			if params.Data["user_language"] != tt.expectedLang {
				t.Errorf("Expected user_language '%s', got %v", tt.expectedLang, params.Data["user_language"])
			}

			// Verify country
			if params.Data["country"] != tt.expectedCountry {
				t.Errorf("Expected country '%s', got %v", tt.expectedCountry, params.Data["country"])
			}

			// Verify timezone
			if params.Data["timezone"] != tt.expectedTimezone {
				t.Errorf("Expected timezone '%s', got %v", tt.expectedTimezone, params.Data["timezone"])
			}

			// Verify additional fields are set
			if params.Data["date_format"] == nil {
				t.Error("Expected date_format to be set")
			}
			if params.Data["number_format"] == nil {
				t.Error("Expected number_format to be set")
			}
		})
	}
}

func TestDetectCountryFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		acceptLanguage string
		expected       string
	}{
		{
			name:           "Chinese from China",
			acceptLanguage: "zh-CN,zh;q=0.9",
			expected:       "CN",
		},
		{
			name:           "English from US",
			acceptLanguage: "en-US,en;q=0.9",
			expected:       "US",
		},
		{
			name:           "English from UK",
			acceptLanguage: "en-GB,en;q=0.9",
			expected:       "GB",
		},
		{
			name:           "Chinese from Hong Kong",
			acceptLanguage: "zh-HK,zh;q=0.9",
			expected:       "HK",
		},
		{
			name:           "Unsupported region",
			acceptLanguage: "fr-FR,fr;q=0.9",
			expected:       "US", // Default fallback
		},
		{
			name:           "No Accept-Language",
			acceptLanguage: "",
			expected:       "US", // Default fallback
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}

			result := detectCountryFromRequest(req)
			if result != tt.expected {
				t.Errorf("Expected country '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDefaultFormatFunctions(t *testing.T) {
	t.Run("Timezone defaults", func(t *testing.T) {
		tests := []struct {
			lang     string
			expected string
		}{
			{"zh", "Asia/Shanghai"},
			{"en", "America/New_York"},
			{"unknown", "UTC"},
		}

		for _, tt := range tests {
			result := getDefaultTimezone(tt.lang)
			if result != tt.expected {
				t.Errorf("getDefaultTimezone(%s) = %s, expected %s", tt.lang, result, tt.expected)
			}
		}
	})

	t.Run("Date format defaults", func(t *testing.T) {
		tests := []struct {
			lang     string
			expected string
		}{
			{"zh", "YYYY年MM月DD日"},
			{"en", "MM/DD/YYYY"},
			{"unknown", "YYYY-MM-DD"},
		}

		for _, tt := range tests {
			result := getDefaultDateFormat(tt.lang)
			if result != tt.expected {
				t.Errorf("getDefaultDateFormat(%s) = %s, expected %s", tt.lang, result, tt.expected)
			}
		}
	})

	t.Run("Number format defaults", func(t *testing.T) {
		tests := []struct {
			lang     string
			expected string
		}{
			{"zh", "zh-CN"},
			{"en", "en-US"},
			{"unknown", "en-US"},
		}

		for _, tt := range tests {
			result := getDefaultNumberFormat(tt.lang)
			if result != tt.expected {
				t.Errorf("getDefaultNumberFormat(%s) = %s, expected %s", tt.lang, result, tt.expected)
			}
		}
	})
}
