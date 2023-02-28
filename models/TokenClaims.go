// Package models model TokenClaims
package models

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
)

// TokenClaims is a TokenClaims
type TokenClaims struct {
	jwt.StandardClaims
	UserID uuid.UUID `json:"user_id"`
}
