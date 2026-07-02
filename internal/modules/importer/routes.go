package importer

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("students.create"))
		r.Use(middleware.RequirePermission("parents.create"))

		r.Post("/students/preview", handler.PreviewStudents)
		r.Post("/students/confirm", handler.ConfirmStudents)
	})
}
