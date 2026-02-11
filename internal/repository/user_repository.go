package repository

import (
	"context"
	"log"
	"time"

	"github.com/echo/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User, password string) (*models.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	CreateRefreshToken(ctx context.Context, user_id int32, refresh_token string, expires_at time.Time) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(dbPool *gorm.DB) UserRepository {
	return &userRepository{
		db: dbPool,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User, password string) (*models.User, error) {
	user.PasswordHash = &password
	err := gorm.G[models.User](r.db).Create(ctx, user)
	if err != nil {
		log.Println("Error inserting user", err.Error())
		return nil, err
	}
	return user, nil
}

func (r *userRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	count, err := gorm.G[models.User](r.db).
		Where("email = ?", email).
		Count(ctx, "id")
	if err != nil {
		log.Println("Error counting user email:", err.Error())
		return false, err
	}

	if count > 0 {
		// Email already exists
		return true, nil
	}

	// Email does not exist
	return false, nil
}

func (r *userRepository) CreateRefreshToken(ctx context.Context, user_id int32, refresh_token string, expires_at time.Time) error {
	_, err := gorm.G[models.User](r.db).Where("id = ?", user_id).Updates(ctx, models.User{RefreshToken: &refresh_token, ExpiresAt: expires_at})
	if err != nil {
		log.Println("Error updating user refresh token", err.Error())
		return err
	}
	return err
}
