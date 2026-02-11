package models

import "github.com/golang-jwt/jwt/v5"

type JWTClaims struct {
	Firstname string
	LastName  string
	Email     string
	Role      string
	UserID    int32
	jwt.RegisteredClaims
}
