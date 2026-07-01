package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Rumm1/eduhub-backend/internal/middleware"
	authmodule "github.com/Rumm1/eduhub-backend/internal/modules/auth"
	organizationmodule "github.com/Rumm1/eduhub-backend/internal/modules/organization"
	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, jwtManager *platformjwt.Manager) http.Handler {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	authRepository := authmodule.NewRepository(db)
	authService := authmodule.NewService(authRepository, jwtManager)
	authHandler := authmodule.NewHandler(authService)

	organizationRepository := organizationmodule.NewRepository(db)
	organizationService := organizationmodule.NewService(organizationRepository)
	organizationHandler := organizationmodule.NewHandler(organizationService)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		response.Message(w, http.StatusOK, "EduHub backend is running")
	})

	r.Get("/health/db", func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			response.Error(w, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "PostgreSQL is not available")
			return
		}

		response.Message(w, http.StatusOK, "PostgreSQL is connected")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			response.Message(w, http.StatusOK, "EduHub API is running")
		})

		r.Get("/health/db", func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
			defer cancel()

			if err := db.Ping(ctx); err != nil {
				response.Error(w, http.StatusServiceUnavailable, "DATABASE_UNAVAILABLE", "PostgreSQL is not available")
				return
			}

			response.Message(w, http.StatusOK, "PostgreSQL is connected")
		})

		r.Route("/auth", func(r chi.Router) {
			authmodule.RegisterRoutes(r, authHandler, jwtManager)
		})

		r.Route("/platform", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireRole("SUPER_ADMIN"))

			organizationmodule.RegisterPlatformRoutes(r, organizationHandler)
		})
	})

	return r
}
