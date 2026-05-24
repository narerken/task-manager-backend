package repository

import (
	"context"
	"fmt"
	"task-manager/internal/models/inputs"

	"gorm.io/gorm"
)

type UsersRolesRepositoryInterface interface {
	SetUserRole(ctx context.Context, in inputs.UsersRolesInput) error
	DeleteUserRole(ctx context.Context, id int) error
}

type UsersRolesRepository struct {
	DB *gorm.DB
}

func NewUsersRolesRepository(db *gorm.DB) *UsersRolesRepository {
	return &UsersRolesRepository{DB: db}
}

func (usersRolesRepo *UsersRolesRepository) SetUserRole(ctx context.Context, in inputs.UsersRolesInput) error {
	res := usersRolesRepo.DB.WithContext(ctx).Exec(
		"INSERT INTO users_roles (user_id, role_id) VALUES (?, ?) ON CONFLICT (user_id) DO UPDATE SET role_id = EXCLUDED.role_id",
		in.UserId,
		in.RoleId,
	)
	if res.Error != nil {
		return fmt.Errorf("failed to scan: %w", res.Error)
	}

	return nil
}

func (usersRolesRepo *UsersRolesRepository) DeleteUserRole(ctx context.Context, userId int) error {
	res := usersRolesRepo.DB.WithContext(ctx).
		Table("users_roles").
		Where("user_id = ?", userId).
		Delete(&struct{}{})
	if res.Error != nil {
		return fmt.Errorf("failed to delete: %w", res.Error)
	}

	return nil
}
