package service

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"
	"task-manager/internal/repository"
)

type CommentServiceInterface interface {
	AddComment(ctx context.Context, in *inputs.AddCommentInput) (*models.Comment, error)
	GetAllComments(ctx context.Context) ([]models.Comment, error)
	GetCommentById(ctx context.Context, id int) (*models.Comment, error)
	UpdateComment(ctx context.Context, id int, in inputs.UpdateCommentInput) error
	DeleteComment(ctx context.Context, userId, id int) error
}

type CommentService struct {
	CommentRepo *repository.CommentRepository
	TaskS       *TaskService
}

func NewCommentService(commentRepo *repository.CommentRepository, taskS *TaskService) *CommentService {
	return &CommentService{CommentRepo: commentRepo, TaskS: taskS}
}

func (commentS *CommentService) AddComment(ctx context.Context, in *inputs.AddCommentInput) (*models.Comment, error) {
	task, err := commentS.TaskS.GetTaskById(ctx, in.TaskId)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrNotFound
		}

		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	if in.TaskId != task.Id {
		return nil, fmt.Errorf("task ids must match")
	}

	comment := models.Comment{
		Comment:   in.Comment,
		TaskId:    in.TaskId,
		CreatorId: in.CreatorId,
	}

	added, err := commentS.CommentRepo.AddComment(ctx, &comment)
	if err != nil {
		return nil, fmt.Errorf("failed to add a comment: %w", err)
	}

	return added, nil
}

func (commentS *CommentService) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	comments, err := commentS.CommentRepo.GetAllComments(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}

	return comments, nil
}

func (commentS *CommentService) GetCommentById(ctx context.Context, id int) (*models.Comment, error) {
	comment, err := commentS.CommentRepo.GetCommentById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get a comment: %w", err)
	}

	return comment, err
}

func (commentS *CommentService) UpdateComment(ctx context.Context, id int, in inputs.UpdateCommentInput) error {
	comment, err := commentS.CommentRepo.GetCommentById(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}

		return fmt.Errorf("failed to get comment by id: %w", err)
	}

	if comment.CreatorId != in.UserId {
		return fmt.Errorf("user is not the owner of the comment")
	}

	err = commentS.CommentRepo.UpdateComment(ctx, id, in.Comment)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}

		return fmt.Errorf("failed to update a comment: %w", err)
	}

	return nil
}

func (commentS *CommentService) DeleteComment(ctx context.Context, userId, id int) error {
	comment, err := commentS.CommentRepo.GetCommentById(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}

		return fmt.Errorf("failed to get comment by id: %w", err)
	}

	if comment.CreatorId != userId {
		return fmt.Errorf("user is not the owner of the comment")
	}

	err = commentS.CommentRepo.DeleteComment(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete a comment: %w", err)
	}

	return nil
}
