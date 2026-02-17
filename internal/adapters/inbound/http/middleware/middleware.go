package middleware

import (
	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/adapters/outbound/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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

func (auth *AuthMiddleware) collectFootPrintsAndLog(ctx *gin.Context) {
	// get client info for logging
	clientIP := ctx.ClientIP()
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	auth.logger.LogInfo("HTTP_SERVER",
		"Auth middleware",
		zap.String("IP address", clientIP),
		zap.String("method", method),
		zap.String("path", path))
}

func setClaimsInContext(ctx *gin.Context, claims *helpers.Claims) {
	ctx.Set("claims", claims)
	ctx.Set("userID", claims.UserID)
	ctx.Set("email", claims.Email)
	ctx.Set("role", claims.Role)
	ctx.Set("firstName", claims.FirstName)
	ctx.Set("lastName", claims.LastName)
	ctx.Set("studentID", claims.StudentID)
	ctx.Set("department", claims.Department)
	ctx.Set("isActive", claims.IsActive)
}
