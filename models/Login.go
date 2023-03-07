// Package models model Login
package models

import "github.com/google/uuid"

// Login is a Login
type Login struct {
	ID       uuid.UUID `json:"id" db:"id"`
	Password string    `json:"password" db:"password"`
	Name     string    `json:"name" db:"name"`
}
