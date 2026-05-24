package handler

import (
	"strconv"
	"task-manager/internal/httpx"
	"task-manager/internal/repository"

	"github.com/gin-gonic/gin"
)

type InternalDataHandler struct {
	userRepo       *repository.UserRepository
	departmentRepo *repository.DepartmentRepository
}

func NewInternalDataHandler(userRepo *repository.UserRepository, departmentRepo *repository.DepartmentRepository) *InternalDataHandler {
	return &InternalDataHandler{
		userRepo:       userRepo,
		departmentRepo: departmentRepo,
	}
}

func (h *InternalDataHandler) GetUserById(c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid user id")
	}

	user, err := h.userRepo.GetUserById(c.Request.Context(), id)
	if err != nil {
		return mapAttendanceSessionError(err)
	}

	httpx.WriteJSON(c, 200, user)
	return nil
}

func (h *InternalDataHandler) GetDepartmentById(c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid department id")
	}

	department, err := h.departmentRepo.GetDepartmentById(c.Request.Context(), id)
	if err != nil {
		return mapAttendanceSessionError(err)
	}

	httpx.WriteJSON(c, 200, department)
	return nil
}

func (h *InternalDataHandler) GetUsersByDepartmentId(c *gin.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return httpx.BadRequest("invalid department id")
	}

	users, err := h.userRepo.GetUsersByDepartmentId(c.Request.Context(), id)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, users)
	return nil
}
