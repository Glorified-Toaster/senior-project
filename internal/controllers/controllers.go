// Package controllers contains handler functions for the web application.
package controllers

import (
	"github.com/Glorified-Toaster/senior-project/internal/config/db/cache"
	"github.com/Glorified-Toaster/senior-project/internal/repository"
	"github.com/go-playground/validator"
)

type Controllers struct {
	validator   *validator.Validate
	StudentRepo repository.StudentRepository
	cache       cache.Cache
}

func NewControllers(valid *validator.Validate, studentRepo repository.StudentRepository, cache cache.Cache) *Controllers {
	return &Controllers{
		valid,
		studentRepo,
		cache,
	}
}
