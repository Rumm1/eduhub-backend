package middleware

import (
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/i18n"
)

const HeaderLanguage = "X-Language"

func Language(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		language := resolveLanguage(r)

		w.Header().Set("Content-Language", string(language))
		w.Header().Set(HeaderLanguage, string(language))
		w.Header().Add("Vary", HeaderLanguage)
		w.Header().Add("Vary", "Accept-Language")

		ctx := i18n.WithLanguage(r.Context(), language)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func resolveLanguage(r *http.Request) i18n.Language {
	headerLanguage := r.Header.Get(HeaderLanguage)
	if headerLanguage != "" {
		return i18n.NormalizeLanguage(headerLanguage)
	}

	acceptLanguage := r.Header.Get("Accept-Language")
	if acceptLanguage != "" {
		return i18n.NormalizeLanguage(acceptLanguage)
	}

	return i18n.DefaultLanguage
}
