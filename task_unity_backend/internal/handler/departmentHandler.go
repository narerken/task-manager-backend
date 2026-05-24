package handler

import (
	"task-manager/internal/httpx"
	"task-manager/internal/models"
	"task-manager/internal/service"

	"github.com/gin-gonic/gin"
)

type DepartmentHandlerInterface interface {
	AddDepartment(c *gin.Context) error
	GetAllDepartments(c *gin.Context) error
}

type DepartmentHandler struct {
	DepartmentS *service.DepartmentService
}

func NewDepartmentHandler(departmentS *service.DepartmentService) *DepartmentHandler {
	return &DepartmentHandler{DepartmentS: departmentS}
}

func (departmentH *DepartmentHandler) AddDepartment(c *gin.Context) error {
	ctx := c.Request.Context()

	var department models.Department
	if err := c.ShouldBindJSON(&department); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	added, err := departmentH.DepartmentS.AddDepartment(ctx, &department)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, added)
	return nil
}

func (departmentH *DepartmentHandler) GetAllDepartments(c *gin.Context) error {
	ctx := c.Request.Context()

	departments, err := departmentH.DepartmentS.GetAllDepartments(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, departments)
	return nil
}
