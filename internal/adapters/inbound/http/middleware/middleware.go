package middleware

import (
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/outbound/logger"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	jwt    *helpers.JWTAuth
	logger *logger.Logger
}

func NewAuthMiddleware(jwt *helpers.JWTAuth, logger *logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwt:    jwt,
		logger: logger,
	}
}

func setClaimsInContext(ctx *gin.Context, claims *helpers.Claims) {
	ctx.Set("claims", claims)
	ctx.Set("userID", claims.UserID)
	ctx.Set("role", claims.Role)
	ctx.Set("fullname", claims.FullName)
	ctx.Set("userName", claims.UserName)
	ctx.Set("username", claims.UserName)
	ctx.Set("isActive", claims.IsActive)
}
