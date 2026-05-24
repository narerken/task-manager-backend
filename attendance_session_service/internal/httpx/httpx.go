package httpx

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func BadRequest(message string) *AppError {
	return &AppError{Code: "BAD_REQUEST", Message: message, Status: http.StatusBadRequest}
}

func Unauthorized(message string) *AppError {
	return &AppError{Code: "UNAUTHORIZED", Message: message, Status: http.StatusUnauthorized}
}

func NotFound(resource string) *AppError {
	return &AppError{Code: "NOT_FOUND", Message: resource + " not found", Status: http.StatusNotFound}
}

func InternalError(err error) *AppError {
	return &AppError{Code: "INTERNAL_ERROR", Message: "internal server error", Status: http.StatusInternalServerError, Err: err}
}

type ValidatorError struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ValidatorError) Error() string {
	return "validation failed"
}

type AppHandler func(c *gin.Context) error

func WrapHandler(h AppHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := h(c); err != nil {
			var appErr *AppError
			if errors.As(err, &appErr) {
				c.JSON(appErr.Status, appErr)
				return
			}

			var validatorErr *ValidatorError
			if errors.As(err, &validatorErr) {
				c.JSON(http.StatusBadRequest, validatorErr)
				return
			}

			c.JSON(http.StatusInternalServerError, InternalError(err))
		}
	}
}

func WriteJSON(c *gin.Context, status int, data any) {
	c.JSON(status, data)
}
