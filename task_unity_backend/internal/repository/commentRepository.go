package repository

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"

	"gorm.io/gorm"
)

type CommentRepositoryInterface interface {
	AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error)
	GetCommentById(ctx context.Context, id int) (*models.Comment, error)
	GetAllComments(ctx context.Context) ([]models.Comment, error)
	UpdateComment(ctx context.Context, id int, newComment string) error
	DeleteComment(ctx context.Context, id int) error
}

type CommentRepository struct {
	DB *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{DB: db}
}

func (commentRepo *CommentRepository) AddComment(ctx context.Context, comment *models.Comment) (*models.Comment, error) {
	err := commentRepo.DB.WithContext(ctx).Table("tasks_comments").Create(comment).Error
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return comment, nil
}

func (commentRepo *CommentRepository) GetCommentById(ctx context.Context, id int) (*models.Comment, error) {
	var comment models.Comment
	err := commentRepo.DB.WithContext(ctx).
		Table("tasks_comments").
		Select("id", "comment", "task_id", "creator_id", "created_at", "updated_at", "deleted_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&comment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &comment, nil
}

func (commentRepo *CommentRepository) GetAllComments(ctx context.Context) ([]models.Comment, error) {
	var comments []models.Comment
	err := commentRepo.DB.WithContext(ctx).
		Table("tasks_comments").
		Select("id", "comment", "task_id", "creator_id", "created_at", "updated_at").
		Where("deleted_at IS NULL").
		Find(&comments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all comments: %w", err)
	}

	return comments, nil
}

func (commentRepo *CommentRepository) UpdateComment(ctx context.Context, id int, comment string) error {
	res := commentRepo.DB.WithContext(ctx).
		Table("tasks_comments").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumns(map[string]any{
			"comment":    comment,
			"updated_at": gorm.Expr("now()"),
		})
	if res.Error != nil {
		return fmt.Errorf("failed to scan: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (commentRepo *CommentRepository) DeleteComment(ctx context.Context, id int) error {
	res := commentRepo.DB.WithContext(ctx).
		Table("tasks_comments").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("deleted_at", gorm.Expr("now()"))
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	if res.Error != nil {
		return fmt.Errorf("failed to delete a comment: %w", res.Error)
	}

	return nil
}
