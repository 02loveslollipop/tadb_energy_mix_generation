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

    // Generator operations
    CreateGenerator(ctx context.Context, req *models.CreateGeneratorRequest) (*models.Generator, error)
    GetGeneratorByID(ctx context.Context, id uuid.UUID) (*models.Generator, error)
    GetAllGenerators(ctx context.Context, typeID *uuid.UUID) ([]*models.Generator, error)
    UpdateGenerator(ctx context.Context, id uuid.UUID, req *models.UpdateGeneratorRequest) (*models.Generator, error)
    DeleteGenerator(ctx context.Context, id uuid.UUID) error

    // Production operations
    CreateProduction(ctx context.Context, req *models.CreateProductionRequest) (*models.Production, error)
    GetProductionByID(ctx context.Context, id uuid.UUID) (*models.Production, error)
    GetAllProductions(ctx context.Context, generatorID *uuid.UUID, startDate, endDate *string) ([]*models.Production, error)
    UpdateProduction(ctx context.Context, id uuid.UUID, req *models.UpdateProductionRequest) (*models.Production, error)
    DeleteProduction(ctx context.Context, id uuid.UUID) error
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

// Helper to scan Generator with joined fields
func scanGenerator(row pgx.Row, g *models.Generator) error {
    return row.Scan(
        &g.ID,
        &g.TypeID,
        &g.TypeName,
        &g.TypeDesc,
        &g.IsRenewable,
        &g.Capacity,
        &g.CreatedAt,
        &g.UpdatedAt,
    )
}

// Helper to scan Production with joined fields
func scanProduction(row pgx.Row, p *models.Production) error {
    return row.Scan(
        &p.ID,
        &p.GeneratorID,
        &p.GeneratorCapacity,
        &p.TypeName,
        &p.IsRenewable,
        &p.Date,
        &p.ProductionMW,
        &p.CreatedAt,
        &p.UpdatedAt,
    )
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

// ===================== Generators =====================
func (r *postgresRepository) CreateGenerator(ctx context.Context, req *models.CreateGeneratorRequest) (*models.Generator, error) {
    query := `
        INSERT INTO generators (id, type, capacity, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id`
    id := uuid.New()
    now := time.Now()
    if _, err := r.db.Exec(ctx, query, id, req.TypeID, req.Capacity, now, now); err != nil {
        return nil, fmt.Errorf("failed to create generator: %w", err)
    }
    return r.GetGeneratorByID(ctx, id)
}

func (r *postgresRepository) GetGeneratorByID(ctx context.Context, id uuid.UUID) (*models.Generator, error) {
    query := `
        SELECT g.id, g.type, t.name, t.description, t.isrenuevable, g.capacity, g.created_at, g.updated_at
        FROM generators g
        JOIN types t ON g.type = t.id
        WHERE g.id = $1`
    var gen models.Generator
    err := scanGenerator(r.db.QueryRow(ctx, query, id), &gen)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, fmt.Errorf("failed to get generator: %w", err)
    }
    return &gen, nil
}

func (r *postgresRepository) GetAllGenerators(ctx context.Context, typeID *uuid.UUID) ([]*models.Generator, error) {
    var (
        query string
        args []any
    )
    if typeID != nil {
        query = `
            SELECT g.id, g.type, t.name, t.description, t.isrenuevable, g.capacity, g.created_at, g.updated_at
            FROM generators g
            JOIN types t ON g.type = t.id
            WHERE g.type = $1
            ORDER BY t.name, g.capacity DESC`
        args = append(args, *typeID)
    } else {
        query = `
            SELECT g.id, g.type, t.name, t.description, t.isrenuevable, g.capacity, g.created_at, g.updated_at
            FROM generators g
            JOIN types t ON g.type = t.id
            ORDER BY t.name, g.capacity DESC`
    }
    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query generators: %w", err)
    }
    defer rows.Close()
    var list []*models.Generator
    for rows.Next() {
        var g models.Generator
        if err := scanGenerator(rows, &g); err != nil {
            return nil, fmt.Errorf("failed to scan generator: %w", err)
        }
        list = append(list, &g)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("row iteration error: %w", err)
    }
    return list, nil
}

func (r *postgresRepository) UpdateGenerator(ctx context.Context, id uuid.UUID, req *models.UpdateGeneratorRequest) (*models.Generator, error) {
    // Build dynamic update
    // For simplicity, set all fields using COALESCE on provided values
    query := `
        UPDATE generators
        SET type = COALESCE($2, type),
            capacity = COALESCE($3, capacity),
            updated_at = $4
        WHERE id = $1`
    now := time.Now()
    if _, err := r.db.Exec(ctx, query, id, req.TypeID, req.Capacity, now); err != nil {
        if err == pgx.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, fmt.Errorf("failed to update generator: %w", err)
    }
    return r.GetGeneratorByID(ctx, id)
}

