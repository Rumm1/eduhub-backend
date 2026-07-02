package role

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("roles.read")).Get("/", handler.List)
	r.With(middleware.RequirePermission("roles.manage")).Post("/", handler.Create)

	r.With(middleware.RequirePermission("roles.read")).Get("/{roleID}", handler.GetByID)
	r.With(middleware.RequirePermission("roles.manage")).Patch("/{roleID}", handler.Update)
	r.With(middleware.RequirePermission("roles.manage")).Delete("/{roleID}", handler.Delete)

	r.With(middleware.RequirePermission("roles.manage")).Post("/{roleID}/permissions", handler.AddPermission)
	r.With(middleware.RequirePermission("roles.manage")).Delete("/{roleID}/permissions/{permissionCode}", handler.RemovePermission)
}
