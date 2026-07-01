package report

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("reports.teacher_schedule.read")).Get("/teacher-schedule", handler.GetTeacherSchedule)
	r.With(middleware.RequirePermission("reports.payments.read")).Get("/payments", handler.GetPaymentsReport)
	r.With(middleware.RequirePermission("reports.student_balance.read")).Get("/student-balances", handler.GetStudentBalancesReport)
}
