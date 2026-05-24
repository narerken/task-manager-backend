package httpx

import (
	"errors"
	"fmt"
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
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}

	return e.Message
}

func NotFound(resource string) *AppError {
	return &AppError{
		Code:    "NOT_FOUND",
		Message: fmt.Sprintf("%s not found", resource),
		Status:  http.StatusNotFound,
	}
}

func BadRequest(message string) *AppError {
	return &AppError{
		Code:    "BAD_REQUEST",
		Message: message,
		Status:  http.StatusBadRequest,
	}
}

func ValidationError(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Status:  http.StatusUnprocessableEntity,
		Err:     err,
	}
}

func InternalError(err error) *AppError {
	return &AppError{
		Code:    "INTERNAL_ERROR",
		Message: "internal server error",
		Status:  http.StatusInternalServerError,
		Err:     err,
	}
}

func Unauthorized(message string) *AppError {
	return &AppError{
		Code:    "UNAUTHORIZED",
		Message: message,
		Status:  http.StatusUnauthorized,
	}
}

func UserAlreadyExists() *AppError {
	return &AppError{
		Code:    "USER_ALREADY_EXISTS",
		Message: "user already exists",
		Status:  http.StatusConflict, // 409
		Err:     nil,
	}
}

type ValidatorError struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ValidatorError) Error() string {
	return "validation failed"
}

func BadRequestValidation(errors map[string][]string) error {
	return &ValidatorError{Errors: errors}
}

type AppHandler func(c *gin.Context) error

func WrapHandler(h AppHandler) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)
		if err == nil {
			return
		}

		var appErr *AppError
		if errors.As(err, &appErr) {
			fmt.Printf("Handled error: %+v\n", appErr.Err)
			WriteJSON(c, appErr.Status, appErr)
			return
		}

		var validatorErr *ValidatorError
		if errors.As(err, &validatorErr) {
			WriteJSON(c, http.StatusBadRequest, validatorErr)
			return
		}

		internalErr := InternalError(err)
		WriteJSON(c, internalErr.Status, internalErr)
	}
}
