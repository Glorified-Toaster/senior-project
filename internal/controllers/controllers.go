// Package controllers contains handler functions for the web application.
package controllers

import (
	"github.com/Glorified-Toaster/senior-project/internal/config"
	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/helpers"
	"github.com/Glorified-Toaster/senior-project/internal/repository"
	"github.com/go-playground/validator"
)

type Controllers struct {
	validator   *validator.Validate
	UserRepo    repository.UserRepository
	cache       cache.Cache
	jwtAuth     *helpers.JWTAuth
	ViperConfig *config.Config
}

func NewControllers(valid *validator.Validate, userRepo repository.UserRepository, cache cache.Cache, jwt *helpers.JWTAuth, viperConfig *config.Config) *Controllers {
	return &Controllers{
		valid,
		userRepo,
		cache,
		jwt,
		viperConfig,
	}
}
