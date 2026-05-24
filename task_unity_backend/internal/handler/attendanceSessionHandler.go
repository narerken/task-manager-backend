package handler

import (
	"errors"
	"strings"
	"task-manager/internal/apperrors"
	"task-manager/internal/client"
	"task-manager/internal/httpx"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
)

type AttendanceSessionHandler struct {
	client *client.AttendanceSessionClient
}

func NewAttendanceSessionHandler(attendanceSessionClient *client.AttendanceSessionClient) *AttendanceSessionHandler {
	return &AttendanceSessionHandler{client: attendanceSessionClient}
}

func (h *AttendanceSessionHandler) CreateSession(c *gin.Context) error {
	return h.proxyWithBody(c, func(userID int, role string, body []byte) (int, []byte, string, error) {
		resp, err := h.client.CreateSession(c.Request.Context(), userID, role, body)
		return proxyResponse(resp, err)
	})
}

func (h *AttendanceSessionHandler) GetSession(c *gin.Context) error {
	return h.proxyWithoutBody(c, func(userID int, role string) (int, []byte, string, error) {
		resp, err := h.client.GetSession(c.Request.Context(), userID, role, c.Param("id"))
		return proxyResponse(resp, err)
	})
}

func (h *AttendanceSessionHandler) ListSessions(c *gin.Context) error {
	return h.proxyWithoutBody(c, func(userID int, role string) (int, []byte, string, error) {
		resp, err := h.client.ListSessions(c.Request.Context(), userID, role, c.Request.URL.Query())
		return proxyResponse(resp, err)
	})
}

func (h *AttendanceSessionHandler) BulkUpsertEntries(c *gin.Context) error {
	return h.proxyWithBody(c, func(userID int, role string, body []byte) (int, []byte, string, error) {
		resp, err := h.client.BulkUpsertEntries(c.Request.Context(), userID, role, c.Param("id"), body)
		return proxyResponse(resp, err)
	})
}

func (h *AttendanceSessionHandler) PublishSession(c *gin.Context) error {
	return h.proxyWithoutBody(c, func(userID int, role string) (int, []byte, string, error) {
		resp, err := h.client.PublishSession(c.Request.Context(), userID, role, c.Param("id"))
		return proxyResponse(resp, err)
	})
}

func (h *AttendanceSessionHandler) proxyWithBody(c *gin.Context, requestFn func(userID int, role string, body []byte) (int, []byte, string, error)) error {
	userID, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	body, err := c.GetRawData()
	if err != nil {
		return httpx.BadRequest("failed to read request body")
	}

	status, responseBody, contentType, err := requestFn(userID, role, body)
	if err != nil {
		return httpx.InternalError(err)
	}

	writeProxyResponse(c, status, responseBody, contentType)
	return nil
}

func (h *AttendanceSessionHandler) proxyWithoutBody(c *gin.Context, requestFn func(userID int, role string) (int, []byte, string, error)) error {
	userID, role, err := getCurrentUserFromContext(c)
	if err != nil {
		return err
	}

	status, responseBody, contentType, err := requestFn(userID, role)
	if err != nil {
		return httpx.InternalError(err)
	}

	writeProxyResponse(c, status, responseBody, contentType)
	return nil
}

func writeProxyResponse(c *gin.Context, status int, body []byte, contentType string) {
	if contentType == "" {
		contentType = "application/json"
	}
	if len(body) == 0 {
		c.Status(status)
		return
	}
	c.Data(status, contentType, body)
}

func proxyResponse(resp *resty.Response, err error) (int, []byte, string, error) {
	if err != nil {
		return 0, nil, "", client.ProxyError(err)
	}

	contentType := resp.Header().Get("Content-Type")

	return resp.StatusCode(), resp.Body(), contentType, nil
}

func getCurrentUserFromContext(c *gin.Context) (int, string, error) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		return 0, "", httpx.BadRequest("invalid user id")
	}

	userID, ok := userIDValue.(int)
	if !ok {
		return 0, "", httpx.BadRequest("invalid user id")
	}

	claimsValue, ok := c.Get("claims")
	if !ok {
		return 0, "", httpx.Unauthorized("claims missing")
	}

	claims, ok := claimsValue.(map[string]interface{})
	if !ok {
		return 0, "", httpx.Unauthorized("claims missing")
	}

	role, ok := claims["role"].(string)
	if !ok {
		return 0, "", httpx.Unauthorized("role missing")
	}

	return userID, role, nil
}

func mapAttendanceSessionError(err error) error {
	switch {
	case errors.Is(err, apperrors.ErrNotFound):
		return httpx.NotFound("resource")
	case errors.Is(err, apperrors.ErrUnauthorized):
		return httpx.Unauthorized("insufficient permissions")
	default:
		if err != nil && (strings.Contains(err.Error(), "invalid date") || strings.Contains(err.Error(), "entries are required")) {
			return httpx.BadRequest(err.Error())
		}
		if err != nil && (strings.Contains(err.Error(), "does not belong") || strings.Contains(err.Error(), "already published") || strings.Contains(err.Error(), "department is not set")) {
			return httpx.BadRequest(err.Error())
		}
		return httpx.InternalError(err)
	}
}
