package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"task-manager/internal/apperrors"
	"task-manager/internal/helpers"
	"task-manager/internal/httpx"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"
	"task-manager/internal/service"
)

type TaskHandlerInterface interface {
	AddTask(c *gin.Context) error
	GetAllTasks(c *gin.Context) error
	GetAllTasksByAssigneeId(c *gin.Context) error
	GetTaskById(c *gin.Context) error
	UpdateTask(c *gin.Context) error
	DeleteTask(c *gin.Context) error
}

type TaskHandler struct {
	TaskS *service.TaskService
}

func NewTaskHandler(taskS *service.TaskService) *TaskHandler {
	return &TaskHandler{TaskS: taskS}
}

func (taskH *TaskHandler) AddTask(c *gin.Context) error {
	ctx := c.Request.Context()

	var req inputs.AddTaskInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userIdValue, ok := c.Get("user_id")
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	userId, ok := userIdValue.(int)
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	v := helpers.NewValidator()
	errs := helpers.Validate(req, v)
	if errs != nil {
		return httpx.BadRequestValidation(errs)
	}

	task := models.Task{
		Title:       req.Title,
		Description: req.Description,
		Deadline:    req.Deadline,
		CreatorId:   userId,
		AssigneeId:  req.AssigneeId,
		Status:      req.Status,
	}

	addedTask, err := taskH.TaskS.AddTask(ctx, &task)
	if err != nil {
		if errors.Is(err, apperrors.ErrCreatorNotFound) {
			return httpx.NotFound("creator user")
		}

		if errors.Is(err, apperrors.ErrAssigneeNotFound) {
			return httpx.NotFound("assignee user")
		}

		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, addedTask)
	return nil
}

func (taskH *TaskHandler) GetAllTasks(c *gin.Context) error {
	ctx := c.Request.Context()

	tasks, err := taskH.TaskS.GetAllTasks(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, tasks)
	return nil
}

func (taskH *TaskHandler) GetAllTasksByAssigneeId(c *gin.Context) error {
	ctx := c.Request.Context()

	userIdValue, ok := c.Get("user_id")
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	userId, ok := userIdValue.(int)
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	tasks, err := taskH.TaskS.GetAllTasksByAssigneeId(ctx, userId)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, tasks)
	return nil
}

func (taskH *TaskHandler) GetTaskById(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid task ID")
	}

	task, err := taskH.TaskS.GetTaskById(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, task)
	return nil
}

func (taskH *TaskHandler) UpdateTask(c *gin.Context) error {
	ctx := c.Request.Context()

	userIdValue, ok := c.Get("user_id")
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	userId, ok := userIdValue.(int)
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid task ID")
	}

	var req inputs.UpdateTaskInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	v := helpers.NewValidator()
	errors := helpers.Validate(req, v)
	if errors != nil {
		return httpx.BadRequest("invalid status")
	}

	updated, err := taskH.TaskS.UpdateTask(ctx, userId, id, req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, updated)
	return nil
}

func (taskH *TaskHandler) DeleteTask(c *gin.Context) error {
	ctx := c.Request.Context()

	userIdValue, ok := c.Get("user_id")
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	userId, ok := userIdValue.(int)
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid task ID")
	}

	err = taskH.TaskS.DeleteTask(ctx, id, userId)
	if err != nil {
		return httpx.InternalError(err)
	}

	c.Status(204)
	return nil
}
