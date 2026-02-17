package router

import (
	"fmt"
	"net/http"
	"time"

	"uot-exam/internal/adapters/inbound/http/handler"
	"uot-exam/internal/adapters/inbound/http/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	csrf "github.com/utrack/gin-csrf"
)

type Router struct {
	router         *gin.Engine
	handler        *handler.Handler
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(ctrl *handler.Handler, authMiddleware *middleware.AuthMiddleware) *Router {
	// useing gin.Default() to create a router with default middleware: logger and recovery (crash-free) middleware
	router := gin.Default()

	store := cookie.NewStore([]byte(ctrl.ViperConfig.GinSession.Secret))
	router.Use(sessions.Sessions("session_token", store))

	// config and enable CORS Middleware
	enableCORS(router, *ctrl)
	enableCSRF(router, *ctrl)

	// setting security headers
	setSecurityHeaders(router)

	// static file handling
	router.Static("/web/static", "./web/static")
	router.Static("/images", "./web/static/images")

	return &Router{
		router:         router,
		handler:        ctrl,
		authMiddleware: authMiddleware,
	}
}

func (router *Router) SetupRoutes() {}

func (r *Router) GetHandler() http.Handler {
	return r.router
}

func enableCORS(router *gin.Engine, ctrl handler.Handler) {
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
}

func enableCSRF(router *gin.Engine, ctrl handler.Handler) {
	router.Use(csrf.Middleware(csrf.Options{
		Secret: ctrl.ViperConfig.CSRF.Secret,
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
