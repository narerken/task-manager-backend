package handler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"strconv"
	"task-manager/internal/apperrors"
	"task-manager/internal/httpx"
	"task-manager/internal/models/inputs"
	"task-manager/internal/service"
)

type CommentHandlerInterface interface {
	AddComment(c *gin.Context) error
	GetAllComments(c *gin.Context) error
	UpdateComment(c *gin.Context) error
	DeleteComment(c *gin.Context) error
}

type CommentHandler struct {
	CommentS *service.CommentService
}

func NewCommentHandler(commentS *service.CommentService) *CommentHandler {
	return &CommentHandler{CommentS: commentS}
}

func (commentH *CommentHandler) AddComment(c *gin.Context) error {
	ctx := c.Request.Context()

	var req inputs.AddCommentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	userIDValue, ok := c.Get("user_id")
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	creatorId, ok := userIDValue.(int)
	if !ok {
		return httpx.BadRequest("invalid user id")
	}

	req.CreatorId = creatorId

	added, err := commentH.CommentS.AddComment(ctx, &req)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 201, added)
	return nil
}

func (commentH *CommentHandler) GetAllComments(c *gin.Context) error {
	ctx := c.Request.Context()

	comments, err := commentH.CommentS.GetAllComments(ctx)
	if err != nil {
		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, comments)
	return nil
}

func (commentH *CommentHandler) UpdateComment(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	userIDValue, ok := c.Get("user_id")
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	userId, ok := userIDValue.(int)
	if !ok {
		return httpx.Unauthorized("user_id missing")
	}

	var req inputs.UpdateCommentInput
	if err := c.ShouldBindJSON(&req); err != nil {
		return httpx.BadRequest("invalid request body")
	}

	req.UserId = userId

	err = commentH.CommentS.UpdateComment(ctx, id, req)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return httpx.NotFound("comment")
		}

		return httpx.InternalError(err)
	}

	httpx.WriteJSON(c, 200, "ok")
	return nil
}

func (commentH *CommentHandler) DeleteComment(c *gin.Context) error {
	ctx := c.Request.Context()

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return httpx.BadRequest("invalid ID")
	}

	userIDValue, ok := c.Get("user_id")
	if !ok {
		return httpx.BadRequest("user_id missing")
	}

	userId, ok := userIDValue.(int)
	if !ok {
		return httpx.BadRequest("user_id missing")
	}

	err = commentH.CommentS.DeleteComment(ctx, userId, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return httpx.NotFound("comment")
		}

		return httpx.InternalError(err)
	}

	c.Status(204)
	return nil
}
