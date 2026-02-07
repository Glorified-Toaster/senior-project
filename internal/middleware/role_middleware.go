package middleware

import (
	"net/http"
	"strings"

	"github.com/Glorified-Toaster/senior-project/internal/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (m *AuthMiddleware) RequireRoles(allowedRoles ...string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieve the role from context (set by AuthenticationMiddleware)
		userRole, exists := ctx.Get("role")
		if !exists {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"error":   "unauthorized",
				"code":    "MISSING_ROLE",
				"message": "User role not found in context context",
			})
			ctx.Abort()
			return
		}

		// Assert that the role is a string
		roleStr, ok := userRole.(string)
		if !ok {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error":   "internal server error",
				"code":    "ROLE_TYPE_ERROR",
				"message": "User role is not of type string",
			})
			ctx.Abort()
			return
		}

		// Check if the user's role is inside the list of allowedRoles
		isAllowed := false
		for _, role := range allowedRoles {
			// standardizing to lowercase to avoid case-sensitivity issues
			if strings.EqualFold(roleStr, role) {
				isAllowed = true
				break
			}
		}

		// If not allowed, return 403 Forbidden
		if !isAllowed {
			utils.LogInfo("HTTP_SERVER",
				"Auth middleware",
				zap.String("msg", "Access denied: Insufficient privileges"),
				zap.String("user_role", roleStr),
				zap.Strings("required_roles", allowedRoles))

			ctx.JSON(http.StatusForbidden, gin.H{
				"error":   "forbidden",
				"code":    "INSUFFICIENT_PERMISSIONS",
				"message": "You do not have permission to access this resource",
			})
			ctx.Abort()
			return
		}

		// User is authorized, proceed
		ctx.Next()
	}
}
