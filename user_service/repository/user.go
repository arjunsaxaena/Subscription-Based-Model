package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/huandu/go-sqlbuilder"
	"github.com/jmoiron/sqlx"

	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/database"
	"github.com/arjunsaxaena/Subscription-Based-Model.git/pkg/models"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		db: database.DB,
	}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	if err := models.ValidateUser(user); err != nil {
		return err
	}

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.IsActive = true

	ib := sqlbuilder.PostgreSQL.NewInsertBuilder()
	ib.InsertInto("users")
	ib.Cols("email", "password_hash", "name", "meta")
	ib.Values(user.Email, user.PasswordHash, user.Name, user.Meta)
	
	query, args := ib.BuildWithFlavor(sqlbuilder.PostgreSQL)
	query += " RETURNING id, created_at, updated_at"

	err := r.db.QueryRowx(query, args...).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return errors.New("user with this email already exists")
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *UserRepository) GetUser(filter models.GetUserFilter) ([]models.User, error) {
	sb := sqlbuilder.PostgreSQL.NewSelectBuilder()
	sb.Select("*")
	sb.From("users")

	if filter.ID != nil {
		sb.Where(sb.Equal("id", *filter.ID))
	}
	if filter.Email != nil {
		sb.Where(sb.Equal("email", *filter.Email))
	}
	if filter.Name != nil {
		sb.Where(sb.Equal("name", *filter.Name))
	}
	if filter.IsActive != nil {
		sb.Where(sb.Equal("is_active", *filter.IsActive))
	} else {
		sb.Where(sb.Equal("is_active", true))
	}

	query, args := sb.BuildWithFlavor(sqlbuilder.PostgreSQL)

	var users []models.User
	err := r.db.Select(&users, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) UpdateUser(id uuid.UUID, updates map[string]interface{}) (*models.User, error) {
	if len(updates) == 0 {
		return nil, errors.New("no updates provided")
	}

	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("users")

	assignments := []string{}
	updated := false

	if email, ok := updates["email"]; ok {
		assignments = append(assignments, ub.Assign("email", email))
		updated = true
	}
	if passwordHash, ok := updates["password_hash"]; ok {
		assignments = append(assignments, ub.Assign("password_hash", passwordHash))
		updated = true
	}
	if name, ok := updates["name"]; ok {
		fmt.Printf("DEBUG: name value = %#v, type = %T\n", name, name)
		if strVal, isString := name.(string); isString {
			assignments = append(assignments, ub.Assign("name", strVal))
			updated = true
		}
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

	var user models.User
	err := r.db.Get(&user, query, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) DeleteUser(id uuid.UUID) error {
	ub := sqlbuilder.PostgreSQL.NewUpdateBuilder()
	ub.Update("users")
	ub.Set(
		ub.Assign("is_active", false),
		ub.Assign("updated_at", sqlbuilder.Raw("CURRENT_TIMESTAMP")),
	)
	ub.Where(ub.Equal("id", id))
	
	query, args := ub.BuildWithFlavor(sqlbuilder.PostgreSQL)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	fmt.Printf("DeleteUser rows affected: %d\n", rowsAffected)
	if err != nil {
		return fmt.Errorf("error checking rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("user not found")
	}

	return nil
}

