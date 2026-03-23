package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (auth *AuthMiddleware) AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// getting the header from the context
		authHeader, err := ctx.Cookie("auth_token")
		if err != nil {
			log.Printf("cannot get the auth header from the context: %v", err)
		}

		if authHeader == "" {
			handleUnauthorized(ctx)
			return
		}

		claims, err := auth.jwt.ValidateToken(authHeader)
		if err != nil {
			auth.logger.LogErrorWithLevel("error", "HTTP_SERVER_ERROR", "TOKEN_VALIDATION_ERROR", "Token validation failed", err)
			// Clear the invalid cookie
			ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
			handleUnauthorized(ctx)
			return
		}

		if claims != nil && claims.ExpiresAt < time.Now().Unix() {
			auth.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "TOKEN_EXPIRED", "Token is expired (manual check)", nil)
			ctx.SetCookie("auth_token", "", -1, "/", "", false, true)
			handleUnauthorized(ctx)
			return
		}

		setClaimsInContext(ctx, claims)
		ctx.Next()
	}
}

func handleUnauthorized(ctx *gin.Context) {
	// For HTMX requests, we use HX-Redirect to force a full page redirect
	if ctx.GetHeader("HX-Request") == "true" {
		ctx.Header("HX-Redirect", "/admin/login")
		ctx.AbortWithStatus(http.StatusOK) // Use 200 with HX-Redirect header
		return
	}

	// For standard requests, use a 303 See Other redirect (avoid 301 Moved Permanently)
	ctx.Redirect(http.StatusSeeOther, "/admin/login")
	ctx.Abort()
}
