// Package routers implements the routing for the web application using the Gin framework.
package routers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/controllers"
	"github.com/Glorified-Toaster/senior-project/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
	ginprometheus "github.com/zsais/go-gin-prometheus"
)

type Router struct {
	router         *gin.Engine
	controllers    *controllers.Controllers
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(ctrl *controllers.Controllers, authMiddleware *middleware.AuthMiddleware) *Router {
	// gin.SetMode(gin.ReleaseMode)

	// useing gin.Default() to create a router with default middleware: logger and recovery (crash-free) middleware
	router := gin.Default()

	store := cookie.NewStore([]byte(ctrl.ViperConfig.GinSession.Secret))
	router.Use(sessions.Sessions("session_token", store))

	// config and enable CORS Middleware
	location := fmt.Sprintf("https://%s:%s", ctrl.ViperConfig.HTTPServer.Addr, ctrl.ViperConfig.HTTPServer.Port)

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

	router.Use(csrf.Middleware(csrf.Options{
		Secret: ctrl.ViperConfig.CSRF.Secret,
		ErrorFunc: func(c *gin.Context) {
			c.String(400, "CSRF token mismatch")
			c.Abort()
		},
	}))

	// setting security headers
	router.Use(func(ctx *gin.Context) {
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-XSS-Protection", "1; mode=block")
		ctx.Header("Referrer-Policy", "strict-origin-when-cross-origin")
		ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		ctx.Next()
	})

	router.Static("/web/static", "./web/static")
	router.Static("/images", "./web/static/images")

	// using prometheus middleware
	prometheus := ginprometheus.NewPrometheus("gin")
	prometheus.Use(router)

	return &Router{
		router:         router,
		controllers:    ctrl,
		authMiddleware: authMiddleware,
	}
}

func (r *Router) GetHandler() http.Handler {
	return r.router
}

func (r *Router) SetupRoutes() {
	public := r.router.Group("/")
	{
		public.GET("/login", r.controllers.StudentLoginPageRender())
		public.GET("/instructor-login", r.controllers.InstructorLoginPageRender())
	}

	publicAPI := r.router.Group("/api/v1")
	{
		publicAPI.POST("/login", r.controllers.StudentLogin())
		// publicAPI.POST("/instructor-login", r.controllers.InstructorLogin())
		publicAPI.POST("/signup", r.controllers.Signup())
	}

	protected := r.router.Group("/api/v1")
	protected.Use(r.authMiddleware.AuthenticationMiddleware())
	{
		protected.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(200, gin.H{"msg": "PONG"})
		})
	}

	adminRoute := protected.Group("/admin")
	adminRoute.Use(r.authMiddleware.RequireRoles("admin"))
	{
		adminRoute.GET("/student/:id", r.controllers.GetStudentByID())
	}
}

func (r *Router) SetCORSConfig() *cors.Config {
	location := fmt.Sprintf("https://%s:%s", r.controllers.ViperConfig.HTTPServer.Addr, r.controllers.ViperConfig.HTTPServer.Port)
	return &cors.Config{
		AllowOrigins:     []string{location},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "X-CSRF-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
		AllowOriginFunc: func(origin string) bool {
			return origin == location
		},
	}
}
