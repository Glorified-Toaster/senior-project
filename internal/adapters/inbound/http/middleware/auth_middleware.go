package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (auth *AuthMiddleware) AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth.collectFootPrintsAndLog(ctx)

		// getting the header from the context
		authHeader, err := ctx.Cookie("auth_token")
		if err != nil {
			log.Printf("cannot get the auth header from the context")
		}

		if authHeader == "" {
			ctx.Redirect(http.StatusMovedPermanently, "/login")
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "authorization header required",
				"code":    "MISSING_AUTH_HEADER",
				"message": "Please include the Authorization header",
			})
			ctx.Abort()
			return
		}

		claims, err := auth.jwt.ValidateToken(authHeader)
		if err != nil {
			auth.logger.LogErrorWithLevel("error", "HTTP_SERVER_ERROR", "TOKEN_VALIDATION_ERROR", "Token validation failed", err)
			ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
			ctx.Redirect(http.StatusSeeOther, "/login")
			ctx.Abort()
			return
		}

		setClaimsInContext(ctx, claims)
		ctx.Next()
	}
}
