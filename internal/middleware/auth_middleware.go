package middleware

import (
	"net/http"
	"recruitment-api/internal/repository"
	"recruitment-api/internal/service"
	"strings"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func AuthMiddleware(jwtService service.JWTService, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			return
		}

		tokenStr := parts[1]
		claims, err := jwtService.ValidateToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		// Verify user exists, is active, and not deleted
		user, err := userRepo.FindByID(claims.UserID)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user not found or inactive"})
			return
		}

		if !user.IsActive || user.IsDeleted {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "user account is inactive or deleted"})
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}
