// Package routers implements the routing for the web application using the Gin framework.
package routers

import (
	"net/http"

	"github.com/Glorified-Toaster/senior-project/internal/controllers"
	"github.com/Glorified-Toaster/senior-project/internal/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router         *gin.Engine
	controllers    *controllers.Controllers
	authMiddleware *middleware.AuthMiddleware
}

func NewRouter(ctrl *controllers.Controllers, authMiddleware *middleware.AuthMiddleware) *Router {
	// gin.SetMode(gin.ReleaseMode)

	// use gin.Default() to create a router with default middleware: logger and recovery (crash-free) middleware
	router := gin.Default()

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
	public := r.router.Group("/api/v1")

	{
		public.GET("/ping", controllers.Ping)
		public.POST("/login", r.controllers.StudentLogin())
		public.POST("/signup", r.controllers.Signup())
	}

	protected := r.router.Group("/api/v1")
	protected.Use(r.authMiddleware.AuthenticationMiddleware())
	{
		protected.GET("/student/:id", r.controllers.GetStudentByID())
	}
}
