package organization

import "github.com/go-chi/chi/v5"

func RegisterPlatformRoutes(r chi.Router, handler *Handler) {
	r.Post("/organizations", handler.Create)
}
