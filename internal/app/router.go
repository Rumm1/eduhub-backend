package app

import (
	"context"
	"net/http"
	"time"

	"github.com/Rumm1/eduhub-backend/internal/middleware"
	aimodule "github.com/Rumm1/eduhub-backend/internal/modules/ai"
	attendancemodule "github.com/Rumm1/eduhub-backend/internal/modules/attendance"
	auditmodule "github.com/Rumm1/eduhub-backend/internal/modules/audit"
	authmodule "github.com/Rumm1/eduhub-backend/internal/modules/auth"
	branchmodule "github.com/Rumm1/eduhub-backend/internal/modules/branch"
	brandingmodule "github.com/Rumm1/eduhub-backend/internal/modules/branding"
	dashboardmodule "github.com/Rumm1/eduhub-backend/internal/modules/dashboard"
	filemodule "github.com/Rumm1/eduhub-backend/internal/modules/file"
	groupmodule "github.com/Rumm1/eduhub-backend/internal/modules/group"
	homeworkmodule "github.com/Rumm1/eduhub-backend/internal/modules/homework"
	importermodule "github.com/Rumm1/eduhub-backend/internal/modules/importer"
	lessonmodule "github.com/Rumm1/eduhub-backend/internal/modules/lesson"
	notificationmodule "github.com/Rumm1/eduhub-backend/internal/modules/notification"
	organizationmodule "github.com/Rumm1/eduhub-backend/internal/modules/organization"
	parentmodule "github.com/Rumm1/eduhub-backend/internal/modules/parent"
	paymentmodule "github.com/Rumm1/eduhub-backend/internal/modules/payment"
	payrollmodule "github.com/Rumm1/eduhub-backend/internal/modules/payroll"
	permissionmodule "github.com/Rumm1/eduhub-backend/internal/modules/permission"
	platformdashboardmodule "github.com/Rumm1/eduhub-backend/internal/modules/platformdashboard"
	platformusermodule "github.com/Rumm1/eduhub-backend/internal/modules/platformuser"
	profilemodule "github.com/Rumm1/eduhub-backend/internal/modules/profile"
	reportmodule "github.com/Rumm1/eduhub-backend/internal/modules/report"
	rolemodule "github.com/Rumm1/eduhub-backend/internal/modules/role"
	schedulemodule "github.com/Rumm1/eduhub-backend/internal/modules/schedule"
	studentmodule "github.com/Rumm1/eduhub-backend/internal/modules/student"
	subjectmodule "github.com/Rumm1/eduhub-backend/internal/modules/subject"
	swaggermodule "github.com/Rumm1/eduhub-backend/internal/modules/swagger"
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

	r.Use(middleware.CORS)

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	aiRepository := aimodule.NewRepository(db)
	aiService := aimodule.NewService(aiRepository)
	aiHandler := aimodule.NewHandler(aiService)

	authRepository := authmodule.NewRepository(db)
	authService := authmodule.NewService(authRepository, jwtManager)
	authHandler := authmodule.NewHandler(authService)

	auditRepository := auditmodule.NewRepository(db)
	auditService := auditmodule.NewService(auditRepository)

	notificationRepository := notificationmodule.NewRepository(db)
	notificationService := notificationmodule.NewService(notificationRepository)
	notificationHandler := notificationmodule.NewHandler(notificationService)

	organizationRepository := organizationmodule.NewRepository(db)
	organizationService := organizationmodule.NewService(organizationRepository)
	organizationHandler := organizationmodule.NewHandler(organizationService)

	brandingRepository := brandingmodule.NewRepository(db)
	brandingService := brandingmodule.NewService(brandingRepository)
	brandingHandler := brandingmodule.NewHandler(brandingService)

	branchRepository := branchmodule.NewRepository(db)
	branchService := branchmodule.NewService(branchRepository)
	branchHandler := branchmodule.NewHandler(branchService)

	userRepository := usermodule.NewRepository(db)
	userService := usermodule.NewService(userRepository)
	userHandler := usermodule.NewHandler(userService)

	profileRepository := profilemodule.NewRepository(db)
	profileService := profilemodule.NewService(profileRepository)
	profileHandler := profilemodule.NewHandler(profileService, auditService)

	roleRepository := rolemodule.NewRepository(db)
	roleService := rolemodule.NewService(roleRepository)
	roleHandler := rolemodule.NewHandler(roleService, auditService)

	permissionRepository := permissionmodule.NewRepository(db)
	permissionService := permissionmodule.NewService(permissionRepository)
	permissionHandler := permissionmodule.NewHandler(permissionService)

	subjectRepository := subjectmodule.NewRepository(db)
	subjectService := subjectmodule.NewService(subjectRepository)
	subjectHandler := subjectmodule.NewHandler(subjectService)

	teacherRepository := teachermodule.NewRepository(db)
	teacherService := teachermodule.NewService(teacherRepository)
	teacherHandler := teachermodule.NewHandler(teacherService)

	studentRepository := studentmodule.NewRepository(db)
	studentService := studentmodule.NewService(studentRepository)
	studentHandler := studentmodule.NewHandler(studentService)

	dashboardRepository := dashboardmodule.NewRepository(db)
	dashboardService := dashboardmodule.NewService(dashboardRepository)
	dashboardHandler := dashboardmodule.NewHandler(dashboardService)

	fileRepository := filemodule.NewRepository(db)
	fileService := filemodule.NewService(fileRepository)
	fileHandler := filemodule.NewHandler(fileService)

	groupRepository := groupmodule.NewRepository(db)
	groupService := groupmodule.NewService(groupRepository)
	groupHandler := groupmodule.NewHandler(groupService)

	lessonRepository := lessonmodule.NewRepository(db)
	lessonService := lessonmodule.NewService(lessonRepository)
	lessonHandler := lessonmodule.NewHandler(lessonService)

	auditHandler := auditmodule.NewHandler(auditService)

	attendanceRepository := attendancemodule.NewRepository(db)
	attendanceService := attendancemodule.NewService(attendanceRepository)
	attendanceHandler := attendancemodule.NewHandler(attendanceService)

	homeworkRepository := homeworkmodule.NewRepository(db)
	homeworkService := homeworkmodule.NewService(homeworkRepository)
	homeworkHandler := homeworkmodule.NewHandler(homeworkService)

	importerRepository := importermodule.NewRepository(db)
	importerService := importermodule.NewService(importerRepository)
	importerHandler := importermodule.NewHandler(importerService)

	parentRepository := parentmodule.NewRepository(db)
	parentService := parentmodule.NewService(parentRepository)
	parentHandler := parentmodule.NewHandler(parentService)

	platformUserRepository := platformusermodule.NewRepository(db)
	platformUserService := platformusermodule.NewService(platformUserRepository)
	platformUserHandler := platformusermodule.NewHandler(platformUserService)

	platformDashboardRepository := platformdashboardmodule.NewRepository(db)
	platformDashboardService := platformdashboardmodule.NewService(platformDashboardRepository)
	platformDashboardHandler := platformdashboardmodule.NewHandler(platformDashboardService)

	paymentRepository := paymentmodule.NewRepository(db)
	paymentService := paymentmodule.NewService(paymentRepository)
	paymentHandler := paymentmodule.NewHandler(paymentService)

	payrollRepository := payrollmodule.NewRepository(db)
	payrollService := payrollmodule.NewService(payrollRepository)
	payrollHandler := payrollmodule.NewHandler(payrollService, auditService)

	reportRepository := reportmodule.NewRepository(db)
	reportService := reportmodule.NewService(reportRepository)
	reportHandler := reportmodule.NewHandler(reportService)

	scheduleRepository := schedulemodule.NewRepository(db)
	scheduleService := schedulemodule.NewService(scheduleRepository)
	scheduleHandler := schedulemodule.NewHandler(scheduleService)

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

	r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	r.Route("/swagger", func(r chi.Router) {
		swaggermodule.RegisterRoutes(r)
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

		r.Route("/ai", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			aimodule.RegisterRoutes(r, aiHandler)
		})

		r.Route("/auth", func(r chi.Router) {
			authmodule.RegisterRoutes(r, authHandler, jwtManager)
		})

		r.Route("/platform", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireRole("SUPER_ADMIN"))
			platformdashboardmodule.RegisterRoutes(r, platformDashboardHandler)
			platformusermodule.RegisterRoutes(r, platformUserHandler)
			organizationmodule.RegisterPlatformRoutes(r, organizationHandler)
		})

		r.Route("/branding", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			brandingmodule.RegisterRoutes(r, brandingHandler)
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
			profilemodule.RegisterUserProfileRoutes(r, profileHandler)
		})

		r.Route("/profiles", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			r.Use(profilemodule.ProfileAuditMiddleware(auditService))

			profilemodule.RegisterRoutes(r, profileHandler)
		})

		r.Route("/roles", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			r.Use(rolemodule.RoleAuditMiddleware(auditService))

			rolemodule.RegisterRoutes(r, roleHandler)
		})

		r.Route("/permissions", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			permissionmodule.RegisterRoutes(r, permissionHandler)
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

		r.Route("/dashboard", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			dashboardmodule.RegisterRoutes(r, dashboardHandler)
		})

		r.Route("/files", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			filemodule.RegisterRoutes(r, fileHandler)
		})

		r.Route("/groups", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			groupmodule.RegisterRoutes(r, groupHandler)
		})

		r.Route("/lessons", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			r.Use(lessonmodule.LessonAuditMiddleware(auditService))

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

		r.Route("/schedules", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			schedulemodule.RegisterRoutes(r, scheduleHandler)
		})

		r.Route("/imports", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			importermodule.RegisterRoutes(r, importerHandler)
		})

		r.Route("/notifications", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			notificationmodule.RegisterRoutes(r, notificationHandler)
		})

		r.Route("/parents", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			parentmodule.RegisterRoutes(r, parentHandler)
		})

		r.Route("/payments", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			r.Use(paymentmodule.PaymentAuditMiddleware(auditService))

			paymentmodule.RegisterRoutes(r, paymentHandler)
		})

		r.Route("/payroll", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			r.Use(payrollmodule.PayrollAuditMiddleware(auditService))

			payrollmodule.RegisterRoutes(r, payrollHandler)
		})
		r.Route("/reports", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			r.Use(reportmodule.ReportAuditMiddleware(auditService))

			reportmodule.RegisterRoutes(r, reportHandler)
		})
		r.Route("/audit-logs", func(r chi.Router) {
			r.Use(middleware.Auth(jwtManager))
			r.Use(middleware.RequireTenant)

			auditmodule.RegisterRoutes(r, auditHandler)
		})
	})

	return r
}
