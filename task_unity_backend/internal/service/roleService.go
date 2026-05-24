package service

import (
	"context"
	"fmt"
	"task-manager/internal/models"
	"task-manager/internal/repository"
)

type RoleServiceInterface interface {
	AddRole(ctx context.Context, role *models.Role) (*models.Role, error)
	GetAllRoles(ctx context.Context) ([]models.Role, error)
	GetRoleById(ctx context.Context, id int) (*models.Role, error)
	UpdateRole(ctx context.Context, id int, newRole *models.Role) error
	DeleteRole(ctx context.Context, id int) error
}

type RoleService struct {
	RoleRepo *repository.RoleRepository
}

func (roleS *RoleService) AddRole(ctx context.Context, role *models.Role) (*models.Role, error) {
	checkRole, err := roleS.RoleRepo.RoleExists(ctx, role.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to check a role: %w", err)
	}

	if checkRole {
		return nil, fmt.Errorf("role is already exists")
	}

	err = roleS.RoleRepo.AddRole(ctx, role)
	if err != nil {
		return nil, fmt.Errorf("failed to add a role: %w", err)
	}

	return role, nil
}

func (roleS *RoleService) GetAllRoles(ctx context.Context) ([]models.Role, error) {
	roles, err := roleS.RoleRepo.GetAllRoles(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all roles: %w", err)
	}

	return roles, nil
}

func (roleS *RoleService) GetRoleById(ctx context.Context, id int) (*models.Role, error) {
	role, err := roleS.RoleRepo.GetRoleById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get a role: %w", err)
	}

	return role, nil
}

func (roleS *RoleService) UpdateRole(ctx context.Context, id int, newRole *models.Role) error {
	checkRole, err := roleS.RoleRepo.RoleExists(ctx, newRole.Name)
	if err != nil {
		return fmt.Errorf("failed to check a role: %w", err)
	}

	if checkRole {
		return fmt.Errorf("role is already exists")
	}

	err = roleS.RoleRepo.UpdateRole(ctx, id, newRole)
	if err != nil {
		return fmt.Errorf("failed to update a role: %w", err)
	}

	return nil
}

func (roleS *RoleService) DeleteRole(ctx context.Context, id int) error {
	err := roleS.RoleRepo.DeleteRole(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete a role: %w", err)
	}

	return nil
}
