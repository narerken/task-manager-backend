package repository

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"

	"gorm.io/gorm"
)

type TaskRepositoryInterface interface {
	AddTask(ctx context.Context, task *models.Task) (*models.Task, error)
	GetTaskById(ctx context.Context, id int) (*models.Task, error)
	GetAllTasks(ctx context.Context) ([]models.Task, error)
	GetAllTasksByAssigneeId(ctx context.Context, assigneeId int) ([]models.Task, error)
	UpdateTask(ctx context.Context, id int, in inputs.UpdateTaskInput) (*models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

type TaskRepository struct {
	DB *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (taskRepo *TaskRepository) AddTask(ctx context.Context, task *models.Task) (*models.Task, error) {
	err := taskRepo.DB.WithContext(ctx).Table("tasks").Create(task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to add a task: %w", err)
	}

	return task, nil
}

func (taskRepo *TaskRepository) GetTaskById(ctx context.Context, id int) (*models.Task, error) {
	var task models.Task
	err := taskRepo.DB.WithContext(ctx).
		Table("tasks").
		Select("id", "title", "description", "deadline", "department_id", "creator_id", "assignee_id", "status", "created_at", "updated_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&task).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get task by id: %w", err)
	}

	return &task, nil
}

func (taskRepo *TaskRepository) GetAllTasks(ctx context.Context) ([]models.Task, error) {
	var tasks []models.Task
	err := taskRepo.DB.WithContext(ctx).
		Table("tasks").
		Select("id", "title", "description", "deadline", "department_id", "creator_id", "assignee_id", "status", "created_at", "updated_at", "deleted_at").
		Where("deleted_at IS NULL").
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to select all tasks: %w", err)
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) GetAllTasksByAssigneeId(ctx context.Context, assigneeId int) ([]models.Task, error) {
	var tasks []models.Task
	err := taskRepo.DB.WithContext(ctx).
		Table("tasks").
		Select("id", "title", "description", "deadline", "department_id", "creator_id", "assignee_id", "status", "created_at", "updated_at", "deleted_at").
		Where("deleted_at IS NULL AND assignee_id = ?", assigneeId).
		Find(&tasks).Error
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	return tasks, nil
}

func (taskRepo *TaskRepository) UpdateTask(ctx context.Context, id int, in inputs.UpdateTaskInput) (*models.Task, error) {
	updates := map[string]any{
		"updated_at": gorm.Expr("NOW()"),
	}
	if in.Title != nil {
		updates["title"] = *in.Title
	}
	if in.Description != nil {
		updates["description"] = *in.Description
	}
	if in.Deadline != nil {
		updates["deadline"] = *in.Deadline
	}
	if in.AssigneeId != nil {
		updates["assignee_id"] = *in.AssigneeId
	}
	if in.Status != nil {
		updates["status"] = *in.Status
	}

	res := taskRepo.DB.WithContext(ctx).
		Table("tasks").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumns(updates)
	if res.Error != nil {
		return nil, fmt.Errorf("failed to scan: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, fmt.Errorf("failed to scan: %w", gorm.ErrRecordNotFound)
	}

	var task models.Task
	err := taskRepo.DB.WithContext(ctx).
		Table("tasks").
		Select("id", "title", "description", "deadline", "department_id", "creator_id", "assignee_id", "status", "created_at", "updated_at", "deleted_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&task).Error
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &task, nil
}

func (taskRepo *TaskRepository) DeleteTask(ctx context.Context, id int) error {
	res := taskRepo.DB.WithContext(ctx).
		Table("tasks").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("deleted_at", gorm.Expr("now()"))
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}
	if res.Error != nil {
		return fmt.Errorf("failed to delete a task: %w", res.Error)
	}

	return nil
}
