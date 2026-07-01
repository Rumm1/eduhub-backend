package payment

import (
"github.com/Rumm1/eduhub-backend/internal/middleware"
"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
r.With(middleware.RequirePermission("payments.read")).Get("/", handler.List)
r.With(middleware.RequirePermission("payments.manage")).Post("/", handler.Create)

r.With(middleware.RequirePermission("payments.read")).Get("/students/{studentID}", handler.ListByStudentID)
}
