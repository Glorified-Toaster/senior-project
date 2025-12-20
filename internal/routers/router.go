// Package routers implements the routing for the web application using the Gin framework.
package routers

import (
	"net/http"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/controllers"
	"github.com/Glorified-Toaster/senior-project/internal/middleware"
	"github.com/Glorified-Toaster/senior-project/internal/templates"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/gin-gonic/gin"
	ginprometheus "github.com/zsais/go-gin-prometheus"
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
		public.GET("/login", func(ctx *gin.Context) {
			render := utils.NewRender(ctx, http.StatusOK, templates.StudentLoginPage())
			ctx.Render(http.StatusOK, render)
		})
	}

	publicAPI := r.router.Group("/api/v1")
	{
		publicAPI.GET("/ping", controllers.Ping())
		publicAPI.POST("/login", r.controllers.StudentLogin())
		publicAPI.POST("/signup", r.controllers.Signup())
		publicAPI.GET("/simple-content", func(c *gin.Context) {
			currentTime := time.Now().Format("15:04:05")
			templates.SimpleContent(currentTime).Render(c.Request.Context(), c.Writer)
		})
	}

	protected := r.router.Group("/api/v1")
	protected.Use(r.authMiddleware.AuthenticationMiddleware())
	{
		protected.GET("/student/:id", r.controllers.GetStudentByID())
	}
}
