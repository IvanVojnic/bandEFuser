package models

import "github.com/google/uuid"

// User is a User
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Password     string    `json:"password" db:"password"`
	Name         string    `json:"name" db:"name"'`
	RefreshToken string    `json:"refreshToken" db:"refreshToken"`
}
