package branding

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.Get("/current", handler.Current)

	r.Patch("/me/avatar", handler.UpdateMyAvatar)
	r.Delete("/me/avatar", handler.ClearMyAvatar)

	r.Group(func(r chi.Router) {
		r.Use(middleware.RequirePermission("files.upload"))

		r.Patch("/organization/logo", handler.UpdateOrganizationLogo)
		r.Delete("/organization/logo", handler.ClearOrganizationLogo)
	})
}
