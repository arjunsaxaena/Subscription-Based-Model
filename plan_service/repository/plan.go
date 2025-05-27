package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/database"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/models"
	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PlanRepository struct {
	db *sqlx.DB
}

func NewPlanRepository() *PlanRepository {
	return &PlanRepository{
		db: database.DB,
	}
}

func (r *PlanRepository) CreatePlan(plan *models.Plan) error {
	if err := models.ValidatePlan(plan); err != nil {
		return err
	}

	plan.ID = uuid.New()
	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()
	plan.IsActive = true

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("plans")
	ib.Cols("name", "price", "features", "duration_days", "meta")
	ib.Values(
		plan.Name,
		plan.Price,
		pq.Array(plan.Features), // plan.Features is now []string
		plan.DurationDays,
		plan.Meta,
	)

	query, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)
	query += " RETURNING id, created_at, updated_at"

	err := r.db.QueryRowx(query, args...).Scan(&plan.ID, &plan.CreatedAt, &plan.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("plan with this name already exists")
		}
		return fmt.Errorf("failed to create plan: %w", err)
	}

	return nil
}

func (r *PlanRepository) GetPlan(filter models.GetPlanFilter) ([]models.Plan, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*")
	sb.From("plans")

	if filter.ID != nil {
		sb.Where(sb.Equal("id", *filter.ID))
	}
	if filter.Name != nil {
		sb.Where(sb.Equal("name", *filter.Name))
	}
	if filter.Price != nil {
		sb.Where(sb.Equal("price", *filter.Price))
	}
	if filter.Features != nil {
		sb.Where(sb.Equal("features", *filter.Features))
	}
	if filter.DurationDays != nil {
		sb.Where(sb.Equal("duration_days", *filter.DurationDays))
	}
	if filter.IsActive != nil {
		sb.Where(sb.Equal("is_active", *filter.IsActive))
	} else {
		sb.Where(sb.Equal("is_active", true))
	}

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var plans []models.Plan
	err := r.db.Select(&plans, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get plans: %w", err)
	}
	return plans, nil
}

func (r *PlanRepository) UpdatePlan(id uuid.UUID, updates map[string]interface{}) (*models.Plan, error) {
	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("plans")

	assignments := []string{}
	updated := false

	if name, ok := updates["name"]; ok {
		assignments = append(assignments, ub.Assign("name", name))
		updated = true
	}
	if price, ok := updates["price"]; ok {
		assignments = append(assignments, ub.Assign("price", price))
		updated = true
	}
	if features, ok := updates["features"]; ok {
		if arr, ok := features.([]string); ok {
			assignments = append(assignments, ub.Assign("features", pq.Array(arr)))
		} else {
			assignments = append(assignments, ub.Assign("features", features))
		}
		updated = true
	}
	if durationDays, ok := updates["duration_days"]; ok {
		assignments = append(assignments, ub.Assign("duration_days", durationDays))
		updated = true
	}
	if meta, ok := updates["meta"]; ok {
		assignments = append(assignments, ub.Assign("meta", meta))
		updated = true
	}

	assignments = append(assignments, ub.Assign("updated_at", sqlbuilder.Raw("CURRENT_TIMESTAMP")))

	if !updated {
		return nil, errors.New("no valid fields to update")
	}

	ub.Set(assignments...)
	ub.Where(ub.Equal("id", id))
	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)
	query += " RETURNING *"

	var plan models.Plan
	err := r.db.Get(&plan, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("plan not found")
		}
		return nil, fmt.Errorf("failed to update plan: %w", err)
	}

	return &plan, nil
}

func (r *PlanRepository) DeletePlan(id uuid.UUID) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("plans")
	ub.Set(
		ub.Assign("is_active", false),
		ub.Assign("updated_at", sqlbuilder.Raw("CURRENT_TIMESTAMP")),
	)
	ub.Where(ub.Equal("id", id))

	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete plan: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("plan not found")
	}

	return nil
}
