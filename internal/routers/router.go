// Package routers implements the routing for the web application using the Gin framework.
package routers

import (
	"net/http"

	"github.com/Glorified-Toaster/senior-project/internal/controllers"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
}

func NewRouter() *Router {
	// gin.SetMode(gin.ReleaseMode)

	// use gin.Default() to create a router with default middleware: logger and recovery (crash-free) middleware
	router := gin.Default()

	return &Router{router: router}
}

func (r *Router) GetHandler() http.Handler {
	return r.router
}

func (r *Router) SetupRoutes() {
	r.router.GET("/ping", controllers.Ping)
	r.router.POST("/signup", controllers.Signup())
}
