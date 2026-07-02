package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Rumm1/eduhub-backend/internal/shared/i18n"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", t(r.Context(),
			"Некорректное тело запроса",
			"Invalid request body",
			"Сұраныс денесі қате",
		))
		return
	}

	result, err := h.service.Login(r.Context(), req)
	if err != nil {
		writeAuthError(w, r.Context(), err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	result, err := h.service.Me(r.Context())
	if err != nil {
		writeAuthError(w, r.Context(), err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) SwitchProfile(w http.ResponseWriter, r *http.Request) {
	var req SwitchProfileRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", t(r.Context(),
			"Некорректное тело запроса",
			"Invalid request body",
			"Сұраныс денесі қате",
		))
		return
	}

	result, err := h.service.SwitchProfile(r.Context(), req)
	if err != nil {
		writeAuthError(w, r.Context(), err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req ChangePasswordRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, http.StatusBadRequest, "INVALID_JSON", t(r.Context(),
			"Некорректное тело запроса",
			"Invalid request body",
			"Сұраныс денесі қате",
		))
		return
	}

	result, err := h.service.ChangePassword(r.Context(), req)
	if err != nil {
		writeAuthError(w, r.Context(), err)
		return
	}

	response.Success(w, http.StatusOK, result)
}

func writeAuthError(w http.ResponseWriter, ctx context.Context, err error) {
	switch {
	case errors.Is(err, ErrInvalidCredentials):
		response.Error(w, http.StatusUnauthorized, "INVALID_CREDENTIALS", t(ctx,
			"Неверный email или пароль",
			"Invalid email or password",
			"Email немесе құпиясөз қате",
		))
	case errors.Is(err, ErrUserInactive):
		response.Error(w, http.StatusForbidden, "USER_INACTIVE", t(ctx,
			"Пользователь неактивен",
			"User is inactive",
			"Пайдаланушы белсенді емес",
		))
	case errors.Is(err, ErrProfileInactive):
		response.Error(w, http.StatusForbidden, "PROFILE_INACTIVE", t(ctx,
			"Профиль неактивен",
			"Profile is inactive",
			"Профиль белсенді емес",
		))
	case errors.Is(err, ErrProfileIDInvalid):
		response.Error(w, http.StatusBadRequest, "PROFILE_ID_INVALID", t(ctx,
			"Некорректный ID профиля",
			"Profile id is invalid",
			"Профиль ID қате",
		))
	case errors.Is(err, ErrUserContextMissing):
		response.Error(w, http.StatusUnauthorized, "UNAUTHORIZED", t(ctx,
			"Необходима авторизация",
			"Unauthorized",
			"Авторизация қажет",
		))
	case errors.Is(err, ErrCurrentPasswordMissing):
		response.Error(w, http.StatusBadRequest, "CURRENT_PASSWORD_REQUIRED", t(ctx,
			"Введите текущий пароль",
			"Current password is required",
			"Қазіргі құпиясөзді енгізіңіз",
		))
	case errors.Is(err, ErrNewPasswordMissing):
		response.Error(w, http.StatusBadRequest, "NEW_PASSWORD_REQUIRED", t(ctx,
			"Введите новый пароль",
			"New password is required",
			"Жаңа құпиясөзді енгізіңіз",
		))
	case errors.Is(err, ErrNewPasswordTooShort):
		response.Error(w, http.StatusBadRequest, "NEW_PASSWORD_TOO_SHORT", t(ctx,
			"Новый пароль должен содержать минимум 8 символов",
			"New password must be at least 8 characters",
			"Жаңа құпиясөз кемінде 8 таңбадан тұруы керек",
		))
	case errors.Is(err, ErrNewPasswordSame):
		response.Error(w, http.StatusBadRequest, "NEW_PASSWORD_SAME_AS_CURRENT", t(ctx,
			"Новый пароль должен отличаться от текущего",
			"New password must be different from current password",
			"Жаңа құпиясөз қазіргі құпиясөзден өзгеше болуы керек",
		))
	default:
		response.Error(w, http.StatusInternalServerError, "INTERNAL_ERROR", t(ctx,
			"Внутренняя ошибка сервера",
			"Internal server error",
			"Сервердің ішкі қатесі",
		))
	}
}

func t(ctx context.Context, ru string, en string, kk string) string {
	return i18n.Pick(ctx, i18n.Messages{
		i18n.LanguageRU: ru,
		i18n.LanguageEN: en,
		i18n.LanguageKK: kk,
	}, ru)
}
