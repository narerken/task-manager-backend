package handler

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"task-manager/internal/helpers"
	"task-manager/internal/httpx"
	"task-manager/internal/models/inputs"
	"task-manager/internal/service"
)

type AttendanceHandlerInterface interface {
	AddAttendance(c *gin.Context) error
	GetAllAttendances(c *gin.Context) error
	UpdateAttendance(c *gin.Context) error
	DeleteAttendance(c *gin.Context) error
}

type AttendanceHandler struct {
	AttendanceS *service.AttendanceService
}

func NewAttendanceHandler(attendanceS *service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{AttendanceS: attendanceS}
}

func (attendanceH *AttendanceHandler) AddAttendance(c *gin.Context) error {
	ctx := c.Request.Context()

	var req inputs.AddAttendanceInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userIDValue, ok := c.Get("user_id")
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	userId, ok := userIDValue.(int)
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	req.Creator = userId
	added, err := attendanceH.AttendanceS.AddAttendance(ctx, &req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 201, added)
	return nil
}

func (attendanceH *AttendanceHandler) GetAllAttendances(c *gin.Context) error {
	ctx := c.Request.Context()

	attendances, err := attendanceH.AttendanceS.GetAllAttendances(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, attendances)
	return nil
}

func (attendanceH *AttendanceHandler) UpdateAttendance(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid attendance id")
	}

	var req inputs.UpdateAttendanceInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userIDValue, ok := c.Get("user_id")
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	userId, ok := userIDValue.(int)
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	req.MarkedBy = &userId

	v := helpers.NewValidator()
	errs := helpers.Validate(req, v)
	if errs != nil {
		return httpx.BadRequestValidation(errs)
	}

	err = attendanceH.AttendanceS.UpdateAttendance(ctx, id, &req)
	if err != nil {
		return httpx.InternalError(err)
	}

	c.Status(200)
	return nil
}

func (attendanceH *AttendanceHandler) DeleteAttendance(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid id")
	}

	err = attendanceH.AttendanceS.DeleteAttendance(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	c.Status(204)
	return nil
}
