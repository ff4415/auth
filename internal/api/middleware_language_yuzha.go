package api

import (
	"context"
	"net/http"
	sortpkg "sort"
	"strconv"
	"strings"

	"github.com/supabase/auth/internal/i18n"
)

// detectAndSetLanguage 检测并设置用户语言偏好到上下文
func (a *API) detectAndSetLanguage(w http.ResponseWriter, req *http.Request) (context.Context, error) {
	ctx := req.Context()

	// 检测用户语言偏好
	userLanguage := detectLanguageFromRequest(req)

	// 将语言偏好存储到上下文中
	ctx = context.WithValue(ctx, i18n.UserLanguageKey, userLanguage)

	return ctx, nil
}

// detectLanguageFromRequest 从HTTP请求中检测用户语言偏好
func detectLanguageFromRequest(req *http.Request) string {
	// 1. 检查查询参数 (最高优先级)
	if lang := req.URL.Query().Get("lang"); lang != "" {
		if normalized := normalizeLanguage(lang); isValidLanguage(normalized) {
			return normalized
		}
	}

	// 2. 检查自定义头
	if lang := req.Header.Get("X-Language"); lang != "" {
		if normalized := normalizeLanguage(lang); isValidLanguage(normalized) {
			return normalized
		}
	}

	// 3. 解析Accept-Language头
	if acceptLang := req.Header.Get("Accept-Language"); acceptLang != "" {
		if lang := parseAcceptLanguageHeader(acceptLang); isValidLanguage(lang) {
			return lang
		}
	}

	// 4. 默认返回英文
	return string(i18n.LanguageEnglish)
}

// normalizeLanguage 标准化语言代码
func normalizeLanguage(lang string) string {
	lang = strings.ToLower(strings.TrimSpace(lang))

	// 处理中文变体
	if strings.HasPrefix(lang, "zh") {
		return string(i18n.LanguageChinese)
	}

	// 处理英文变体
	if strings.HasPrefix(lang, "en") {
		return string(i18n.LanguageEnglish)
	}

	// 检查精确匹配
	switch lang {
	case "zh", "zh-cn", "zh-hans", "chinese":
		return string(i18n.LanguageChinese)
	case "en", "en-us", "en-gb", "english":
		return string(i18n.LanguageEnglish)
	}

	// 不支持的语言默认返回英文
	return string(i18n.LanguageEnglish)
}

// isValidLanguage 验证语言代码是否有效
func isValidLanguage(lang string) bool {
	validLanguages := map[string]bool{
		string(i18n.LanguageEnglish): true,
		string(i18n.LanguageChinese): true,
	}
	return validLanguages[lang]
}

// parseAcceptLanguageHeader 解析Accept-Language头，支持权重值
func parseAcceptLanguageHeader(acceptLang string) string {
	type langQuality struct {
		lang    string
		quality float64
	}

	var languages []langQuality

	// 按逗号分割获取各个语言偏好
	parts := strings.Split(acceptLang, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)

		// 按分号分割语言和权重
		langParts := strings.Split(part, ";")
		lang := strings.TrimSpace(langParts[0])
		quality := 1.0 // 默认权重

		// 解析权重值
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

	// 按权重排序（从高到低）
	sortpkg.Slice(languages, func(i, j int) bool {
		return languages[i].quality > languages[j].quality
	})

	// 返回第一个支持的语言
	for _, langQ := range languages {
		normalizedLang := normalizeLanguage(langQ.lang)
		if isValidLanguage(normalizedLang) {
			// 验证是否确实匹配原始语言前缀
			if strings.HasPrefix(strings.ToLower(langQ.lang), normalizedLang) {
				return normalizedLang
			}
		}
	}

	// 如果没有找到支持的语言，默认返回英文
	return string(i18n.LanguageEnglish)
}

// getLanguageFromContext 从上下文中获取用户语言偏好
func getLanguageFromContext(ctx context.Context) string {
	if lang, ok := ctx.Value(i18n.UserLanguageKey).(string); ok {
		return lang
	}
	return string(i18n.LanguageEnglish)
}

// getLanguageFromContextHTTP 从HTTP请求上下文中获取用户语言偏好
func getLanguageFromContextHTTP(req *http.Request) string {
	return getLanguageFromContext(req.Context())
}
