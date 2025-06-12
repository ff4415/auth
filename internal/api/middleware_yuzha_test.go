package api

import (
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supabase/auth/internal/i18n"
)

func TestAPI_verifyAndSetLanguage(t *testing.T) {
	api := &API{}

	tests := []struct {
		name           string
		acceptLanguage string
		langParam      string
		expected       i18n.Language
		description    string
	}{
		{
			name:           "Default English",
			acceptLanguage: "",
			langParam:      "",
			expected:       i18n.LanguageEnglish,
			description:    "无任何语言信息时默认英文",
		},
		{
			name:           "English via Accept-Language",
			acceptLanguage: "en-US,en;q=0.9",
			langParam:      "",
			expected:       i18n.LanguageEnglish,
			description:    "通过Accept-Language头检测英文",
		},
		{
			name:           "Chinese via Accept-Language",
			acceptLanguage: "zh-CN,zh;q=0.9,en;q=0.8",
			langParam:      "",
			expected:       i18n.LanguageChinese,
			description:    "通过Accept-Language头检测中文",
		},
		{
			name:           "Chinese via query parameter (higher priority)",
			acceptLanguage: "en-US,en;q=0.9",
			langParam:      "zh",
			expected:       i18n.LanguageChinese,
			description:    "查询参数优先级高于Accept-Language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 创建测试请求
			req := httptest.NewRequest("GET", "/test", nil)
			if tt.acceptLanguage != "" {
				req.Header.Set("Accept-Language", tt.acceptLanguage)
			}
			if tt.langParam != "" {
				q := req.URL.Query()
				q.Add("lang", tt.langParam)
				req.URL.RawQuery = q.Encode()
			}

			w := httptest.NewRecorder()

			// 调用中间件函数
			ctx, err := api.verifyAndSetLanguage(w, req)

			// 验证结果
			assert.NoError(t, err, "中间件不应该返回错误")
			assert.NotNil(t, ctx, "应该返回有效的上下文")

			// 验证语言是否正确设置到上下文中
			lang := ctx.Value(i18n.UserLanguageKey)
			assert.NotNil(t, lang, "语言应该被设置到上下文中")
			assert.Equal(t, tt.expected, lang, "语言检测结果应该匹配期望值: %s", tt.description)
		})
	}
}

func TestAPI_verifyAndSetLanguage_ContextIntegration(t *testing.T) {
	api := &API{}

	// 测试上下文集成
	req := httptest.NewRequest("GET", "/test?lang=zh", nil)
	w := httptest.NewRecorder()

	ctx, err := api.verifyAndSetLanguage(w, req)
	assert.NoError(t, err)

	// 验证可以从上下文中获取语言
	lang := i18n.GetLanguageFromContext(ctx)
	assert.Equal(t, i18n.LanguageChinese, lang, "应该能够从上下文中获取中文语言设置")
}
