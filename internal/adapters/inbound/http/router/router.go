package router

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"uot-exam/internal/adapters/inbound/http/handler"
	"uot-exam/internal/adapters/inbound/http/middleware"
	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/domain"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Router struct {
	router         *gin.Engine
	userHandler    *handler.UserHandler
	authMiddleware *middleware.AuthMiddleware
	viperConfig    *config.Config
}

func NewRouter(userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware, viperConfig *config.Config) *Router {

	gin.DisableConsoleColor()
	// Logging to a file.
	gin.DefaultWriter = io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   viperConfig.GinLogger.Filename,
		MaxSize:    viperConfig.Lumberjack.MaxSize,
		MaxBackups: viperConfig.Lumberjack.MaxBackups,
		MaxAge:     viperConfig.Lumberjack.MaxAge,
		Compress:   viperConfig.Lumberjack.Compress,
	})

	// useing gin.Default() to create a router with default middleware: logger and recovery (crash-free) middleware
	router := gin.Default()

	router.MaxMultipartMemory = 8 << 20 // 8MB

	store := cookie.NewStore([]byte(viperConfig.GinSession.Secret))
	router.Use(sessions.Sessions("session_token", store))

	// config and enable CORS Middleware
	enableCORS(router, viperConfig)
	enableCSRF(router, viperConfig)

	// setting security headers
	setSecurityHeaders(router)

	// static file handling
	router.Static("/src", "./web/static/src")
	router.Static("/web/static/src", "./web/static/src")
	router.Static("/images", "./web/static/images")
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/static/*.html")

	router.StaticFile("/favicon.ico", "./web/static/images/favicon.ico")

	return &Router{
		router:         router,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes() {

	// Public API
	publicAPI := r.router.Group("/api/v1")
	{
		publicAPI.POST("/login", r.userHandler.Login())

	}

	// Public routes
	publicRoutes := r.router.Group("/")
	{
		publicRoutes.POST("/login", r.userHandler.Login())
		publicRoutes.GET("/login", r.userHandler.AdminLogin())
		publicRoutes.GET("/logout", r.userHandler.Logout())
		publicRoutes.GET("/", r.userHandler.LandingPage())
		publicRoutes.POST("/setup/admin", r.userHandler.SetupAdmin())
	}

	adminRoutes := r.router.Group("/admin")
	adminRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	adminRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleAdmin))
	{
		adminRoutes.GET("/users/search", r.userHandler.SearchUsers())
		adminRoutes.POST("/users/search", r.userHandler.SearchUsers())
		dashboardRoutes := adminRoutes.Group("/dashboard")
		{
			dashboardRoutes.GET("/", r.userHandler.AdminDashboardMainRender())
			dashboardRoutes.GET("/users", r.userHandler.UserPageRender())
			dashboardRoutes.GET("/user/:id", r.userHandler.EditUserPageRender())
			dashboardRoutes.POST("/user/edit/:id", r.userHandler.EditUserInfo())
			dashboardRoutes.POST("/users/create", r.userHandler.Create())
			dashboardRoutes.POST("/users/create-csv", r.userHandler.CreateUserCSV())
			dashboardRoutes.GET("/users/generate-username", r.userHandler.GenerateRandomUsername())
			dashboardRoutes.GET("/users/generate-password", r.userHandler.GenerateRandomPassword())
			dashboardRoutes.GET("/users/deleted/search", r.userHandler.SearchDeletedUsers())
			dashboardRoutes.POST("/users/deleted/search", r.userHandler.SearchDeletedUsers())
			dashboardRoutes.GET("/exams", r.userHandler.AllExamsPageRender())
			dashboardRoutes.GET("/exams/search", r.userHandler.SearchExams())
			dashboardRoutes.POST("/exams/search", r.userHandler.SearchExams())
			dashboardRoutes.GET("/subjects", r.userHandler.AllSubjectsPageRender())
			dashboardRoutes.GET("/subjects/search", r.userHandler.SearchSubjects())
			dashboardRoutes.POST("/subjects/search", r.userHandler.SearchSubjects())
			dashboardRoutes.POST("/subjects/create", r.userHandler.CreateSubject())
			dashboardRoutes.GET("/subject/:id", r.userHandler.EditSubjectPageRender())
			dashboardRoutes.POST("/subject/edit/:id", r.userHandler.EditSubjectInfo())
			dashboardRoutes.POST("/subject/publish/:id", r.userHandler.PublishSubject())
			dashboardRoutes.POST("/subject/delete/:id", r.userHandler.DeleteSubject())
			dashboardRoutes.GET("/subject/:id/ws", r.userHandler.AdminSubjectTrackerWS())
			dashboardRoutes.GET("/subject/:id/export-pdf", r.userHandler.SubjectOverallAttemptsPDF())
			dashboardRoutes.POST("/subject/:id/timer/end", r.userHandler.AdminSubjectEndTimer())
			dashboardRoutes.POST("/subject/:id/timer/extend", r.userHandler.AdminSubjectExtendTimer())
			dashboardRoutes.POST("/subject/:id/instructors/assign", r.userHandler.AssignInstructorToSubject())
			dashboardRoutes.POST("/subject/:id/instructors/unassign/:user_id", r.userHandler.UnassignInstructorFromSubject())
			dashboardRoutes.GET("/subject/:id/instructors/search", r.userHandler.SearchSubjectInstructors())
			dashboardRoutes.POST("/subject/:id/instructors/search", r.userHandler.SearchSubjectInstructors())
			dashboardRoutes.POST("/subject/:id/students/assign", r.userHandler.AssignStudentToSubject())
			dashboardRoutes.GET("/subject/:id/students/search", r.userHandler.SearchSubjectStudents())
			dashboardRoutes.POST("/subject/:id/students/search", r.userHandler.SearchSubjectStudents())
			dashboardRoutes.POST("/subject/:id/students/assign-csv", r.userHandler.AssignStudentToSubjectCSV())
			dashboardRoutes.GET("/subject/:id/students/export", r.userHandler.ExportSubjectStudentsCSV())
			dashboardRoutes.POST("/subject/:id/students/unassign/:user_id", r.userHandler.UnassignStudentFromSubject())
			dashboardRoutes.GET("/subject/:id/exams/search", r.userHandler.SearchExamsBySubject())
			dashboardRoutes.POST("/subject/:id/exams/search", r.userHandler.SearchExamsBySubject())
			dashboardRoutes.POST("/exams/create", r.userHandler.CreateExam())
			dashboardRoutes.GET("/exam/:id", r.userHandler.EditExamPageRender())
			dashboardRoutes.POST("/exam/edit/:id", r.userHandler.EditExamInfo())
			dashboardRoutes.POST("/exam/delete/:id", r.userHandler.SoftDeleteExam())
			dashboardRoutes.GET("/exam/:id/export-csv", r.userHandler.ExportExamCSV())
		}
		adminRoutes.GET("/settings", r.userHandler.AdminSettingsPageRender())
		adminRoutes.POST("/settings/general", r.userHandler.UpdateSettings())
		adminRoutes.GET("/settings/backup", r.userHandler.ExportBackup())
	}

	// User routes
	userRoutes := r.router.Group("/users")
	userRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	userRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleAdmin))
	{
		userRoutes.POST("/create", r.userHandler.Create())
		userRoutes.GET("/id/:id", r.userHandler.GetUserByID())
		userRoutes.GET("/username/:username", r.userHandler.GetUserByUsername())
		userRoutes.GET("/list-all", r.userHandler.ListAllUsers())
		userRoutes.DELETE("/delete/:id", r.userHandler.SoftDeleteUser())
		userRoutes.POST("/toggle/:id", r.userHandler.ToggleUserActive())
		userRoutes.POST("/update-password/:id", r.userHandler.UpdateUserPassword())
		userRoutes.POST("/restore/:id", r.userHandler.RestoreUser())
	}

	// Student routes
	studentRoutes := r.router.Group("/student")
	studentRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	studentRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleStudent))
	{
		studentRoutes.GET("/dashboard", r.userHandler.StudentDashboardRender())
		
		// Subject-specific routes
		subjectRoutes := studentRoutes.Group("/subject/:id")
		subjectRoutes.Use(r.authMiddleware.SubjectAccessMiddleware())
		{
			subjectRoutes.GET("", r.userHandler.StudentSubjectView())
			subjectRoutes.GET("/ws", r.userHandler.StudentSubjectTrackerWS())
			subjectRoutes.POST("/auto-submit", r.userHandler.StudentSubjectAutoSubmit())
			subjectRoutes.GET("/result", r.userHandler.StudentSubjectResult())
			subjectRoutes.GET("/pdf", r.userHandler.StudentSubjectPDF())
		}

		// Exam-specific routes
		examRoutes := studentRoutes.Group("/exam/:id")
		examRoutes.Use(r.authMiddleware.ExamAccessMiddleware())
		{
			examRoutes.POST("/start", r.userHandler.StudentStartExam())
			examRoutes.GET("/take", r.userHandler.StudentExamView())
			examRoutes.POST("/answer", r.userHandler.StudentSaveAnswer())
			examRoutes.POST("/submit", r.userHandler.StudentSubmitExam())
			examRoutes.POST("/auto-submit", r.userHandler.StudentAutoSubmit())
		}
	}

	// Instructor routes
	instructorRoutes := r.router.Group("/instructor")
	instructorRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	instructorRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleInstructor))
	{
		instructorRoutes.GET("/dashboard", r.userHandler.InstructorDashboardRender())
		instructorRoutes.POST("/subjects/search", r.userHandler.SearchInstructorSubjects())

		// Subject-specific routes
		subjectRoutes := instructorRoutes.Group("/subject/:id")
		subjectRoutes.Use(r.authMiddleware.SubjectAccessMiddleware())
		{
			subjectRoutes.GET("", r.userHandler.InstructorSubjectView())
			subjectRoutes.GET("/students/search", r.userHandler.SearchSubjectStudentsInstructor())
			subjectRoutes.POST("/students/search", r.userHandler.SearchSubjectStudentsInstructor())
			subjectRoutes.POST("/exams/search", r.userHandler.SearchExamsBySubject())
			subjectRoutes.GET("/students/export", r.userHandler.ExportSubjectStudentsCSV())
			subjectRoutes.POST("/students/assign", r.userHandler.AssignStudentToSubjectInstructor())
			subjectRoutes.POST("/students/upload-csv", r.userHandler.AssignStudentToSubjectCSV())
			subjectRoutes.POST("/students/unassign/:user_id", r.userHandler.UnassignStudentFromSubjectInstructor())
			
			// Exam creation within a subject
			subjectRoutes.POST("/exams/create", r.userHandler.CreateExam())
		}

		// Exam-specific routes (Instructor restricted)
		examRoutes := instructorRoutes.Group("/exam/:id")
		examRoutes.Use(r.authMiddleware.ExamAccessMiddleware())
		{
			examRoutes.GET("", r.userHandler.InstructorEditExamPageRender())
			examRoutes.POST("/edit", r.userHandler.EditExamInfo())
			examRoutes.POST("/delete", r.userHandler.SoftDeleteExam())
			examRoutes.GET("/export-csv", r.userHandler.ExportExamCSV())

			// Question Management
			examRoutes.POST("/question/create", r.userHandler.CreateQuestion())
			examRoutes.POST("/question/upload-csv", r.userHandler.UploadQuestionCSV())
			examRoutes.POST("/question/upload-csv-random", r.userHandler.UploadQuestionCSVRandom())
			examRoutes.POST("/question/edit/:question-id", r.userHandler.UpdateQuestion())
			examRoutes.POST("/question/delete/:question-id", r.userHandler.DeleteQuestion())
		}
	}

	// Shared routes for exam management (Admin + Instructor)
	// We use the same path prefix as admin to avoid changing components for now
	sharedExamRoutes := r.router.Group("/admin/dashboard/exam")
	sharedExamRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	sharedExamRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleAdmin, domain.RoleInstructor))
	{
		sharedExamRoutes.POST("/preview/question-text", r.userHandler.PreviewQuestionText())
		sharedExamRoutes.POST("/preview/question-choice", r.userHandler.PreviewQuestionChoice())
		sharedExamRoutes.POST("/preview/question-form", r.userHandler.GetQuestionForm())
		sharedExamRoutes.POST("/:id/question/create", r.userHandler.CreateQuestion())
		sharedExamRoutes.POST("/:id/question/upload-csv", r.userHandler.UploadQuestionCSV())
		sharedExamRoutes.POST("/:id/question/upload-csv-random", r.userHandler.UploadQuestionCSVRandom())
		sharedExamRoutes.POST("/:id/question/edit/:question-id", r.userHandler.UpdateQuestion())
		sharedExamRoutes.DELETE("/:id/question/delete/:question-id", r.userHandler.DeleteQuestion())
		sharedExamRoutes.POST("/:id/question/delete/:question-id", r.userHandler.DeleteQuestion())
		sharedExamRoutes.GET("/:id/export-pdf", r.userHandler.ExamAttemptsPDF())
	}

	sharedQuestionRoutes := r.router.Group("/admin/dashboard/questions")
	sharedQuestionRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	sharedQuestionRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleAdmin, domain.RoleInstructor))
	{
		sharedQuestionRoutes.POST("/edit/:id", r.userHandler.UpdateQuestion())
	}

	r.router.NoRoute(func(c *gin.Context) {
		c.HTML(http.StatusNotFound, "404.html", nil)
	})
}

func (r *Router) GetHandler() http.Handler {
	return r.router
}

func enableCORS(router *gin.Engine, viperConfig *config.Config) {
	// config and enable CORS Middleware
	location := fmt.Sprintf("https://%s:%s", viperConfig.HTTPServer.Addr, viperConfig.HTTPServer.Port)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{location},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-CSRF-Token", "Authorization"},
		ExposeHeaders:    []string{},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return origin == location
		},
	}))
}

func enableCSRF(router *gin.Engine, viperConfig *config.Config) {
	router.Use(csrf.Middleware(csrf.Options{
		Secret: viperConfig.CSRF.Secret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))
}

func setSecurityHeaders(router *gin.Engine) {
	router.Use(func(ctx *gin.Context) {
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		ctx.Next()
	})
}
