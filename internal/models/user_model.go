package models

import "time"

type User struct {
	ID                int32     `json:"id"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	Email             string    `json:"email"`
	Role              string    `json:"role"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	PasswordHash      *string   `json:"-"`
	PhoneNumber       *string   `json:"phone_number,omitempty"`
	RefreshToken      *string   `json:"refresh_token"`
	ExpiresAt         time.Time `json:"expires_at"`
	VerificationToken *string   `json:"verification_token"`
	EmailVerified     bool      `json:"email_verified"`
}

type RegisterRequest struct {
	FirstName   string  `json:"first_name" validate:"required,min=2,max=100"`
	LastName    string  `json:"last_name" validate:"required,min=2,max=100"`
	Email       string  `json:"email" validate:"required,email"`
	PhoneNumber *string `json:"phone_number,omitempty" validate:"omitempty,min=10,max=20"`
	Password    string  `json:"password" validate:"required,min=8,max=100"`
	Role        string  `json:"role" validate:"required"`
}

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}
