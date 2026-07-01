package group

import (
	"github.com/Rumm1/eduhub-backend/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r chi.Router, handler *Handler) {
	r.With(middleware.RequirePermission("groups.read")).Get("/", handler.List)
	r.With(middleware.RequirePermission("groups.create")).Post("/", handler.Create)

	r.With(middleware.RequirePermission("groups.read")).Get("/{groupID}/students", handler.ListStudents)
	r.With(middleware.RequirePermission("groups.update")).Post("/{groupID}/students", handler.AddStudent)
}
