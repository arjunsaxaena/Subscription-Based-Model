package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/database"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/models"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
)

type SubscriptionRepository struct {
	db *sqlx.DB
}

func NewSubscriptionRepository() *SubscriptionRepository {
	return &SubscriptionRepository{
		db: database.DB,
	}
}

func (r *SubscriptionRepository) CreateSubscription(subscription *models.Subscription) error {
	subscription.ID = uuid.New()
	subscription.CreatedAt = time.Now()
	subscription.UpdatedAt = time.Now()
	subscription.IsActive = true

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("subscriptions")
	ib.Cols("user_id", "plan_id", "status", "start_date", "end_date", "meta")
	ib.Values(
		subscription.UserID,
		subscription.PlanID,
		subscription.Status,
		subscription.StartDate,
		subscription.EndDate,
		subscription.Meta,
	)

	query, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)
	query += " RETURNING id, created_at, updated_at"

	err := r.db.QueryRowx(query, args...).Scan(&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *SubscriptionRepository) GetSubscription(filter models.GetSubscriptionFilter) ([]models.Subscription, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*")
	sb.From("subscriptions")

	if filter.ID != nil {
		sb.Where(sb.Equal("id", *filter.ID))
	}
	if filter.UserID != nil {
		sb.Where(sb.Equal("user_id", *filter.UserID))
	}
	if filter.PlanID != nil {
		sb.Where(sb.Equal("plan_id", *filter.PlanID))
	}
	if filter.Status != nil {
		sb.Where(sb.Equal("status", *filter.Status))
	}
	if filter.IsActive != nil {
		sb.Where(sb.Equal("is_active", *filter.IsActive))
	} else {
		sb.Where(sb.Equal("is_active", true))
	}

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var subscriptions []models.Subscription
	err := r.db.Select(&subscriptions, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions: %w", err)
	}
	return subscriptions, nil
}

func (r *SubscriptionRepository) UpdateSubscription(id uuid.UUID, updates map[string]interface{}) (*models.Subscription, error) {
	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("subscriptions")

	assignments := []string{}
	for key, value := range updates {
		assignments = append(assignments, ub.Assign(key, value))
	}
	assignments = append(assignments, ub.Assign("updated_at", sqlbuilder.Raw("CURRENT_TIMESTAMP")))

	ub.Set(assignments...)
	ub.Where(ub.Equal("id", id))
	
	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)
	query += " RETURNING *"

	var subscription models.Subscription
	err := r.db.Get(&subscription, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("subscription not found")
		}
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	return &subscription, nil
}

func (r *SubscriptionRepository) DeleteSubscription(id uuid.UUID) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("subscriptions")
	ub.Set(
		ub.Assign("is_active", false),
		ub.Assign("updated_at", sqlbuilder.Raw("CURRENT_TIMESTAMP")),
	)
	ub.Where(ub.Equal("id", id))

	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("subscription not found")
	}

	return nil
}
