package service

import (
	"context"
	"log"
	"time"

	"github.com/echo/internal/models"
	"github.com/echo/internal/pkg/jwt"
	"github.com/echo/internal/pkg/password"
	"github.com/echo/internal/repository"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwt      jwt.Manager
}

func NewAuthService(userRepo repository.UserRepository, jwt jwt.Manager) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

func (s *authService) RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error) {
	emailExists, err := s.userRepo.EmailExists(ctx, req.Email)

	if err != nil {
		return nil, err
	}

	if emailExists {
		log.Println("User already exists!")
		return nil, err
	}

	HashPassword, err := password.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
	}

	createdUser, err := s.userRepo.Create(ctx, user, HashPassword)
	if err != nil {
		return nil, err
	}

	accessToken, err := s.jwt.GetAccessToken(createdUser.FirstName, createdUser.LastName, createdUser.Email, createdUser.Role, createdUser.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := s.jwt.GenerateRefreshToken()
	expiresAt := time.Now().Add(s.jwt.GetRefreshTokenExpiry())

	if err := s.userRepo.CreateRefreshToken(ctx, createdUser.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	user.PasswordHash = nil

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}
