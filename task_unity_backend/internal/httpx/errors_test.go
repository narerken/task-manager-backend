package httpx

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAppErrorErrorIncludesWrappedError(t *testing.T) {
	err := &AppError{
		Message: "request failed",
		Err:     errors.New("db down"),
	}

	if got := err.Error(); got != "request failed: db down" {
		t.Fatalf("unexpected error string: %q", got)
	}
}

func TestWrapHandlerWritesAppErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler := WrapHandler(func(c *gin.Context) error {
		return NotFound("task")
	})

	handler(ctx)

	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, recorder.Code)
	}

	var body AppError
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Code != "NOT_FOUND" || body.Message != "task not found" {
		t.Fatalf("unexpected response body: %+v", body)
	}
}

func TestWrapHandlerWritesValidatorErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler := WrapHandler(func(c *gin.Context) error {
		return BadRequestValidation(map[string][]string{
			"status": {"Field 'status' is required"},
		})
	})

	handler(ctx)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, recorder.Code)
	}

	var body ValidatorError
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(body.Errors["status"]) != 1 {
		t.Fatalf("unexpected response body: %+v", body)
	}
}

func TestWrapHandlerWritesInternalErrorForUnknownError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)

	handler := WrapHandler(func(c *gin.Context) error {
		return errors.New("boom")
	})

	handler(ctx)

	if recorder.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, recorder.Code)
	}

	var body AppError
	if err := json.Unmarshal(recorder.Body.Bytes(), &body); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if body.Code != "INTERNAL_ERROR" || body.Message != "internal server error" {
		t.Fatalf("unexpected response body: %+v", body)
	}
}