func (r *postgresRepository) DeleteGenerator(ctx context.Context, id uuid.UUID) error {
    res, err := r.db.Exec(ctx, `DELETE FROM generators WHERE id = $1`, id)
    if err != nil {
        return fmt.Errorf("failed to delete generator: %w", err)
    }
    if res.RowsAffected() == 0 {
        return sql.ErrNoRows
    }
    return nil
}

// ===================== Productions =====================
func (r *postgresRepository) CreateProduction(ctx context.Context, req *models.CreateProductionRequest) (*models.Production, error) {
    query := `
        INSERT INTO productions (id, generator_id, date, production_mw, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`
    id := uuid.New()
    now := time.Now()
    if _, err := r.db.Exec(ctx, query, id, req.GeneratorID, req.Date, req.ProductionMW, now, now); err != nil {
        return nil, fmt.Errorf("failed to create production: %w", err)
    }
    return r.GetProductionByID(ctx, id)
}

func (r *postgresRepository) GetProductionByID(ctx context.Context, id uuid.UUID) (*models.Production, error) {
    query := `
        SELECT p.id, p.generator_id, g.capacity, t.name, t.isrenuevable, p.date, p.production_mw, p.created_at, p.updated_at
        FROM productions p
        JOIN generators g ON p.generator_id = g.id
        JOIN types t ON g.type = t.id
        WHERE p.id = $1`
    var pr models.Production
    err := scanProduction(r.db.QueryRow(ctx, query, id), &pr)
    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, fmt.Errorf("failed to get production: %w", err)
    }
    return &pr, nil
}

func (r *postgresRepository) GetAllProductions(ctx context.Context, generatorID *uuid.UUID, startDate, endDate *string) ([]*models.Production, error) {
    var (
        query string
        args []any
    )
    base := `
        SELECT p.id, p.generator_id, g.capacity, t.name, t.isrenuevable, p.date, p.production_mw, p.created_at, p.updated_at
        FROM productions p
        JOIN generators g ON p.generator_id = g.id
        JOIN types t ON g.type = t.id`
    where := ""
    idx := 1
    if generatorID != nil {
        where += fmt.Sprintf(" WHERE p.generator_id = $%d", idx)
        args = append(args, *generatorID)
        idx++
    }
    if startDate != nil && *startDate != "" {
        if where == "" { where = " WHERE" } else { where += " AND" }
        where += fmt.Sprintf(" p.date >= $%d", idx)
        args = append(args, *startDate)
        idx++
    }
    if endDate != nil && *endDate != "" {
        if where == "" { where = " WHERE" } else { where += " AND" }
        where += fmt.Sprintf(" p.date <= $%d", idx)
        args = append(args, *endDate)
        idx++
    }
    order := " ORDER BY p.date DESC, t.name"
    query = base + where + order

    rows, err := r.db.Query(ctx, query, args...)
    if err != nil {
        return nil, fmt.Errorf("failed to query productions: %w", err)
    }
    defer rows.Close()
    var list []*models.Production
    for rows.Next() {
        var p models.Production
        if err := scanProduction(rows, &p); err != nil {
            return nil, fmt.Errorf("failed to scan production: %w", err)
        }
        list = append(list, &p)
    }
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("row iteration error: %w", err)
    }
    return list, nil
}

func (r *postgresRepository) UpdateProduction(ctx context.Context, id uuid.UUID, req *models.UpdateProductionRequest) (*models.Production, error) {
    query := `
        UPDATE productions
        SET generator_id = COALESCE($2, generator_id),
            date = COALESCE($3, date),
            production_mw = COALESCE($4, production_mw),
            updated_at = $5
        WHERE id = $1`
    now := time.Now()
    if _, err := r.db.Exec(ctx, query, id, req.GeneratorID, req.Date, req.ProductionMW, now); err != nil {
        if err == pgx.ErrNoRows {
            return nil, sql.ErrNoRows
        }
        return nil, fmt.Errorf("failed to update production: %w", err)
    }
    return r.GetProductionByID(ctx, id)
}

func (r *postgresRepository) DeleteProduction(ctx context.Context, id uuid.UUID) error {
    res, err := r.db.Exec(ctx, `DELETE FROM productions WHERE id = $1`, id)
    if err != nil {
        return fmt.Errorf("failed to delete production: %w", err)
    }
    if res.RowsAffected() == 0 {
        return sql.ErrNoRows
    }
    return nil
}
