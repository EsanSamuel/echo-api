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
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByID(ctx context.Context, id int32) (*models.User, error)
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UserVerificationToken(ctx context.Context, email string, token string) error
	FindUserByToken(ctx context.Context, token string) (*models.User, error)
	VerifyUser(ctx context.Context, user *models.User) error
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

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := gorm.G[models.User](r.db).Where("email = ?", email).First(ctx)

	if err != nil {
		log.Println("User not found!", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) FindUserByID(ctx context.Context, id int32) (*models.User, error) {
	user, err := gorm.G[models.User](r.db).Where("id = ?", id).First(ctx)

	if err != nil {
		log.Println("User not found!", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context) ([]models.User, error) {
	users, err := gorm.G[models.User](r.db).Find(ctx)

	if err != nil {
		log.Println("Users not found!", err)
		return nil, err
	}

	return users, nil
}

func (r *userRepository) UserVerificationToken(ctx context.Context, email string, token string) error {
	_, err := gorm.G[models.User](r.db).Where("email = ?", email).Updates(ctx, models.User{VerificationToken: &token})

	if err != nil {
		log.Println("Error generating verification token!", err)
		return err
	}

	return nil
}

func (r *userRepository) FindUserByToken(ctx context.Context, token string) (*models.User, error) {
	user, err := gorm.G[models.User](r.db).
		Where("verification_token = ?", token).
		First(ctx)

	if err != nil {
		log.Println("User not found!", err)
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) VerifyUser(ctx context.Context, user *models.User) error {
	user.EmailVerified = true
	*user.VerificationToken = ""
	_, err := gorm.G[models.User](r.db).Where("email = ?", user.Email).Updates(ctx, *user)
	if err != nil {
		log.Println("Error updating user token", err.Error())
		return err
	}
	return err
}
