package router

import (
	"fmt"
	"net/http"
	"time"

	"uot-exam/internal/adapters/inbound/http/handler"
	"uot-exam/internal/adapters/inbound/http/middleware"
	"uot-exam/internal/adapters/outbound/config"

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
	router.Static("/web/static", "./web/static")
	router.Static("/images", "./web/static/images")

	return &Router{
		router:         router,
		userHandler:    userHandler,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) SetupRoutes() {
	// User routes
	userRoutes := r.router.Group("/users")
	{
		userRoutes.POST("/create", r.userHandler.Create())
		userRoutes.GET("/id/:id", r.userHandler.GetUserByID())
		userRoutes.GET("/username/:username", r.userHandler.GetUserByUsername())
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
