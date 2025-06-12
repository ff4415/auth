package i18n

import (
	"context"
	"net/http"
)

// LanguageContextKey is the context key for storing user language
type LanguageContextKey string

const UserLanguageKey LanguageContextKey = "user_language"

// LanguageMiddleware detects and sets user language preference in context
func LanguageMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Detect user language preference
		lang := GetLanguageFromRequest(r)

		// Store in context for later use
		ctx := context.WithValue(r.Context(), UserLanguageKey, lang)

		// Continue with next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetLanguageFromContext retrieves user language from request context
func GetLanguageFromContext(ctx context.Context) Language {
	if lang, ok := ctx.Value(UserLanguageKey).(Language); ok {
		return lang
	}
	return LanguageEnglish
}

// GetLanguageFromContextHTTP is a convenience function for HTTP handlers
func GetLanguageFromContextHTTP(r *http.Request) Language {
	return GetLanguageFromContext(r.Context())
}
