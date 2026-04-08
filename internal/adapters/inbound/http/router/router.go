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
	"github.com/thinkerou/favicon"
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
	//enableCSRF(router, viperConfig)

	// setting security headers
	setSecurityHeaders(router)

	// static file handling
	router.Static("/src", "./web/static/src")
	router.Static("/web/static/src", "./web/static/src")
	router.Static("/images", "./web/static/images")
	router.Static("/static", "./web/static")
	router.LoadHTMLGlob("web/static/*.html")

	router.Use(favicon.New("./web/static/images/favicon.ico"))

	return &Router{
		router:         router,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes() {

	// Public API
	publicAPI := r.router.Group("/admin/api/v1")
	{
		publicAPI.POST("/admin/login", r.userHandler.Login())
	}

	// Public routes
	publicRoutes := r.router.Group("/")
	{
		publicRoutes.POST("/login", r.userHandler.Login())
		publicRoutes.GET("/test", r.userHandler.TestPage())
		publicRoutes.GET("/admin/login", r.userHandler.AdminLogin())
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
			dashboardRoutes.POST("/users/create", r.userHandler.Create())
			dashboardRoutes.GET("/users/generate-username", r.userHandler.GenerateRandomUsername())
			dashboardRoutes.GET("/users/generate-password", r.userHandler.GenerateRandomPassword())
			dashboardRoutes.GET("/users/deleted", r.userHandler.DeletedUsersPageRender())
			dashboardRoutes.GET("/users/deleted/search", r.userHandler.SearchDeletedUsers())
			dashboardRoutes.POST("/users/deleted/search", r.userHandler.SearchDeletedUsers())
			dashboardRoutes.GET("/exams", r.userHandler.AllExamsPageRender())
			dashboardRoutes.GET("/exams/search", r.userHandler.SearchExams())
			dashboardRoutes.POST("/exams/search", r.userHandler.SearchExams())
			dashboardRoutes.GET("/subjects", r.userHandler.AllSubjectsPageRender())
			dashboardRoutes.GET("/subjects/search", r.userHandler.SearchSubjects())
			dashboardRoutes.POST("/subjects/search", r.userHandler.SearchSubjects())
			dashboardRoutes.POST("/subjects/create", r.userHandler.CreateSubject())
			dashboardRoutes.GET("/logout", r.userHandler.Logout())
			dashboardRoutes.GET("/subject/:id", r.userHandler.EditSubjectPageRender())
			dashboardRoutes.POST("/subject/edit/:id", r.userHandler.EditSubjectInfo())
			dashboardRoutes.POST("/subject/delete/:id", r.userHandler.DeleteSubject())
			dashboardRoutes.GET("/subject/:id/exams/search", r.userHandler.SearchExamsBySubject())
			dashboardRoutes.POST("/subject/:id/exams/search", r.userHandler.SearchExamsBySubject())
			dashboardRoutes.POST("/exams/create", r.userHandler.CreateExam())
			dashboardRoutes.GET("/exam/:id", r.userHandler.EditExamPageRender())
			dashboardRoutes.POST("/exam/edit/:id", r.userHandler.EditExamInfo())
			dashboardRoutes.POST("/exam/delete/:id", r.userHandler.SoftDeleteExam())
			dashboardRoutes.POST("/exam/preview/question-text", r.userHandler.PreviewQuestionText())
			dashboardRoutes.POST("/exam/preview/question-choice", r.userHandler.PreviewQuestionChoice())
			dashboardRoutes.POST("/exam/preview/question-form", r.userHandler.GetQuestionForm())
			dashboardRoutes.POST("/exam/:id/question/create", r.userHandler.CreateQuestion())
			dashboardRoutes.POST("/exam/:id/question/upload-csv", r.userHandler.UploadQuestionCSV())
			dashboardRoutes.DELETE("/exam/:id/question/delete/:question-id", r.userHandler.DeleteQuestion())
			dashboardRoutes.GET("/exam/:id/export-csv", r.userHandler.ExportExamCSV())
			dashboardRoutes.POST("/questions/edit/:id", r.userHandler.UpdateQuestion())
		}
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
		ExposeHeaders:    []string{"Content-Length"},
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
