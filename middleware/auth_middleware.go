package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"stable/database/migrations"
	"stable/modules/users"
	"stable/packages/utils"
)

const (
	RoleUser = "USER"
	RoleTrainer   = "TRAINER"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"Missing authorization header",
				nil,
			))
			return
		}

		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == authHeader {
			// Jika header tidak diawali Bearer
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"Authorization header must use Bearer token",
				nil,
			))
			return
		}

		token, err := utils.ValidateToken(tokenString)
		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"Invalid token",
				nil,
			))
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"Invalid token claims",
				nil,
			))
			return
		}

		userIDValue, ok := claims["user_id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"User ID not found in token",
				nil,
			))
			return
		}

		userIDFloat, ok := userIDValue.(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"Invalid user ID in token",
				nil,
			))
			return
		}

		c.Set("user_id", int(userIDFloat))
		c.Next()
	}
}

func RequireAuthorization(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("user_id")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"User ID not found in context",
				nil,
			))
			return
		}

		userRepo := users.NewRepository(migrations.GetDB())
		user, err := userRepo.FindByID(userID.(int))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.BuildResponseFailed(
				"Unauthorized",
				"User not found",
				nil,
			))
			return
		}

		userRole := ""
		if user.Role != "" {
			userRole = user.Role
		}

		hasRole := false
		for _, role := range allowedRoles {
			if strings.EqualFold(userRole, role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, utils.BuildResponseFailed(
				"Forbidden",
				fmt.Sprintf("Access denied. Required role: %v, your role: %s", allowedRoles, userRole),
				nil,
			))
			return
		}

		c.Set("user_role", userRole)
		c.Next()
	}
}

// Shortcut roles
func RequireUser() gin.HandlerFunc {
	return RequireAuthorization(RoleUser)
}

func RequireUserOrTrainer() gin.HandlerFunc {
	return RequireAuthorization(RoleUser, RoleTrainer)
}