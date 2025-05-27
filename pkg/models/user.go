package models

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID       `json:"id" db:"id"`
	Email        string          `json:"email" db:"email"`
	PasswordHash string          `json:"password" db:"password_hash"`
	Name         *string          `json:"name" db:"name"`
	Meta         *json.RawMessage `json:"meta" db:"meta"`
	BaseEntity
} 

type GetUserFilter struct {
	ID       *uuid.UUID `json:"id"`
	Email    *string    `json:"email"`
	Name     *string    `json:"name"`
	IsActive *bool      `json:"is_active"`
}

func ValidateUser(user *User) error {
	if user.Email == "" {
		return errors.New("email is required")
	}
	if user.PasswordHash == "" {
		return errors.New("password is required")
	}
	return nil
}

