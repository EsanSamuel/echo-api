package service

import (
	"context"

	"github.com/echo/internal/models"
	"github.com/echo/internal/pkg/jwt"
	"github.com/echo/internal/repository"
)

type UserService interface {
	FindUser(ctx context.Context, id string) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
	jwt      jwt.Manager
}

func NewUserService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}
