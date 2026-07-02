package file

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("files.read"))

		r.Get("/", handler.List)
		r.Get("/{fileID}", handler.GetByID)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("files.upload"))

		r.Post("/upload", handler.Upload)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("files.delete"))

		r.Delete("/{fileID}", handler.Delete)
	})
}
