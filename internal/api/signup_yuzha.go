package api

import (
	"net/http"
	"strings"

	"github.com/supabase/auth/internal/i18n"
)

// ConfigureDefaultsYuzhaWithContext 根据请求上下文智能配置默认值
func (p *SignupParams) ConfigureDefaultsYuzhaWithContext(r *http.Request) {
	if p.Data == nil {
		p.Data = make(map[string]interface{})
	}

	// 🌐 智能语言检测 - 优先级策略
	if _, ok := p.Data["user_language"]; !ok {
		// 1. 检测用户首选语言
		detectedLang := i18n.GetLanguageFromRequest(r)
		p.Data["user_language"] = string(detectedLang)
	}

	// 🌍 智能地区检测
	if _, ok := p.Data["country"]; !ok {
		country := detectCountryFromRequest(r)
		p.Data["country"] = country
	}

	// 🎨 根据语言设置其他默认偏好
	userLang := p.Data["user_language"].(string)
	if _, ok := p.Data["timezone"]; !ok {
		p.Data["timezone"] = getDefaultTimezone(userLang)
	}

	if _, ok := p.Data["date_format"]; !ok {
		p.Data["date_format"] = getDefaultDateFormat(userLang)
	}

	if _, ok := p.Data["number_format"]; !ok {
		p.Data["number_format"] = getDefaultNumberFormat(userLang)
	}
}

// detectCountryFromRequest 从请求中检测用户地区
func detectCountryFromRequest(r *http.Request) string {
	// 1. 检查Accept-Language头中的地区信息
	acceptLang := r.Header.Get("Accept-Language")
	if acceptLang != "" {
		// 解析如 "zh-CN", "en-US" 等格式
		parts := strings.Split(acceptLang, ",")
		if len(parts) > 0 {
			lang := strings.TrimSpace(strings.Split(parts[0], ";")[0])
			if strings.Contains(lang, "-") {
				langParts := strings.Split(lang, "-")
				if len(langParts) > 1 {
					region := strings.ToUpper(langParts[1])
					return mapRegionToCountry(region)
				}
			}
		}
	}

	// 2. 根据语言推断默认地区
	userLang := i18n.GetLanguageFromRequest(r)
	switch userLang {
	case i18n.LanguageChinese:
		return "CN"
	case i18n.LanguageEnglish:
		return "US"
	default:
		return "US"
	}
}

// getDefaultTimezone 根据语言获取默认时区
func getDefaultTimezone(lang string) string {
	switch lang {
	case "zh":
		return "Asia/Shanghai"
	case "en":
		return "America/New_York"
	default:
		return "UTC"
	}
}

// getDefaultDateFormat 根据语言获取默认日期格式
func getDefaultDateFormat(lang string) string {
	switch lang {
	case "zh":
		return "YYYY年MM月DD日"
	case "en":
		return "MM/DD/YYYY"
	default:
		return "YYYY-MM-DD"
	}
}

// getDefaultNumberFormat 根据语言获取默认数字格式
func getDefaultNumberFormat(lang string) string {
	switch lang {
	case "zh":
		return "zh-CN"
	case "en":
		return "en-US"
	default:
		return "en-US"
	}
}

// mapRegionToCountry 地区代码映射
func mapRegionToCountry(region string) string {
	regionMap := map[string]string{
		"CN": "CN", // 中国
		"HK": "HK", // 香港
		"TW": "TW", // 台湾
		"US": "US", // 美国
		"GB": "GB", // 英国
		"CA": "CA", // 加拿大
		"AU": "AU", // 澳大利亚
		"SG": "SG", // 新加坡
	}

	if country, exists := regionMap[region]; exists {
		return country
	}
	return "US" // 默认美国
}
