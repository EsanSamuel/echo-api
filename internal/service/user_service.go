package service

import (
	"context"

	"github.com/echo/internal/models"
	"github.com/echo/internal/repository"
)

type UserService interface {
	FindUser(ctx context.Context, id int32) (*models.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) FindUser(ctx context.Context, id int32) (*models.User, error) {
	user, err := s.userRepo.FindUserByID(ctx, id)

	if err != nil {
		return nil, err
	}

	return user, nil
}
