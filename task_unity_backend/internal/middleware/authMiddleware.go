package middleware

import (
	"strings"
	"task-manager/internal/auth"
	"task-manager/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret *auth.JWTSecret, next httpx.AppHandler) httpx.AppHandler {
	return func(c *gin.Context) error {
		tokenStr := c.GetHeader("Authorization")

		if tokenStr == "" {
			return httpx.Unauthorized("empty token")
		}

		const bearerPrefix = "Bearer "

		if !strings.HasPrefix(tokenStr, bearerPrefix) {
			return httpx.Unauthorized("invalid token format")
		}

		tokenStr = strings.TrimPrefix(tokenStr, bearerPrefix)

		token, err := jwtSecret.ValidateToken(tokenStr)
		if err != nil {
			return httpx.ValidationError(err)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return httpx.Unauthorized("invalid token claims")
		}

		userId := int(claims["user_id"].(float64))
		role, ok := claims["role"].(string)
		if !ok {
			return httpx.Unauthorized("role missing in token")
		}

		c.Set("user_id", userId)
		c.Set("claims", map[string]interface{}{
			"user_id": userId,
			"role":    role,
		})
		return next(c)
	}
}
