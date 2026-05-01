package service

import (
	"auth-service/config"
	"auth-service/models"
	"auth-service/repo"
	"auth-service/utils"
	"fmt"
)

type AuthService struct {
	Repo *repo.UserRepository
	Cfg  *config.Config
}

func NewAuthService(r *repo.UserRepository, cfg *config.Config) *AuthService {
	return &AuthService{
		Repo: r,
		Cfg:  cfg,
	}
}

func (s *AuthService) Register(user *models.User) (*models.User, error) {
	existingUser, err := s.Repo.GetByEmail(user.Email)

	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("email already exists")
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user.Password = hashedPassword

	err = s.Repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *AuthService) Login(email, password string) (string, error) {
	user, err := s.Repo.GetByEmail(email)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	err = utils.CheckPassword(password, user.Password)
	if err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}

	tokenStr, err := utils.GenerateToken(user, s.Cfg.JWTSecret)
	if err != nil {
		return "", err
	}

	return tokenStr, nil
}
