package permission

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("profiles.manage")).Get("/", handler.List)
	r.With(middleware.RequirePermission("profiles.manage")).Get("/groups", handler.ListGroups)
}
