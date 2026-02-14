package service

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/echo/internal/models"
	"github.com/echo/internal/pkg/jwt"
	"github.com/echo/internal/pkg/password"
	"github.com/echo/internal/pkg/token"
	"github.com/echo/internal/repository"
)

type AuthService interface {
	RegisterUser(ctx context.Context, req *models.RegisterRequest) (*models.AuthResponse, error)
	Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error)
	VerifyEmail(ctx context.Context, token string) error
	//FindUser(ctx context.Context, id string) (*models.User, error)
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

	verified := false
	user.EmailVerified = verified

	createdUser, err := s.userRepo.Create(ctx, user, HashPassword)
	if err != nil {
		return nil, err
	}

	verificationToken, err := token.GenerateVerificationToken()
	if err != nil {
		return nil, err
	}

	err = s.userRepo.UserVerificationToken(ctx, createdUser.Email, verificationToken)
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

func (s *authService) VerifyEmail(ctx context.Context, token string) error {
	user, err := s.userRepo.FindUserByToken(ctx, token)
	if err != nil {
		return err
	}

	err = s.userRepo.VerifyUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) Login(ctx context.Context, req *models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.New("user not found")
	}

	if user.PasswordHash == nil {
		return nil, errors.New("password not set")
	}
	err = password.Compare(*user.PasswordHash, req.Password)
	if err != nil {
		return nil, errors.New("Password does not match")
	}
	if !user.EmailVerified {
		return nil, errors.New("Email not verified!")
	}

	accessToken, err := s.jwt.GetAccessToken(user.FirstName, user.LastName, user.Email, user.Role, user.ID)
	if err != nil {
		return nil, err
	}

	refreshToken := s.jwt.GenerateRefreshToken()
	expiresAt := time.Now().Add(s.jwt.GetRefreshTokenExpiry())

	if err := s.userRepo.CreateRefreshToken(ctx, user.ID, refreshToken, expiresAt); err != nil {
		return nil, err
	}

	user.PasswordHash = nil

	//s.userRepo.UpdateUser(ctx, user)

	return &models.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil

}
