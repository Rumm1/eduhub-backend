package attendance

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("attendance.read")).Get("/lessons/{lessonID}", handler.ListByLessonID)
	r.With(middleware.RequirePermission("attendance.manage")).Post("/lessons/{lessonID}/mark", handler.MarkLessonAttendance)
}
