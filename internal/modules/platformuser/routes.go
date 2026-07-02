package platformuser

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequireRole("SUPER_ADMIN"))

		r.Post("/users/{userID}/reset-password", handler.ResetPassword)
	})
}
