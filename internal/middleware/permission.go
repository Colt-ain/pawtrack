package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/repository"
)

var permissionRepo repository.PermissionRepository

// InitPermissionMiddleware initializes the permission repository for middleware
func InitPermissionMiddleware(repo repository.PermissionRepository) {
	permissionRepo = repo
}

// RequirePermission creates middleware that checks if user has the required permission
func RequirePermission(permissionName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has the required permission
		hasPermission, err := permissionRepo.HasPermission(userID, permissionName)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RequireAnyPermission checks if user has at least one of the given permissions
func RequireAnyPermission(permissionNames ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := GetUserIDFromContext(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has any of the required permissions
		hasPermission, err := permissionRepo.HasAnyPermission(userID, permissionNames)
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserPermissions helper to get user permissions in handlers
func GetUserPermissions(c *gin.Context) ([]string, error) {
	userID, err := GetUserIDFromContext(c)
	if err != nil {
		return nil, err
	}

	return permissionRepo.GetUserPermissions(userID)
}
