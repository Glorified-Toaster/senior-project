package middleware

import (
	"log"
	"net/http"
	"strings"

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

func (m *AuthMiddleware) AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// get client info for logging
		clientIP := ctx.ClientIP()
		path := ctx.Request.URL.Path
		method := ctx.Request.Method

		utils.LogInfo("HTTP_SERVER",
			"Auth middleware",
			zap.String("IP address", clientIP),
			zap.String("method", method),
			zap.String("path", path))

		// getting the header from the context
		authHeader, err := ctx.Cookie("auth_token")

		if authHeader == "" {
			ctx.Redirect(http.StatusMovedPermanently, "/api/v1/login")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "authorization header required",
				"code":    "MISSING_AUTH_HEADER",
				"message": "Please include the Authorization header",
			})
			ctx.Abort()
			return
		}

		var tokenString string

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "Bearer ")
		} else if strings.HasPrefix(authHeader, "bearer ") {
			tokenString = strings.TrimPrefix(authHeader, "bearer ")
		} else {

			tokenString = strings.TrimSpace(authHeader)

			log.Printf("No Bearer prefix found, using entire header as token")

			utils.LogInfo("HTTP_SERVER",
				"Auth middleware",
				zap.String("IP address", clientIP),
				zap.String("msg", "No Bearer prefix found, using entire header as token"))
		}

		if tokenString == "" {
			log.Printf("Auth Failed - Empty token after extraction")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "invalid authorization header format",
				"code":    "INVALID_AUTH_FORMAT",
				"message": "Authorization header should be: Bearer <token>",
			})
			ctx.Abort()
			return
		}

		claims, err := m.jwt.ValidateToken(tokenString)
		if err != nil {
			utils.LogErrorWithLevel("error", "HTTP_SERVER_ERROR", "TOKEN_VALIDATION_ERROR", "Token validation failed", err)

			var errorMsg string
			var errorCode string

			if strings.Contains(err.Error(), "expired") {
				errorMsg = "token has expired"
				errorCode = "TOKEN_EXPIRED"
			} else if strings.Contains(err.Error(), "signature") {
				errorMsg = "invalid token signature"
				errorCode = "INVALID_SIGNATURE"
			} else {
				errorMsg = "invalid token"
				errorCode = "INVALID_TOKEN"
			}

			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   errorMsg,
				"code":    errorCode,
				"message": "Please login again to get a new token",
			})
			ctx.Abort()
			return
		}

		setClaimsInContext(ctx, claims)
		ctx.Next()
	}
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
