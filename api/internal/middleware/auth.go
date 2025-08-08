package middleware

import (
	"api/internal/models"
	"api/internal/services"
	"api/internal/utils"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.UnauthorizedResponse(c, "Authorization header is required")
			c.Abort()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			utils.UnauthorizedResponse(c, "Invalid authorization header format")
			c.Abort()
			return
		}

		claims, err := authService.ValidateToken(bearerToken[1])
		if err != nil {
			utils.UnauthorizedResponse(c, "Invalid token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("microsoft_id", claims.MicrosoftID)
		c.Set("claims", claims)

		c.Next()
	}
}

func RequireRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ForbiddenResponse(c, "User role not found")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.ForbiddenResponse(c, "Invalid user role format")
			c.Abort()
			return
		}

		userRoleEnum := models.UserRole(roleStr)

		switch role {
		case models.RoleAdmin:
			if userRoleEnum != models.RoleAdmin {
				utils.ForbiddenResponse(c, "Admin role required")
				c.Abort()
				return
			}
		case models.RoleManager:
			if userRoleEnum != models.RoleAdmin && userRoleEnum != models.RoleManager {
				utils.ForbiddenResponse(c, "Manager role or higher required")
				c.Abort()
				return
			}
		case models.RoleEmployee:
		}

		c.Next()
	}
}

func RequireOwnerOrRole(role models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			utils.ForbiddenResponse(c, "User ID not found")
			c.Abort()
			return
		}

		currentUserID, ok := userID.(uint)
		if !ok {
			utils.ForbiddenResponse(c, "Invalid user ID format")
			c.Abort()
			return
		}

		targetUserIDStr := c.Param("id")
		if targetUserIDStr == "" {
			targetUserIDStr = c.Param("user_id")
		}

		if targetUserIDStr != "" {
			var targetUserID uint
			if _, err := utils.GetValidator().Var(&targetUserID, "required,numeric"); err == nil {
				if currentUserID == targetUserID {
					c.Next()
					return
				}
			}
		}

		userRole, exists := c.Get("user_role")
		if !exists {
			utils.ForbiddenResponse(c, "User role not found")
			c.Abort()
			return
		}

		roleStr, ok := userRole.(string)
		if !ok {
			utils.ForbiddenResponse(c, "Invalid user role format")
			c.Abort()
			return
		}

		userRoleEnum := models.UserRole(roleStr)

		switch role {
		case models.RoleAdmin:
			if userRoleEnum != models.RoleAdmin {
				utils.ForbiddenResponse(c, "Admin role required")
				c.Abort()
				return
			}
		case models.RoleManager:
			if userRoleEnum != models.RoleAdmin && userRoleEnum != models.RoleManager {
				utils.ForbiddenResponse(c, "Manager role or higher required")
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func OptionalAuth(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := authService.ValidateToken(bearerToken[1])
		if err != nil {
			c.Next()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_email", claims.Email)
		c.Set("user_role", claims.Role)
		c.Set("microsoft_id", claims.MicrosoftID)
		c.Set("claims", claims)

		c.Next()
	}
}