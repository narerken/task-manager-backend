package repository

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"
	"time"

	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserById(ctx context.Context, id int) (*models.User, error)
	GetUsersByDepartmentId(ctx context.Context, departmentId int) ([]models.User, error)
	GetAuthUserByEmail(ctx context.Context, email string) (*inputs.AuthUser, string, error)
	EmailExists(ctx context.Context, email *string) (bool, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	AddUser(ctx context.Context, user *models.User) error
	UpdateUserProfile(ctx context.Context, id int, in inputs.UpdateUserProfileInput) error
	UpdateUserPassword(ctx context.Context, id int, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (userRepo *UserRepository) GetUserById(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	err := userRepo.DB.WithContext(ctx).
		Table("users").
		Select("id", "full_name", "email", "department_id", "created_at", "updated_at", "deleted_at").
		Where("id = ?", id).
		Take(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (userRepo *UserRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	err := userRepo.DB.WithContext(ctx).
		Table("users").
		Select("id", "full_name", "email", "department_id", "created_at", "updated_at", "deleted_at").
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}

	return users, nil
}

func (userRepo *UserRepository) GetUsersByDepartmentId(ctx context.Context, departmentId int) ([]models.User, error) {
	var users []models.User
	err := userRepo.DB.WithContext(ctx).
		Table("users").
		Select("id", "full_name", "email", "department_id", "created_at", "updated_at", "deleted_at").
		Where("department_id = ? AND deleted_at IS NULL", departmentId).
		Order("full_name ASC").
		Find(&users).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch users by department: %w", err)
	}

	return users, nil
}

func (userRepo *UserRepository) GetAuthUserByEmail(ctx context.Context, email string) (*inputs.AuthUser, string, error) {
	var row struct {
		Id        int
		Email     string
		Password  string
		DeletedAt *time.Time
		Role      string
	}

	err := userRepo.DB.WithContext(ctx).
		Table("users").
		Select("users.id", "users.email", "users.password", "users.deleted_at", "roles.name as role").
		Joins("JOIN users_roles ON users.id = users_roles.user_id").
		Joins("JOIN roles ON users_roles.role_id = roles.id").
		Where("users.email = ? AND users.deleted_at IS NULL", email).
		Take(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", apperrors.ErrNotFound
		}
		return nil, "", fmt.Errorf("failed to get auth user by email: %w", err)
	}

	return &inputs.AuthUser{
		Id:        row.Id,
		Email:     row.Email,
		Password:  row.Password,
		DeletedAt: row.DeletedAt,
	}, row.Role, nil
}

func (userRepo *UserRepository) EmailExists(ctx context.Context, email *string) (bool, error) {
	var count int64
	err := userRepo.DB.WithContext(ctx).
		Table("users").
		Where("email = ?", email).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check if email exists: %w", err)
	}

	return count > 0, nil
}

func (userRepo *UserRepository) AddUser(ctx context.Context, user *models.User) error {
	return userRepo.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Table("users").Create(user).Error; err != nil {
			return fmt.Errorf("failed to add user: %w", err)
		}

		res := tx.Exec(
			"INSERT INTO users_roles (user_id, role_id) VALUES (?, (SELECT id FROM roles WHERE name = 'user'))",
			user.Id,
		)
		if res.Error != nil {
			return fmt.Errorf("failed to assign role to user: %w", res.Error)
		}

		return nil
	})
}

func (userRepo *UserRepository) UpdateUserProfile(ctx context.Context, id int, in inputs.UpdateUserProfileInput) error {
	updates := map[string]any{}
	if in.FullName != nil {
		updates["full_name"] = *in.FullName
	}
	if in.Email != nil {
		updates["email"] = *in.Email
	}
	if len(updates) == 0 {
		return fmt.Errorf("no fields to update")
	}

	res := userRepo.DB.WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		UpdateColumns(updates)
	if res.Error != nil {
		return fmt.Errorf("failed to update user profile: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (userRepo *UserRepository) UpdateUserPassword(ctx context.Context, id int, newPassword string) error {
	res := userRepo.DB.WithContext(ctx).
		Table("users").
		Where("id = ?", id).
		UpdateColumn("password", newPassword)
	if res.Error != nil {
		return fmt.Errorf("failed to update user password: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}

func (userRepo *UserRepository) DeleteUser(ctx context.Context, id int) error {
	res := userRepo.DB.WithContext(ctx).
		Table("users").
		Where("id = ? AND deleted_at IS NULL", id).
		UpdateColumn("deleted_at", gorm.Expr("now()"))
	if res.Error != nil {
		return fmt.Errorf("failed to delete user: %w", res.Error)
	}
	if res.RowsAffected == 0 {
		return apperrors.ErrNotFound
	}

	return nil
}
