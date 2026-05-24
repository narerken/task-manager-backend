package httpx

import "github.com/gin-gonic/gin"

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}

func WriteError(c *gin.Context, status int, message string) {
	c.JSON(status, ErrorResponse{Error: message})
}
