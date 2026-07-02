package ai

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("dashboard.overview.read"))

		r.Get("/insights/dashboard", handler.DashboardInsights)
	})
}
