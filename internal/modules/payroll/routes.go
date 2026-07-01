package payroll

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("payroll.read_all")).Get("/entries", handler.ListEntries)
	r.With(middleware.RequirePermission("payroll.read_own")).Get("/entries/my", handler.ListMyEntries)

	r.With(middleware.RequirePermission("payroll.generate")).Post("/periods", handler.CreatePeriod)
	r.With(middleware.RequirePermission("payroll.generate")).Post("/periods/{periodID}/generate", handler.GenerateForPeriod)

	r.With(middleware.RequirePermission("payroll.adjustments.manage")).Post("/entries/{entryID}/adjustments", handler.CreateAdjustment)
	r.With(middleware.RequirePermission("payroll.send_to_teacher")).Post("/entries/{entryID}/send-to-teacher", handler.SendToTeacher)

	r.With(middleware.RequirePermission("payroll.confirm")).Post("/entries/{entryID}/confirm", handler.ConfirmByTeacher)
	r.With(middleware.RequirePermission("payroll.dispute")).Post("/entries/{entryID}/dispute", handler.DisputeByTeacher)

	r.With(middleware.RequirePermission("payroll.approve")).Post("/entries/{entryID}/approve", handler.ApproveByFinance)
	r.With(middleware.RequirePermission("payroll.mark_paid")).Post("/entries/{entryID}/mark-paid", handler.MarkPaid)
}
