package profile

import (
"github.com/Rumm1/eduhub-backend/internal/middleware"
"github.com/go-chi/chi/v5"
)

func RegisterUserProfileRoutes(r chi.Router, handler *Handler) {
r.With(middleware.RequirePermission("profiles.read")).Get("/{userID}/profiles", handler.ListByUserID)
r.With(middleware.RequirePermission("profiles.manage")).Post("/{userID}/profiles", handler.Create)
}

func RegisterRoutes(r chi.Router, handler *Handler) {
r.With(middleware.RequirePermission("profiles.read")).Get("/{profileID}", handler.GetByID)

r.With(middleware.RequirePermission("profiles.manage")).Patch("/{profileID}", handler.Update)
r.With(middleware.RequirePermission("profiles.manage")).Delete("/{profileID}", handler.Disable)

r.With(middleware.RequirePermission("profiles.manage")).Post("/{profileID}/set-default", handler.SetDefault)

r.With(middleware.RequirePermission("profiles.manage")).Post("/{profileID}/roles", handler.AddRole)
r.With(middleware.RequirePermission("profiles.manage")).Delete("/{profileID}/roles/{roleCode}", handler.RemoveRole)

r.With(middleware.RequirePermission("profiles.manage")).Post("/{profileID}/branches", handler.AddBranch)
r.With(middleware.RequirePermission("profiles.manage")).Delete("/{profileID}/branches/{branchID}", handler.RemoveBranch)
}
