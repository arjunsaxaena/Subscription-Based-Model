package models

import (
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type Plan struct {
	ID           uuid.UUID        `json:"id" db:"id"`
	Name         string           `json:"name" db:"name"`
	Price        float64          `json:"price" db:"price"`
	Features     pq.StringArray   `json:"features,omitempty" db:"features"`
	DurationDays int              `json:"duration_days" db:"duration_days"`
	Meta         *json.RawMessage `json:"meta" db:"meta"`
	BaseEntity
} 

type GetPlanFilter struct {
	ID          *uuid.UUID    `json:"id"`
	Name        *string       `json:"name"`
	Price       *float64      `json:"price"`
	Features    *pq.StringArray `json:"features"`
	DurationDays *int         `json:"duration_days"`
	IsActive    *bool         `json:"is_active"`
}

func ValidatePlan(plan *Plan) error {
	if plan.Name == "" {
		return errors.New("name is required")
	}
	if plan.Price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if plan.DurationDays <= 0 {
		return errors.New("duration_days must be greater than 0")
	}
	return nil
}