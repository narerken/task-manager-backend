package middleware

import (
	"attendance_session_service/internal/httpx"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ServiceAuthMiddleware(serviceToken string, next httpx.AppHandler) httpx.AppHandler {
	return func(c *gin.Context) error {
		if c.GetHeader("X-Service-Token") != serviceToken {
			return httpx.Unauthorized("invalid service token")
		}

		userID, err := strconv.Atoi(c.GetHeader("X-User-Id"))
		if err != nil {
			return httpx.Unauthorized("invalid user id")
		}
		role := c.GetHeader("X-User-Role")
		if role == "" {
			return httpx.Unauthorized("user role is missing")
		}

		c.Set("user_id", userID)
		c.Set("user_role", role)
		return next(c)
	}
}
