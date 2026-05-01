package middleware

import (
	"log"
	"net/http"
	"strings"
	"todo-service/client"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authClient *client.AuthClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid auth header"})
			c.Abort()
			return
		}

		token := parts[1]

		log.Println("TOKEN:", token)

		res, err := authClient.ValidateToken(token)

		log.Println("AUTH RESPONSE:", res, err)

		if err != nil || res == nil || !res.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("user_id", res.UserID)
		c.Next()
	}
}
