package api

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supabase/auth/internal/i18n"
)

func TestDetectLanguageFromRequest(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		customHeader   string
		acceptLanguage string
		expected       string
		description    string
	}{
		{
			name:           "Query parameter takes priority",
			queryParam:     "zh",
			customHeader:   "en",
			acceptLanguage: "fr",
			expected:       "zh",
			description:    "查询参数具有最高优先级",
		},
		{
			name:           "Custom header when no query param",
			customHeader:   "zh",
			acceptLanguage: "en",
			expected:       "zh",
			description:    "自定义头优先于Accept-Language",
		},
		{
			name:           "Accept-Language when others missing",
			acceptLanguage: "zh-CN,zh;q=0.9,en;q=0.8",
			expected:       "zh",
			description:    "Accept-Language头解析",
		},
		{
			name:        "Default to English when all missing",
			expected:    "en",
			description: "默认英文",
		},
		{
			name:           "Invalid language defaults to English",
			queryParam:     "fr",
			acceptLanguage: "fr-FR,fr;q=0.9",
			expected:       "en",
			description:    "不支持的语言回退到英文",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)

			// 设置查询参数
			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Add("lang", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			// 设置自定义头
			if tt.customHeader != "" {
				req.Header.Set("X-Language", tt.customHeader)
			}

			// 设置Accept-Language头
			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}

			result := detectLanguageFromRequest(req)
			assert.Equal(t, tt.expected, result, "场景: %s", tt.description)
		})
	}
}

func TestNormalizeLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"en", "en"},
		{"EN", "en"},
		{"en-US", "en"},
		{"en-GB", "en"},
		{"english", "en"},
		{"zh", "zh"},
		{"ZH", "zh"},
		{"zh-CN", "zh"},
		{"zh-TW", "zh"},
		{"chinese", "zh"},
		{"fr", "en"}, // 不支持的语言默认英文
		{"", "en"},   // 空字符串默认英文
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := normalizeLanguage(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsValidLanguage(t *testing.T) {
	tests := []struct {
		lang     string
		expected bool
	}{
		{"en", true},
		{"zh", true},
		{"fr", false},
		{"", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(tt.lang, func(t *testing.T) {
			result := isValidLanguage(tt.lang)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestParseAcceptLanguageHeader(t *testing.T) {
	tests := []struct {
		name        string
		acceptLang  string
		expected    string
		description string
	}{
		{
			name:        "Chinese with quality values",
			acceptLang:  "zh-CN,zh;q=0.9,en;q=0.8",
			expected:    "zh",
			description: "中文具有最高权重",
		},
		{
			name:        "English with quality values",
			acceptLang:  "en-US,en;q=0.9,zh;q=0.8",
			expected:    "en",
			description: "英文具有最高权重",
		},
		{
			name:        "Mixed languages with Chinese priority",
			acceptLang:  "fr;q=0.7,zh-CN;q=0.9,en;q=0.8",
			expected:    "zh",
			description: "中文权重最高",
		},
		{
			name:        "Unsupported language falls back to English",
			acceptLang:  "fr-FR,fr;q=0.9,de;q=0.8",
			expected:    "en",
			description: "不支持的语言回退到英文",
		},
		{
			name:        "Empty header",
			acceptLang:  "",
			expected:    "en",
			description: "空头默认英文",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseAcceptLanguageHeader(tt.acceptLang)
			assert.Equal(t, tt.expected, result, "场景: %s", tt.description)
		})
	}
}

func TestDetectAndSetLanguage(t *testing.T) {
	// 创建一个模拟的API实例
	api := &API{}

	tests := []struct {
		name           string
		acceptLanguage string
		queryParam     string
		expected       string
	}{
		{
			name:           "Chinese from Accept-Language",
			acceptLanguage: "zh-CN,zh;q=0.9,en;q=0.8",
			expected:       "zh",
		},
		{
			name:       "English from query parameter",
			queryParam: "en",
			expected:   "en",
		},
		{
			name:     "Default English",
			expected: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// 设置测试参数
			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}
			if tt.queryParam != "" {
				q := req.URL.Query()
				q.Add("lang", tt.queryParam)
				req.URL.RawQuery = q.Encode()
			}

			// 调用中间件
			ctx, err := api.detectAndSetLanguage(w, req)

			// 验证结果
			assert.NoError(t, err)
			assert.NotNil(t, ctx)

			// 验证语言是否正确设置到上下文中
			lang := getLanguageFromContext(ctx)
			assert.Equal(t, tt.expected, lang)
		})
	}
}

func TestGetLanguageFromContext(t *testing.T) {
	// 测试有语言的上下文
	ctx := context.WithValue(context.Background(), i18n.UserLanguageKey, "zh")
	lang := getLanguageFromContext(ctx)
	assert.Equal(t, "zh", lang)

	// 测试空上下文
	emptyCtx := context.Background()
	lang = getLanguageFromContext(emptyCtx)
	assert.Equal(t, "en", lang)

	// 测试错误类型的值
	invalidCtx := context.WithValue(context.Background(), i18n.UserLanguageKey, 123)
	lang = getLanguageFromContext(invalidCtx)
	assert.Equal(t, "en", lang)
}

func TestGetLanguageFromContextHTTP(t *testing.T) {
	// 创建带有语言上下文的请求
	ctx := context.WithValue(context.Background(), i18n.UserLanguageKey, "zh")
	req := httptest.NewRequest("GET", "/test", nil)
	req = req.WithContext(ctx)

	lang := getLanguageFromContextHTTP(req)
	assert.Equal(t, "zh", lang)

	// 测试没有语言上下文的请求
	emptyReq := httptest.NewRequest("GET", "/test", nil)
	lang = getLanguageFromContextHTTP(emptyReq)
	assert.Equal(t, "en", lang)
}
