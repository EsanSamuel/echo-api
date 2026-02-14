package password

import (
	"errors"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	minPasswordLength = 8
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		errors.New("Failed to hash password")
		return "", err
	}

	return string(hashedPassword), nil
}

func Compare(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func Validate(password string) error {
	if len(password) < minPasswordLength {
		return errors.New("password must be at least 8 characters")
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsSymbol(char) || unicode.IsPunct(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		return errors.New("Password should contain lowercase")
	}
	if !hasUpper {
		return errors.New("Password should contain uppercase")
	}
	if !hasSpecial {
		return errors.New("Password should contain special characters")
	}
	if !hasNumber {
		return errors.New("Password should contain numbers")
	}

	return nil
}
