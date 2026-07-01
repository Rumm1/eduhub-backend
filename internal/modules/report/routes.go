package report

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("reports.teacher_schedule.read")).Get("/teacher-schedule", handler.GetTeacherSchedule)
}
