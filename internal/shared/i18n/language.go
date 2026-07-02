package i18n

import (
	"context"
	"strings"
)

type Language string

const (
	LanguageRU Language = "ru"
	LanguageEN Language = "en"
	LanguageKK Language = "kk"
)

const DefaultLanguage = LanguageRU

type contextKey string

const languageContextKey contextKey = "language"

var supportedLanguages = map[Language]struct{}{
	LanguageRU: {},
	LanguageEN: {},
	LanguageKK: {},
}

func WithLanguage(ctx context.Context, language Language) context.Context {
	return context.WithValue(ctx, languageContextKey, NormalizeLanguage(string(language)))
}

func FromContext(ctx context.Context) Language {
	value, ok := ctx.Value(languageContextKey).(Language)
	if !ok {
		return DefaultLanguage
	}

	return NormalizeLanguage(string(value))
}

func NormalizeLanguage(value string) Language {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return DefaultLanguage
	}

	parts := strings.Split(value, ",")
	for _, part := range parts {
		language := normalizeSingleLanguage(part)
		if _, ok := supportedLanguages[language]; ok {
			return language
		}
	}

	return DefaultLanguage
}

func IsSupported(language Language) bool {
	_, ok := supportedLanguages[language]
	return ok
}

func normalizeSingleLanguage(value string) Language {
	value = strings.TrimSpace(strings.ToLower(value))
	if value == "" {
		return DefaultLanguage
	}

	value = strings.Split(value, ";")[0]
	value = strings.TrimSpace(value)
	value = strings.ReplaceAll(value, "_", "-")

	if strings.Contains(value, "-") {
		value = strings.Split(value, "-")[0]
	}

	switch value {
	case "ru":
		return LanguageRU
	case "en":
		return LanguageEN
	case "kk", "kz":
		return LanguageKK
	default:
		return Language(value)
	}
}
