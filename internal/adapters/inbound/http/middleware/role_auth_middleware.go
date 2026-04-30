package middleware

import (
	"net/http"

	"uot-exam/internal/adapters/inbound/http/helpers"
	"uot-exam/internal/domain"

	"github.com/gin-gonic/gin"
)

// RoleAuthMiddleware restricts access to routes based on one or more allowed roles.
// It must be used after AuthenticationMiddleware, which sets the claims in context.
func (m *AuthMiddleware) RoleAuthMiddleware(allowedRoles ...domain.UserRole) gin.HandlerFunc {
	roleSet := make(map[domain.UserRole]struct{}, len(allowedRoles))
	for _, r := range allowedRoles {
		roleSet[r] = struct{}{}
	}

	return func(c *gin.Context) {
		// Retrieve the claims set by AuthenticationMiddleware.
		raw, exists := c.Get("claims")
		if !exists {
			m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "MISSING_CLAIMS", "Claims not found in context", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		claims, ok := raw.(*helpers.Claims)
		if !ok || claims == nil {
			m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "INVALID_CLAIMS", "Claims type assertion failed", nil)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		if !claims.IsActive {
			m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "INACTIVE_USER", "Inactive user attempted access", nil)
			c.JSON(http.StatusForbidden, gin.H{"error": "Account is disabled"})
			c.Abort()
			return
		}

		userRole := domain.UserRole(claims.Role)
		if _, allowed := roleSet[userRole]; !allowed {
			m.logger.LogErrorWithLevel("warn", "AUTH_ERROR", "FORBIDDEN", "User role not permitted", nil)
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
			c.Abort()
			return
		}

		c.Next()
	}
}
