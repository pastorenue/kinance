package common

import (
	"github.com/google/uuid"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Email  string    `json:"email"`
	Role   string    `json:"role"`
	jwt.RegisteredClaims
}
