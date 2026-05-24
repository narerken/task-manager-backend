package service

import (
	"context"
	"errors"
	"fmt"
	"task-manager/internal/apperrors"
	"task-manager/internal/auth"
	"task-manager/internal/models"
	"task-manager/internal/models/inputs"
	"task-manager/internal/repository"
	"task-manager/internal/utils"
	"unicode/utf8"
)

type UserServiceInterface interface {
	Register(ctx context.Context, input inputs.RegisterInput) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, id int) (*models.User, error)
	UpdateUserProfile(ctx context.Context, id int, in inputs.UpdateUserProfileInput) (*models.User, error)
	UpdateUserPassword(ctx context.Context, id int, newPassword string) error
	DeleteUser(ctx context.Context, id int) error
}

type UserService struct {
	UserRepo  *repository.UserRepository
	JwtSecret *auth.JWTSecret
}

func NewUserService(userR *repository.UserRepository, jwtSecret *auth.JWTSecret) *UserService {
	return &UserService{
		UserRepo:  userR,
		JwtSecret: jwtSecret,
	}
}

func (userS *UserService) Register(ctx context.Context, input inputs.RegisterInput) (*models.User, error) {
	checkEmail, err := userS.UserRepo.EmailExists(ctx, &input.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	if checkEmail {
		return nil, apperrors.ErrEmailAlreadyExists
	}

	if utf8.RuneCountInString(input.Password) < 8 {
		return nil, apperrors.ErrWeakPassword
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := models.User{
		FullName:     input.FullName,
		Email:        input.Email,
		Password:     hashedPassword,
		DepartmentId: input.DepartmentId,
	}

	err = userS.UserRepo.AddUser(ctx, &user)
	if err != nil {
		return nil, fmt.Errorf("failed to add user: %w", err)
	}

	user.Password = ""
	return &user, nil
}

func (userS *UserService) Login(ctx context.Context, email, password string) (string, error) {
	authUser, role, err := userS.UserRepo.GetAuthUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return "", apperrors.ErrUnauthorized
		}
		return "", fmt.Errorf("failed to find user by email: %w", err)
	}

	valid, err := utils.ValidatePassword(password, authUser.Password)
	if err != nil {
		return "", fmt.Errorf("failed to validate password: %w", err)
	}

	if !valid {
		return "", apperrors.ErrUnauthorized
	}

	token, err := userS.JwtSecret.GenerateToken(authUser, role)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}

	return token, nil
}

func (userS *UserService) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := userS.UserRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	return users, nil
}

func (userS *UserService) GetUserById(ctx context.Context, id int) (*models.User, error) {
	user, err := userS.UserRepo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return user, nil
}

func (userS *UserService) UpdateUserProfile(ctx context.Context, id int, in inputs.UpdateUserProfileInput) (*models.User, error) {
	if in.Email != nil {
		exists, err := userS.UserRepo.EmailExists(ctx, in.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email: %w", err)
		}
		if exists {
			return nil, apperrors.ErrEmailAlreadyExists
		}
	}

	err := userS.UserRepo.UpdateUserProfile(ctx, id, in)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	user, err := userS.UserRepo.GetUserById(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return nil, apperrors.ErrNotFound
		}
		return nil, fmt.Errorf("failed to fetch updated user: %w", err)
	}

	return user, nil
}

func (userS *UserService) UpdateUserPassword(ctx context.Context, id int, newPassword string) error {
	if utf8.RuneCountInString(newPassword) < 8 {
		return apperrors.ErrWeakPassword
	}

	hash, err := utils.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = userS.UserRepo.UpdateUserPassword(ctx, id, hash)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (userS *UserService) DeleteUser(ctx context.Context, id int) error {
	err := userS.UserRepo.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, apperrors.ErrNotFound) {
			return apperrors.ErrNotFound
		}
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
