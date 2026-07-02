package parent

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("parents.read"))

		r.Get("/", handler.List)
		r.Get("/{parentID}", handler.GetByID)
		r.Get("/{parentID}/students", handler.ListStudents)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("parents.manage"))

		r.Post("/", handler.Create)
		r.Patch("/{parentID}", handler.Update)
		r.Delete("/{parentID}", handler.Delete)

		r.Post("/{parentID}/students/{studentID}", handler.AttachStudent)
		r.Delete("/{parentID}/students/{studentID}", handler.DetachStudent)
	})
}
