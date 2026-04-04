package handler

import (
	"net/http"
	"strconv"
	"task-manager/models"
	"task-manager/service"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	Service *service.TodoService
}

func NewTodoHandler(s *service.TodoService) *TodoHandler {
	return &TodoHandler{Service: s}
}

func (h *TodoHandler) RegisterRoutes(r *gin.Engine) {
	todos := r.Group("/todos")

	todos.POST("", h.Create)
	todos.GET("", h.GetAll)
	todos.GET("/:id", h.GetByID)
	todos.PUT("/:id", h.Update)
	todos.DELETE("/:id", h.Delete)

	todos.PATCH("/:id/complete", h.MarkCompleted)
	todos.GET("/completed", h.GetCompleted)
	todos.GET("/priority/:priority", h.GetByPriority)
}

func (h *TodoHandler) Create(c *gin.Context) {
	var todo models.Todo

	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.Service.Create(&todo)
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) GetAll(c *gin.Context) {
	todos, _ := h.Service.GetAll()
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	todo, err := h.Service.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var todo models.Todo
	c.ShouldBindJSON(&todo)
	todo.ID = uint(id)

	h.Service.Update(&todo)
	c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.Service.Delete(uint(id))

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (h *TodoHandler) MarkCompleted(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	h.Service.MarkCompleted(uint(id))

	c.JSON(http.StatusOK, gin.H{"message": "completed"})
}

func (h *TodoHandler) GetCompleted(c *gin.Context) {
	todos, _ := h.Service.GetCompleted()
	c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) GetByPriority(c *gin.Context) {
	p, _ := strconv.Atoi(c.Param("priority"))
	todos, _ := h.Service.GetByPriority(p)

	c.JSON(http.StatusOK, todos)
}
