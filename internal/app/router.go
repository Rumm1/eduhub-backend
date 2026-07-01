package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Rumm1/eduhub-backend/internal/middleware"
	attendancemodule "github.com/Rumm1/eduhub-backend/internal/modules/attendance"
	authmodule "github.com/Rumm1/eduhub-backend/internal/modules/auth"
	branchmodule "github.com/Rumm1/eduhub-backend/internal/modules/branch"
	groupmodule "github.com/Rumm1/eduhub-backend/internal/modules/group"
	homeworkmodule "github.com/Rumm1/eduhub-backend/internal/modules/homework"
	lessonmodule "github.com/Rumm1/eduhub-backend/internal/modules/lesson"
	organizationmodule "github.com/Rumm1/eduhub-backend/internal/modules/organization"
	studentmodule "github.com/Rumm1/eduhub-backend/internal/modules/student"
	subjectmodule "github.com/Rumm1/eduhub-backend/internal/modules/subject"
	teachermodule "github.com/Rumm1/eduhub-backend/internal/modules/teacher"
	usermodule "github.com/Rumm1/eduhub-backend/internal/modules/user"
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

	branchRepository := branchmodule.NewRepository(db)
	branchService := branchmodule.NewService(branchRepository)
	branchHandler := branchmodule.NewHandler(branchService)

	userRepository := usermodule.NewRepository(db)
	userService := usermodule.NewService(userRepository)
	userHandler := usermodule.NewHandler(userService)

	subjectRepository := subjectmodule.NewRepository(db)
	subjectService := subjectmodule.NewService(subjectRepository)
	subjectHandler := subjectmodule.NewHandler(subjectService)

	teacherRepository := teachermodule.NewRepository(db)
	teacherService := teachermodule.NewService(teacherRepository)
	teacherHandler := teachermodule.NewHandler(teacherService)

	studentRepository := studentmodule.NewRepository(db)
	studentService := studentmodule.NewService(studentRepository)
	studentHandler := studentmodule.NewHandler(studentService)

	groupRepository := groupmodule.NewRepository(db)
	groupService := groupmodule.NewService(groupRepository)
	groupHandler := groupmodule.NewHandler(groupService)

	lessonRepository := lessonmodule.NewRepository(db)
	lessonService := lessonmodule.NewService(lessonRepository)
	lessonHandler := lessonmodule.NewHandler(lessonService)

	attendanceRepository := attendancemodule.NewRepository(db)
	attendanceService := attendancemodule.NewService(attendanceRepository)
	attendanceHandler := attendancemodule.NewHandler(attendanceService)

	homeworkRepository := homeworkmodule.NewRepository(db)
	homeworkService := homeworkmodule.NewService(homeworkRepository)
	homeworkHandler := homeworkmodule.NewHandler(homeworkService)

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

		r.Route("/branches", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			branchmodule.RegisterRoutes(r, branchHandler)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			usermodule.RegisterRoutes(r, userHandler)
		})

		r.Route("/subjects", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			subjectmodule.RegisterRoutes(r, subjectHandler)
		})

		r.Route("/teachers", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			teachermodule.RegisterRoutes(r, teacherHandler)
		})

		r.Route("/students", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			studentmodule.RegisterRoutes(r, studentHandler)
		})

		r.Route("/groups", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			groupmodule.RegisterRoutes(r, groupHandler)
		})

		r.Route("/lessons", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			lessonmodule.RegisterRoutes(r, lessonHandler)
		})

		r.Route("/attendance", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			attendancemodule.RegisterRoutes(r, attendanceHandler)
		})

		r.Route("/homeworks", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			homeworkmodule.RegisterRoutes(r, homeworkHandler)
		})
	})

	return r
}
