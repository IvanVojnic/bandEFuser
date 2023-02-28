package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}
