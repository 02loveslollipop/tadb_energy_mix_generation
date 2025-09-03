package handlers

import (
	"net/http"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TypeHandler handles energy generator type-related HTTP requests
type TypeHandler struct {
	// Add dependencies like database connection, services, etc.
}

// NewTypeHandler creates a new TypeHandler instance
func NewTypeHandler() *TypeHandler {
	return &TypeHandler{}
}

// GetTypes retrieves all energy generator types
// @Summary Get all energy generator types
// @Description Retrieve a list of all energy generator types
// @Tags Types
// @Accept json
// @Produce json
// @Success 200 {object} models.SuccessResponse{data=[]models.Type}
// @Failure 500 {object} models.ErrorResponse
// @Router /api/types [get]
func (h *TypeHandler) GetTypes(c *gin.Context) {
	// TODO: Implement database query using core.get_all_types()
	types := []models.Type{
		{
			ID:          uuid.New(),
			Name:        "Solar",
			Description: "Solar photovoltaic panels",
			IsRenewable: true,
		},
		{
			ID:          uuid.New(),
			Name:        "Wind",
			Description: "Wind turbines",
			IsRenewable: true,
		},
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Types retrieved successfully",
		Data:    types,
	})
}

// GetTypeByID retrieves a specific energy generator type by ID
// @Summary Get energy generator type by ID
// @Description Retrieve a specific energy generator type by its ID
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID" format(uuid)
// @Success 200 {object} models.SuccessResponse{data=models.Type}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/types/{id} [get]
func (h *TypeHandler) GetTypeByID(c *gin.Context) {
	idParam := c.Param("id")
	typeID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid type ID",
			Message: "The provided ID is not a valid UUID",
			Code:    400,
		})
		return
	}

	// TODO: Implement database query using core.get_type_by_id()
	typeData := models.Type{
		ID:          typeID,
		Name:        "Solar",
		Description: "Solar photovoltaic panels",
		IsRenewable: true,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Type retrieved successfully",
		Data:    typeData,
	})
}

// CreateType creates a new energy generator type
// @Summary Create a new energy generator type
// @Description Create a new energy generator type
// @Tags Types
// @Accept json
// @Produce json
// @Param type body models.CreateTypeRequest true "Type data"
// @Success 201 {object} models.SuccessResponse{data=models.Type}
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/types [post]
func (h *TypeHandler) CreateType(c *gin.Context) {
	var req models.CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
			Code:    400,
		})
		return
	}

	// TODO: Implement database insertion using core.insert_type()
	newType := models.Type{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsRenewable: req.IsRenewable,
	}

	c.JSON(http.StatusCreated, models.SuccessResponse{
		Message: "Type created successfully",
		Data:    newType,
	})
}

// UpdateType updates an existing energy generator type
// @Summary Update an energy generator type
// @Description Update an existing energy generator type by ID
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID" format(uuid)
// @Param type body models.UpdateTypeRequest true "Updated type data"
// @Success 200 {object} models.SuccessResponse{data=models.Type}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/types/{id} [put]
func (h *TypeHandler) UpdateType(c *gin.Context) {
	idParam := c.Param("id")
	typeID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid type ID",
			Message: "The provided ID is not a valid UUID",
			Code:    400,
		})
		return
	}

	var req models.UpdateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid input",
			Message: err.Error(),
			Code:    400,
		})
		return
	}

	// TODO: Implement database update using core.update_type()
	updatedType := models.Type{
		ID:          typeID,
		Name:        req.Name,
		Description: req.Description,
		IsRenewable: *req.IsRenewable,
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Type updated successfully",
		Data:    updatedType,
	})
}

// DeleteType deletes an energy generator type
// @Summary Delete an energy generator type
// @Description Delete an energy generator type by ID
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID" format(uuid)
// @Success 200 {object} models.SuccessResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/types/{id} [delete]
func (h *TypeHandler) DeleteType(c *gin.Context) {
	idParam := c.Param("id")
	typeID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid type ID",
			Message: "The provided ID is not a valid UUID",
			Code:    400,
		})
		return
	}

	// TODO: Implement database deletion using core.delete_type()
	_ = typeID // Use the typeID for database deletion

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Type deleted successfully",
	})
}

// GetGeneratorsByType retrieves all generators of a specific type
// @Summary Get generators by type
// @Description Retrieve all generators of a specific energy generator type
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID" format(uuid)
// @Success 200 {object} models.SuccessResponse{data=[]models.Generator}
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/types/{id}/generators [get]
func (h *TypeHandler) GetGeneratorsByType(c *gin.Context) {
	idParam := c.Param("id")
	typeID, err := uuid.Parse(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error:   "Invalid type ID",
			Message: "The provided ID is not a valid UUID",
			Code:    400,
		})
		return
	}

	// TODO: Implement database query using core.get_generators_by_type()
	generators := []models.Generator{
		{
			ID:       uuid.New(),
			TypeID:   typeID,
			TypeName: "Solar",
			Capacity: 100.5,
		},
	}

	c.JSON(http.StatusOK, models.SuccessResponse{
		Message: "Generators retrieved successfully",
		Data:    generators,
	})
}
