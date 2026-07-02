package i18n

import "context"

type Messages map[Language]string

func Pick(ctx context.Context, messages Messages, fallback string) string {
	language := FromContext(ctx)

	if message, ok := messages[language]; ok && message != "" {
		return message
	}

	if message, ok := messages[DefaultLanguage]; ok && message != "" {
		return message
	}

	return fallback
}

func ByLanguage(language Language, messages Messages, fallback string) string {
	language = NormalizeLanguage(string(language))

	if message, ok := messages[language]; ok && message != "" {
		return message
	}

	if message, ok := messages[DefaultLanguage]; ok && message != "" {
		return message
	}

	return fallback
}
