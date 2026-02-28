package handler

import (
	"net/http"
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/outbound/config"
	"uot-exam/internal/adapters/outbound/logger"
	"uot-exam/internal/application"
	"uot-exam/internal/ports"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type UserHandler struct {
	userApp     *application.Application
	validate    *validator.Validate
	jwt         *helpers.JWTAuth
	viperConfig *config.Config
	logger      *logger.Logger
}

func NewUserHandler(userApp *application.Application, validate *validator.Validate, jwt *helpers.JWTAuth, viperConfig *config.Config, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userApp:     userApp,
		validate:    validate,
		jwt:         jwt,
		viperConfig: viperConfig,
		logger:      logger,
	}
}

func (h *UserHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ports.CreateUserParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (json)"})
			return
		}

		if err := h.validate.Struct(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (struct validation)"})
			return
		}

		user, err := h.userApp.CreateUser(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user : " + err.Error()})
			return
		}

		token, err := h.jwt.GenerateToken(user)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":        "User created successfully. Please login to get access token.",
				"student_id": user.ID,
				"warning":    "Token generation failed - please login manually",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{
			"msg":          "User created successfully",
			"access_token": token,
			"token_type":   "Bearer",
			"user": gin.H{
				"id":        user.ID,
				"full_name": user.FullName,
				"user_name": user.Username,
				"is_active": user.IsActive,
				"role":      user.Role,
			},
		})
	}

}

func (h *UserHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ports.LoginParams
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (json)"})
			return
		}

		if err := h.validate.Struct(req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body (struct validation)"})
			return
		}

		user, err := h.userApp.Login(ctx, req)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to login user : " + err.Error()})
			return
		}

		token, err := h.jwt.GenerateToken(user)
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"msg":        "User logged in successfully. Please login to get access token.",
				"student_id": user.ID,
				"warning":    "Token generation failed - please login manually",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"msg":          "User logged in successfully",
			"access_token": token,
			"token_type":   "Bearer",
			"user": gin.H{
				"id":        user.ID,
				"full_name": user.FullName,
				"user_name": user.Username,
				"is_active": user.IsActive,
				"role":      user.Role,
			},
		})
	}

}

func (h *UserHandler) GetUserByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := uuid.Parse(ctx.Param("id"))
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "INVALID_REQUEST", "INVALID_REQUEST", "Invalid request body", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := h.userApp.GetUserByID(ctx, id)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to get user by ID", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get user by ID"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (h *UserHandler) GetUserByUsername() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := ctx.Param("username")
		if username == "" {
			h.logger.LogErrorWithLevel("warn", "INVALID_REQUEST", "INVALID_REQUEST", "Invalid request body", nil)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Invalid request body"})
			return
		}

		user, err := h.userApp.GetUserByUsername(ctx, username)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to get user by username", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to get user by username"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	}
}

func (h *UserHandler) ListAllUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		users, err := h.userApp.ListAllUsers(ctx)
		if err != nil {
			h.logger.LogErrorWithLevel("warn", "DATABASE_ERROR", "DATABASE_ERROR", "Failed to list all users", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "Failed to list all users"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	}
}
