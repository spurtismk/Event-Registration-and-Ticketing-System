package services

import (
	"context"
	"errors"

	"event_registration/internal/models"
	"event_registration/internal/repositories"
	"event_registration/internal/utils"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password string, role models.Role) (*models.User, error)
	Login(ctx context.Context, email, password, jwtSecret string) (string, error)
}

type authService struct {
	userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Register(ctx context.Context, name, email, password string, role models.Role) (*models.User, error) {
	// Check if user exists
	existingUser, _ := s.userRepo.FindByEmail(ctx, email)
	if existingUser != nil {
		return nil, errors.New("email already in use")
	}

	hash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	if role == "" {
		role = models.RoleAudience
	}

	user := &models.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         role,
	}

	err = s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password, jwtSecret string) (string, error) {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if !user.IsActive {
		return "", errors.New("user is inactive")
	}

	if !utils.CheckPasswordHash(password, user.PasswordHash) {
		return "", errors.New("invalid credentials")
	}

	token, err := utils.GenerateToken(user.ID.String(), user.Role, jwtSecret)
	if err != nil {
		return "", err
	}

	return token, nil
}
