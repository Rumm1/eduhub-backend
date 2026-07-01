package lesson

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("lessons.read")).Get("/", handler.List)
	r.With(middleware.RequirePermission("lessons.create")).Post("/", handler.Create)

	r.With(middleware.RequirePermission("lessons.update")).Patch("/{lessonID}/teacher", handler.UpdateActualTeacher)
}
