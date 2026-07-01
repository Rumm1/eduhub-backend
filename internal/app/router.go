package app

import (
	"context"
	"net/http"
	"time"

	authmodule "github.com/Rumm1/eduhub-backend/internal/modules/auth"
	platformjwt "github.com/Rumm1/eduhub-backend/internal/platform/jwt"
	"github.com/Rumm1/eduhub-backend/internal/shared/response"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(db *pgxpool.Pool, jwtManager *platformjwt.Manager) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	authRepository := authmodule.NewRepository(db)
	authService := authmodule.NewService(authRepository, jwtManager)
	authHandler := authmodule.NewHandler(authService)

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
			authmodule.RegisterRoutes(r, authHandler)
		})
	})

	return r
}
