// Package controllers contains handler functions for the web application.
package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/models"
	"github.com/Glorified-Toaster/senior-project/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

type Controllers struct {
	validator *validator.Validate
	repo      repository.UserRepository
	cache     cache.Cache
}

var ctrl *Controllers

func NewControllers(valid *validator.Validate, repo repository.UserRepository, cache cache.Cache) {
	ctrl = &Controllers{
		valid,
		repo,
		cache,
	}
}

func Ping(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "pong",
	})
}

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var user models.User

		// Get user input
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Validate user input
		if validationErr := ctrl.validator.Struct(user); validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		_, err := ctrl.repo.CreateUser(c, &user)
		if err != nil {
			return
		}
	}
}
