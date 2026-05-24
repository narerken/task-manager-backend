package repository

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"

	"gorm.io/gorm"
)

type DepartmentRepositoryInterface interface {
	AddDepartment(ctx context.Context, department *models.Department) error
	GetDepartmentById(ctx context.Context, id int) (*models.Department, error)
	GetAllDepartments(ctx context.Context) ([]models.Department, error)
	UpdateDepartment(ctx context.Context, id int, newDepartment *models.Department) error
	DeleteDepartment(ctx context.Context, id int) error
}

type DepartmentRepository struct {
	DB *gorm.DB
}

func NewDepartmentRepository(db *gorm.DB) *DepartmentRepository {
	return &DepartmentRepository{DB: db}
}

func (departmentRepo *DepartmentRepository) AddDepartment(ctx context.Context, department *models.Department) error {
	err := departmentRepo.DB.WithContext(ctx).Table("departments").Create(department).Error
	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (departmentRepo *DepartmentRepository) GetDepartmentById(ctx context.Context, id int) (*models.Department, error) {
	var department models.Department
	err := departmentRepo.DB.WithContext(ctx).
		Table("departments").
		Select("id", "name", "created_at", "updated_at", "deleted_at").
		Where("id = ? AND deleted_at IS NULL", id).
		Take(&department).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &department, nil
}

func (departmentRepo *DepartmentRepository) GetAllDepartments(ctx context.Context) ([]models.Department, error) {
	var departments []models.Department
	err := departmentRepo.DB.WithContext(ctx).
		Table("departments").
		Select("id", "name", "created_at", "updated_at").
		Where("deleted_at IS NULL").
		Find(&departments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to select all departments: %w", err)
	}

	return departments, nil
}

func (departmentRepo *DepartmentRepository) UpdateDepartment(ctx context.Context, id int, newDepartment *models.Department) error {
	res := departmentRepo.DB.WithContext(ctx).
		Table("departments").
		Where("id = ?", id).
		UpdateColumns(map[string]any{
			"name":       newDepartment.Name,
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

func (departmentRepo *DepartmentRepository) DeleteDepartment(ctx context.Context, id int) error {
	res := departmentRepo.DB.WithContext(ctx).
		Table("departments").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("deleted_at", gorm.Expr("now()"))
	if res.Error != nil {
		return fmt.Errorf("failed to delete a department: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (departmentRepo *DepartmentRepository) DepartmentExists(ctx context.Context, departmentName string) (bool, error) {
	var count int64
	err := departmentRepo.DB.WithContext(ctx).
		Table("departments").
		Where("name = ?", departmentName).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to scan: %w", err)
	}

	return count > 0, nil
}
