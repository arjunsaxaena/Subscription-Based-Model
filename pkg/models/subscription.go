package models

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "ACTIVE"
	SubscriptionStatusInactive  SubscriptionStatus = "INACTIVE"
	SubscriptionStatusCancelled SubscriptionStatus = "CANCELLED"
	SubscriptionStatusExpired   SubscriptionStatus = "EXPIRED"
)

type Subscription struct {
	ID        uuid.UUID          `json:"id" db:"id"`
	UserID    uuid.UUID          `json:"user_id" db:"user_id"`
	PlanID    uuid.UUID          `json:"plan_id" db:"plan_id"`
	Status    SubscriptionStatus `json:"status" db:"status"`
	StartDate time.Time          `json:"start_date" db:"start_date"`
	EndDate   time.Time          `json:"end_date" db:"end_date"`
	Meta      *json.RawMessage    `json:"meta" db:"meta"`
	BaseEntity
} 

type GetSubscriptionFilter struct {
	ID       *uuid.UUID `json:"id"`
	UserID   *uuid.UUID `json:"user_id"`
	PlanID   *uuid.UUID `json:"plan_id"`
	Status   *SubscriptionStatus `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	IsActive  *bool      `json:"is_active"`
}

func ValidateSubscription(subscription *Subscription) error {
	if subscription.UserID == uuid.Nil {
		return errors.New("user_id is required")
	}
	if subscription.PlanID == uuid.Nil {
		return errors.New("plan_id is required")
	}
	if subscription.Status == "" {
		return errors.New("status is required")
	}
	if subscription.StartDate.IsZero() {
		return errors.New("start_date is required")
	}
	if subscription.EndDate.IsZero() {
		return errors.New("end_date is required")
	}

	return nil
}

