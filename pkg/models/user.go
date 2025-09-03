package models

import (
	"time"

	"github.com/google/uuid"
)

// Type represents an energy generator type
// @Description Energy generator type (renewable/non-renewable)
type Type struct {
	ID          uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name        string    `json:"name" db:"name" binding:"required,max=20" example:"Solar"`
	Description string    `json:"description" db:"description" binding:"required,max=80" example:"Solar photovoltaic panels"`
	IsRenewable bool      `json:"isRenewable" db:"isrenuevable" example:"true"`
	CreatedAt   time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// CreateTypeRequest represents the request payload for creating a type
// @Description Request body for creating a new energy generator type
type CreateTypeRequest struct {
	Name        string `json:"name" binding:"required,max=20" example:"Solar"`
	Description string `json:"description" binding:"required,max=80" example:"Solar photovoltaic panels"`
	IsRenewable bool   `json:"isRenewable" example:"true"`
}

// UpdateTypeRequest represents the request payload for updating a type
// @Description Request body for updating an energy generator type
type UpdateTypeRequest struct {
	Name        string `json:"name,omitempty" binding:"omitempty,max=20" example:"Solar"`
	Description string `json:"description,omitempty" binding:"omitempty,max=80" example:"Solar photovoltaic panels"`
	IsRenewable *bool  `json:"isRenewable,omitempty" example:"true"`
}

// Generator represents an energy generator
// @Description Energy generator with capacity and type information
type Generator struct {
	ID          uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440001"`
	TypeID      uuid.UUID `json:"typeId" db:"type" example:"550e8400-e29b-41d4-a716-446655440000"`
	TypeName    string    `json:"typeName,omitempty" db:"type_name" example:"Solar"`
	TypeDesc    string    `json:"typeDescription,omitempty" db:"type_description" example:"Solar photovoltaic panels"`
	IsRenewable bool      `json:"isRenewable,omitempty" db:"isrenuevable" example:"true"`
	Capacity    float64   `json:"capacity" db:"capacity" binding:"required,gt=0" example:"100.5"`
	CreatedAt   time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// CreateGeneratorRequest represents the request payload for creating a generator
// @Description Request body for creating a new energy generator
type CreateGeneratorRequest struct {
	TypeID   uuid.UUID `json:"typeId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Capacity float64   `json:"capacity" binding:"required,gt=0" example:"100.5"`
}

// UpdateGeneratorRequest represents the request payload for updating a generator
// @Description Request body for updating an energy generator
type UpdateGeneratorRequest struct {
	TypeID   *uuid.UUID `json:"typeId,omitempty" example:"550e8400-e29b-41d4-a716-446655440000"`
	Capacity *float64   `json:"capacity,omitempty" binding:"omitempty,gt=0" example:"100.5"`
}

// Production represents energy production data
// @Description Daily energy production record for a generator
type Production struct {
	ID                uuid.UUID `json:"id" db:"id" example:"550e8400-e29b-41d4-a716-446655440002"`
	GeneratorID       uuid.UUID `json:"generatorId" db:"generator_id" example:"550e8400-e29b-41d4-a716-446655440001"`
	GeneratorCapacity float64   `json:"generatorCapacity,omitempty" db:"generator_capacity" example:"100.5"`
	TypeName          string    `json:"typeName,omitempty" db:"type_name" example:"Solar"`
	IsRenewable       bool      `json:"isRenewable,omitempty" db:"isrenuevable" example:"true"`
	Date              string    `json:"date" db:"date" binding:"required" example:"2025-09-03"`
	ProductionMW      float64   `json:"productionMw" db:"production_mw" binding:"required,gte=0" example:"85.3"`
	CreatedAt         time.Time `json:"createdAt,omitempty" db:"created_at"`
	UpdatedAt         time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}

// CreateProductionRequest represents the request payload for creating a production record
// @Description Request body for creating a new production record
type CreateProductionRequest struct {
	GeneratorID  uuid.UUID `json:"generatorId" binding:"required" example:"550e8400-e29b-41d4-a716-446655440001"`
	Date         string    `json:"date" binding:"required" example:"2025-09-03"`
	ProductionMW float64   `json:"productionMw" binding:"required,gte=0" example:"85.3"`
}

// UpdateProductionRequest represents the request payload for updating a production record
// @Description Request body for updating a production record
type UpdateProductionRequest struct {
	GeneratorID  *uuid.UUID `json:"generatorId,omitempty" example:"550e8400-e29b-41d4-a716-446655440001"`
	Date         *string    `json:"date,omitempty" example:"2025-09-03"`
	ProductionMW *float64   `json:"productionMw,omitempty" binding:"omitempty,gte=0" example:"85.3"`
}

// ErrorResponse represents an error response
// @Description Error response structure
type ErrorResponse struct {
	Error   string `json:"error" example:"Invalid input"`
	Message string `json:"message,omitempty" example:"The provided data is invalid"`
	Code    int    `json:"code,omitempty" example:"400"`
}

// SuccessResponse represents a success response
// @Description Success response structure
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// Analytics structures for reporting endpoints

// TotalProductionByDate represents daily production totals
// @Description Daily production summary with renewable breakdown
type TotalProductionByDate struct {
	Date                   string  `json:"date" example:"2025-09-03"`
	TotalProduction        float64 `json:"totalProduction" example:"1250.5"`
	RenewableProduction    float64 `json:"renewableProduction" example:"850.3"`
	NonRenewableProduction float64 `json:"nonRenewableProduction" example:"400.2"`
}

// GeneratorEfficiency represents generator performance metrics
// @Description Generator efficiency and performance data
type GeneratorEfficiency struct {
	GeneratorID          uuid.UUID `json:"generatorId" example:"550e8400-e29b-41d4-a716-446655440001"`
	TypeName             string    `json:"typeName" example:"Solar"`
	Capacity             float64   `json:"capacity" example:"100.5"`
	TotalProduction      float64   `json:"totalProduction" example:"2850.7"`
	AvgDailyProduction   float64   `json:"avgDailyProduction" example:"85.3"`
	EfficiencyPercentage float64   `json:"efficiencyPercentage" example:"84.87"`
}

// RenewableSummary represents renewable vs non-renewable summary
// @Description Summary of renewable vs non-renewable energy production
type RenewableSummary struct {
	EnergyType        string  `json:"energyType" example:"Renewable"`
	TotalCapacity     float64 `json:"totalCapacity" example:"500.0"`
	GeneratorCount    int64   `json:"generatorCount" example:"5"`
	TotalProduction   float64 `json:"totalProduction" example:"12750.5"`
	AvgProduction     float64 `json:"avgProduction" example:"85.0"`
	PercentageOfTotal float64 `json:"percentageOfTotal" example:"68.5"`
}
