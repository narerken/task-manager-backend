package middleware

import (
	"task-manager/internal/httpx"
	"task-manager/internal/models"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(next httpx.AppHandler, roles ...models.Role) httpx.AppHandler {
	return func(c *gin.Context) error {
		claimsValue, exists := c.Get("claims")
		if !exists {
			return httpx.Unauthorized("claims missing")
		}

		claims, ok := claimsValue.(map[string]interface{})
		if !ok {
			return httpx.Unauthorized("claims missing")
		}

		roleValue, ok := claims["role"].(string)
		if !ok {
			return httpx.Unauthorized("role missing")
		}

		allowed := false
		for _, role := range roles {
			if roleValue == role.Name {
				allowed = true
				break
			}
		}

		if !allowed {
			return httpx.Unauthorized("insufficient permissions")
		}

		return next(c)
	}
}
