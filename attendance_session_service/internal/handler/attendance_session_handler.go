package handler

import (
	"attendance_session_service/internal/httpx"
	"attendance_session_service/internal/models"
	"attendance_session_service/internal/repository"
	"attendance_session_service/internal/service"
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AttendanceSessionHandler struct {
	service   *service.AttendanceSessionService
	validator *validator.Validate
}

func NewAttendanceSessionHandler(service *service.AttendanceSessionService) *AttendanceSessionHandler {
	return &AttendanceSessionHandler{
		service:   service,
		validator: validator.New(),
	}
}

func (h *AttendanceSessionHandler) CreateSession(c *gin.Context) error {
	var req models.CreateAttendanceSessionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}
	if errs := validateStruct(h.validator, req); errs != nil {
		return errs
	}

	currentUserID, role, err := currentUser(c)
	if err != nil {
		return err
	}

	session, err := h.service.CreateSession(c.Request.Context(), currentUserID, role, req)
	if err != nil {
		return mapAttendanceSessionError(err)
	}
	httpx.WriteJSON(c, 201, session)
	return nil
}

func (h *AttendanceSessionHandler) GetSession(c *gin.Context) error {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid session id")
	}
	currentUserID, role, err := currentUser(c)
	if err != nil {
		return err
	}

	session, err := h.service.GetSession(c.Request.Context(), currentUserID, role, sessionID)
	if err != nil {
		return mapAttendanceSessionError(err)
	}
	httpx.WriteJSON(c, 200, session)
	return nil
}

func (h *AttendanceSessionHandler) ListSessions(c *gin.Context) error {
	currentUserID, role, err := currentUser(c)
	if err != nil {
		return err
	}

	var departmentID *int
	if value := c.Query("department_id"); value != "" {
		parsed, parseErr := strconv.Atoi(value)
		if parseErr != nil {
			return httpx.BadRequest("invalid department_id")
		}
		departmentID = &parsed
	}

	var dateFrom *string
	if value := c.Query("date_from"); value != "" {
		dateFrom = &value
	}
	var dateTo *string
	if value := c.Query("date_to"); value != "" {
		dateTo = &value
	}

	sessions, err := h.service.ListSessions(c.Request.Context(), currentUserID, role, departmentID, dateFrom, dateTo)
	if err != nil {
		return mapAttendanceSessionError(err)
	}
	httpx.WriteJSON(c, 200, sessions)
	return nil
}

func (h *AttendanceSessionHandler) BulkUpsertEntries(c *gin.Context) error {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid session id")
	}

	var req models.BulkUpsertAttendanceEntriesInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}
	if len(req.Entries) == 0 {
		return httpx.BadRequest("entries are required")
	}
	if errs := validateSlice(h.validator, req.Entries); errs != nil {
		return errs
	}

	currentUserID, role, err := currentUser(c)
	if err != nil {
		return err
	}

	if err := h.service.BulkUpsertEntries(c.Request.Context(), currentUserID, role, sessionID, req); err != nil {
		return mapAttendanceSessionError(err)
	}
	c.Status(200)
	return nil
}

func (h *AttendanceSessionHandler) PublishSession(c *gin.Context) error {
	sessionID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid session id")
	}
	currentUserID, role, err := currentUser(c)
	if err != nil {
		return err
	}

	if err := h.service.PublishSession(c.Request.Context(), currentUserID, role, sessionID); err != nil {
		return mapAttendanceSessionError(err)
	}
	c.Status(200)
	return nil
}

func currentUser(c *gin.Context) (int, string, error) {
	userIDValue, ok := c.Get("user_id")
	if !ok {
		return 0, "", httpx.Unauthorized("user id is missing")
	}
	userID, ok := userIDValue.(int)
	if !ok {
		return 0, "", httpx.Unauthorized("invalid user id")
	}

	roleValue, ok := c.Get("user_role")
	if !ok {
		return 0, "", httpx.Unauthorized("user role is missing")
	}
	role, ok := roleValue.(string)
	if !ok {
		return 0, "", httpx.Unauthorized("invalid user role")
	}

	return userID, role, nil
}

func validateStruct(v *validator.Validate, payload any) error {
	if err := v.Struct(payload); err != nil {
		return buildValidationError(err)
	}
	return nil
}

func validateSlice(v *validator.Validate, items any) error {
	switch typed := items.(type) {
	case []models.AttendanceEntryUpsertInput:
		for _, item := range typed {
			if err := v.Struct(item); err != nil {
				return buildValidationError(err)
			}
		}
	default:
		if err := v.Var(items, "required"); err != nil {
			return buildValidationError(err)
		}
	}
	return nil
}

func buildValidationError(err error) error {
	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return httpx.BadRequest("validation failed")
	}
	result := make(map[string][]string)
	for _, item := range validationErrs {
		field := strings.ToLower(item.Field())
		result[field] = append(result[field], item.Tag())
	}
	return &httpx.ValidatorError{Errors: result}
}

func mapAttendanceSessionError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return httpx.NotFound("resource")
	case errors.Is(err, service.ErrUnauthorized):
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
