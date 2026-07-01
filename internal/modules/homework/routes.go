package homework

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("homeworks.read")).Get("/", handler.List)
	r.With(middleware.RequirePermission("homeworks.manage")).Post("/", handler.Create)

	r.With(middleware.RequirePermission("homeworks.read")).Get("/lessons/{lessonID}", handler.ListByLessonID)
}
