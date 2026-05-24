package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
	"task-manager/internal/apperrors"
	"task-manager/internal/httpx"
	"task-manager/internal/models/inputs"
	"task-manager/internal/service"
)

type UserHandlerInterface interface {
	Register(c *gin.Context) error
	GetAllUsers(c *gin.Context) error
	GetUserById(c *gin.Context) error
	Login(c *gin.Context) error
	UpdateUserProfile(c *gin.Context) error
	UpdateUserPassword(c *gin.Context) error
	DeleteUser(c *gin.Context) error
}

type UserHandler struct {
	UserService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{UserService: userService}
}

func (userH UserHandler) Register(c *gin.Context) error {
	ctx := c.Request.Context()
	var req inputs.RegisterInput

	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid JSON")
	}

	user, err := userH.UserService.Register(ctx, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrEmailAlreadyExists) {
			return httpx.UserAlreadyExists()
		}

		log.Printf("failed to add user: %v", err)
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 201, user)
	return nil
}

func (userH *UserHandler) GetAllUsers(c *gin.Context) error {
	ctx := c.Request.Context()
	users, err := userH.UserService.GetAllUsers(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, users)
	return nil
}

func (userH *UserHandler) Login(c *gin.Context) error {
	ctx := c.Request.Context()

	var req inputs.LoginInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	tokenStr, err := userH.UserService.Login(ctx, req.Email, req.Password)
	if err != nil {
		return httpx.Unauthorized("invalid email or password")
	}

	httpx.WriteJSON(c, 200, map[string]string{"token": tokenStr})
	return nil
}

func (userH *UserHandler) GetUserById(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	user, err := userH.UserService.GetUserById(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, user)
	return nil
}

func (userH *UserHandler) UpdateUserProfile(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	var req inputs.UpdateUserProfileInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	updated, err := userH.UserService.UpdateUserProfile(ctx, id, req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, updated)
	return nil
}

func (userH *UserHandler) UpdateUserPassword(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	var req inputs.UpdatePasswordInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	err = userH.UserService.UpdateUserPassword(ctx, id, req.Password)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, "OK")
	return nil
}

func (userH *UserHandler) DeleteUser(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	err = userH.UserService.DeleteUser(ctx, id)
	if err != nil {
		return httpx.InternalError(err)
	}

	c.Status(204)
	return nil
}
