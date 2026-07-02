package notification

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("notifications.read"))

		r.Get("/", handler.List)
		r.Get("/types", handler.Types)
		r.Patch("/read-all", handler.MarkAllRead)
		r.Patch("/{notificationID}/read", handler.MarkRead)
		r.Delete("/{notificationID}", handler.Delete)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("notifications.manage"))

		r.Post("/", handler.Create)
	})
}
