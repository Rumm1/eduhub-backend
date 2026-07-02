package auth

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler, jwtManager *platformjwt.Manager) {
	r.Post("/login", handler.Login)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth(jwtManager))

		r.Get("/me", handler.Me)
		r.Post("/switch-profile", handler.SwitchProfile)
		r.Post("/change-password", handler.ChangePassword)
	})
}
