package repository

import (
	"context"
	"fmt"
	"task-manager/internal/models"

	"gorm.io/gorm"
)

type RoleRepositoryInterface interface {
	AddRole(ctx context.Context, role *models.Role) error
	GetRoleById(ctx context.Context, id int) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]models.Role, error)
	UpdateRole(ctx context.Context, id int, newRole *models.Role) error
	DeleteRole(ctx context.Context, id int) error
	RoleExists(ctx context.Context, roleName string) (bool, error)
}

type RoleRepository struct {
	DB *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{DB: db}
}

func (roleRepo *RoleRepository) AddRole(ctx context.Context, role *models.Role) error {
	err := roleRepo.DB.WithContext(ctx).Table("roles").Create(role).Error
	if err != nil {
		return fmt.Errorf("failed to scan: %w", err)
	}

	return nil
}

func (roleRepo *RoleRepository) GetRoleById(ctx context.Context, id int) (*models.Role, error) {
	var role models.Role
	err := roleRepo.DB.WithContext(ctx).
		Table("roles").
		Select("id", "name", "department_id").
		Where("id = ?", id).
		Take(&role).Error
	if err != nil {
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	return &role, nil
}

func (roleRepo *RoleRepository) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	var roles []models.Role
	err := roleRepo.DB.WithContext(ctx).
		Table("roles").
		Select("id", "name", "department_id").
		Find(&roles).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	return roles, nil
}

func (roleRepo *RoleRepository) UpdateRole(ctx context.Context, id int, newRole *models.Role) error {
	res := roleRepo.DB.WithContext(ctx).
		Table("roles").
		Where("id = ?", id).
		UpdateColumns(map[string]any{
			"name":          newRole.Name,
			"department_id": newRole.DepartmentId,
		})
	if res.Error != nil {
		return fmt.Errorf("failed to update a role: %w", res.Error)
	}

	return nil
}

func (roleRepo *RoleRepository) DeleteRole(ctx context.Context, id int) error {
	res := roleRepo.DB.WithContext(ctx).
		Table("roles").
		Where("id = ?", id).
		Delete(&struct{}{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete a role: %w", res.Error)
	}

	return nil
}

func (roleRepo *RoleRepository) RoleExists(ctx context.Context, roleName string) (bool, error) {
	var count int64
	err := roleRepo.DB.WithContext(ctx).
		Table("roles").
		Where("name = ?", roleName).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to scan: %w", err)
	}

	return count > 0, nil
}
