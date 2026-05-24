package middleware

import (
	"github.com/gin-gonic/gin"
)

func InternalServiceMiddleware(serviceToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if serviceToken == "" {
			c.AbortWithStatusJSON(500, gin.H{"code": "INTERNAL_ERROR", "message": "internal service token is not configured"})
			return
		}

		if c.GetHeader("X-Service-Token") != serviceToken {
			c.AbortWithStatusJSON(401, gin.H{"code": "UNAUTHORIZED", "message": "invalid service token"})
			return
		}

		c.Next()
	}
}
