package i18n

import (
	"context"
	"net/http"
)

// JWTClaims represents JWT claims structure
type JWTClaims interface {
	GetUserMetadata() map[string]interface{}
}

// JWTClaimsKey is the context key for storing JWT claims
type JWTClaimsKey string

const ClaimsKey JWTClaimsKey = "jwt_claims"

// GetLanguageFromJWT extracts language preference from JWT claims
func GetLanguageFromJWT(claims map[string]interface{}) Language {
	if userMeta, ok := claims["user_metadata"].(map[string]interface{}); ok {
		if lang, ok := userMeta["user_language"].(string); ok && lang != "" {
			return normalizeLanguage(lang)
		}
	}

	// Check app_metadata as fallback
	if appMeta, ok := claims["app_metadata"].(map[string]interface{}); ok {
		if lang, ok := appMeta["user_language"].(string); ok && lang != "" {
			return normalizeLanguage(lang)
		}
	}

	return LanguageEnglish
}

// GetLanguageFromRequestWithJWT detects language with JWT claims priority
func GetLanguageFromRequestWithJWT(r *http.Request) Language {
	// 1. Check query parameter first (highest priority)
	if lang := r.URL.Query().Get("lang"); lang != "" {
		return normalizeLanguage(lang)
	}

	// 2. Check custom header
	if lang := r.Header.Get("X-Language"); lang != "" {
		return normalizeLanguage(lang)
	}

	// 3. Check JWT claims from context
	if claims, ok := r.Context().Value(ClaimsKey).(map[string]interface{}); ok {
		if lang := GetLanguageFromJWT(claims); lang != LanguageEnglish {
			return lang
		}
	}

	// 4. Check Accept-Language header
	if acceptLang := r.Header.Get("Accept-Language"); acceptLang != "" {
		return parseAcceptLanguage(acceptLang)
	}

	// 5. Default to English
	return LanguageEnglish
}

// EnhancedLanguageMiddleware detects language with JWT support
func EnhancedLanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Use enhanced language detection with JWT support
		lang := GetLanguageFromRequestWithJWT(r)

		// Store in context
		ctx := context.WithValue(r.Context(), UserLanguageKey, lang)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
