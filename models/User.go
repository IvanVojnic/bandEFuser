package models

import "github.com/google/uuid"

// User is a User
type User struct {
	UserID       uuid.UUID `json:"id" db:"id"`
	UserEmail    string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	RefreshToken string    `json:"refreshToken" db:"refreshToken"`
}
