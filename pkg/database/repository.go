package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository interface defines all database operations
type Repository interface {
	// Type operations
	CreateType(ctx context.Context, req *models.CreateTypeRequest) (*models.Type, error)
	GetTypeByID(ctx context.Context, id uuid.UUID) (*models.Type, error)
	GetAllTypes(ctx context.Context, isRenewable *bool) ([]*models.Type, error)
	UpdateType(ctx context.Context, id uuid.UUID, req *models.UpdateTypeRequest) (*models.Type, error)
	DeleteType(ctx context.Context, id uuid.UUID) error

	// User operations (placeholder for future implementation)
	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
}

// postgresRepository implements Repository interface
type postgresRepository struct {
	db *pgxpool.Pool
}

// NewRepository creates a new repository instance
func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{
		db: db,
	}
}

// CreateType creates a new energy generator type
func (r *postgresRepository) CreateType(ctx context.Context, req *models.CreateTypeRequest) (*models.Type, error) {
	query := `
		INSERT INTO types (id, name, description, isrenuevable, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, isrenuevable, created_at, updated_at`

	id := uuid.New()
	now := time.Now()

	var typeRecord models.Type
	err := r.db.QueryRow(ctx, query, id, req.Name, req.Description, req.IsRenewable, now, now).Scan(
		&typeRecord.ID,
		&typeRecord.Name,
		&typeRecord.Description,
		&typeRecord.IsRenewable,
		&typeRecord.CreatedAt,
		&typeRecord.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create type: %w", err)
	}

	return &typeRecord, nil
}

// GetTypeByID retrieves a type by its ID
func (r *postgresRepository) GetTypeByID(ctx context.Context, id uuid.UUID) (*models.Type, error) {
	query := `
		SELECT id, name, description, isrenuevable, created_at, updated_at
		FROM types
		WHERE id = $1`

	var typeRecord models.Type
	err := r.db.QueryRow(ctx, query, id).Scan(
		&typeRecord.ID,
		&typeRecord.Name,
		&typeRecord.Description,
		&typeRecord.IsRenewable,
		&typeRecord.CreatedAt,
		&typeRecord.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to get type: %w", err)
	}

	return &typeRecord, nil
}

// GetAllTypes retrieves all types, optionally filtered by renewable status
func (r *postgresRepository) GetAllTypes(ctx context.Context, isRenewable *bool) ([]*models.Type, error) {
	var query string
	var args []interface{}

	if isRenewable != nil {
		query = `
			SELECT id, name, description, isrenuevable, created_at, updated_at
			FROM types
			WHERE isrenuevable = $1
			ORDER BY name`
		args = append(args, *isRenewable)
	} else {
		query = `
			SELECT id, name, description, isrenuevable, created_at, updated_at
			FROM types
			ORDER BY name`
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query types: %w", err)
	}
	defer rows.Close()

	var types []*models.Type
	for rows.Next() {
		var typeRecord models.Type
		err := rows.Scan(
			&typeRecord.ID,
			&typeRecord.Name,
			&typeRecord.Description,
			&typeRecord.IsRenewable,
			&typeRecord.CreatedAt,
			&typeRecord.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan type: %w", err)
		}
		types = append(types, &typeRecord)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return types, nil
}

// UpdateType updates an existing type
func (r *postgresRepository) UpdateType(ctx context.Context, id uuid.UUID, req *models.UpdateTypeRequest) (*models.Type, error) {
	query := `
		UPDATE types
		SET name = $2, description = $3, isrenuevable = $4, updated_at = $5
		WHERE id = $1
		RETURNING id, name, description, isrenuevable, created_at, updated_at`

	now := time.Now()

	var typeRecord models.Type
	err := r.db.QueryRow(ctx, query, id, req.Name, req.Description, req.IsRenewable, now).Scan(
		&typeRecord.ID,
		&typeRecord.Name,
		&typeRecord.Description,
		&typeRecord.IsRenewable,
		&typeRecord.CreatedAt,
		&typeRecord.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, fmt.Errorf("failed to update type: %w", err)
	}

	return &typeRecord, nil
}

// DeleteType deletes a type by its ID
func (r *postgresRepository) DeleteType(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM types WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete type: %w", err)
	}

	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetUserByID is a placeholder implementation
func (r *postgresRepository) GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	// TODO: Implement user operations when User model is ready
	return nil, fmt.Errorf("user operations not implemented yet")
}
