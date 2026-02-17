package handler

import (
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/adapters/outbound/logger"

	"github.com/go-playground/validator"
)

type Handler struct {
	validator *validator.Validate
	// UserRepo  	repository.UserRepository
	jwtAuth     *helpers.JWTAuth
	ViperConfig *config.Config
	logger      *logger.Logger
}

func NewHandler(valid *validator.Validate, jwt *helpers.JWTAuth, viperConfig *config.Config, logger *logger.Logger) *Handler {
	return &Handler{
		valid,
		jwt,
		viperConfig,
		logger,
	}
}
