package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/you/pawtrack/internal/models"
	"github.com/you/pawtrack/internal/service"
)

const (
	userIDKey   = "userID"
	userEmailKey = "userEmail"
	userRoleKey  = "userRole"
)

// AuthMiddleware validates JWT token and adds user info to context
func AuthMiddleware(authService service.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		// Check Bearer prefix
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Store user info in context
		c.Set(userIDKey, claims.UserID)
		c.Set(userEmailKey, claims.Email)
		c.Set(userRoleKey, claims.Role)

		c.Next()
	}
}

// RequireRole middleware checks if user has required role
func RequireRole(requiredRole models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get(userRoleKey)
		if !exists {
			c.JSON(http.StatusForbidden, gin.H{"error": "user role not found in context"})
			c.Abort()
			return
		}

		userRole, ok := role.(models.UserRole)
		if !ok || userRole != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get(userIDKey)
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, errors.New("invalid user ID type")
	}

	return id, nil
}

// GetUserRoleFromContext extracts user role from context
func GetUserRoleFromContext(c *gin.Context) (models.UserRole, error) {
	role, exists := c.Get(userRoleKey)
	if !exists {
		return "", errors.New("user role not found in context")
	}

	userRole, ok := role.(models.UserRole)
	if !ok {
		return "", errors.New("invalid user role type")
	}

	return userRole, nil
}
