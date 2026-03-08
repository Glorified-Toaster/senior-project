package router

import (
	"fmt"
	"net/http"
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
)

type Router struct {
	router         *gin.Engine
	userHandler    *handler.UserHandler
	authMiddleware *middleware.AuthMiddleware
	viperConfig    *config.Config
}

func NewRouter(userHandler *handler.UserHandler, authMiddleware *middleware.AuthMiddleware, viperConfig *config.Config) *Router {
	// useing gin.Default() to create a router with default middleware: logger and recovery (crash-free) middleware
	router := gin.Default()

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
		publicRoutes.POST("/users/search", r.userHandler.SearchUsers())
		publicRoutes.POST("/users/new", r.userHandler.Create())
		publicRoutes.GET("/admin/login", r.userHandler.AdminLogin())
	}

	adminRoutes := r.router.Group("/admin")
	adminRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	adminRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleAdmin))
	{
		dashboardRoutes := adminRoutes.Group("/dashboard")
		{
			dashboardRoutes.GET("/", r.userHandler.AdminDashboardMainRender())
			dashboardRoutes.GET("/users", r.userHandler.UserPageRender())
		}
	}

	// User routes
	userRoutes := r.router.Group("/users")
	//userRoutes.Use(r.authMiddleware.AuthenticationMiddleware())
	//userRoutes.Use(r.authMiddleware.RoleAuthMiddleware(domain.RoleAdmin))
	{
		userRoutes.POST("/create", r.userHandler.Create())
		userRoutes.GET("/id/:id", r.userHandler.GetUserByID())
		userRoutes.GET("/username/:username", r.userHandler.GetUserByUsername())
		userRoutes.GET("/list-all", r.userHandler.ListAllUsers())
		userRoutes.DELETE("/delete/:id", r.userHandler.SoftDeleteUser())
	}
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
