package middleware

import (
	"github.com/Glorified-Toaster/senior-project/internal/helpers"
	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthMiddleware struct {
	jwt *helpers.JWTAuth
}

func NewAuthMiddleware(jwt *helpers.JWTAuth) *AuthMiddleware {
	return &AuthMiddleware{
		jwt: jwt,
	}
}

func collectFootPrintsAndLog(ctx *gin.Context) {
	// get client info for logging
	clientIP := ctx.ClientIP()
	path := ctx.Request.URL.Path
	method := ctx.Request.Method

	utils.LogInfo("HTTP_SERVER",
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
