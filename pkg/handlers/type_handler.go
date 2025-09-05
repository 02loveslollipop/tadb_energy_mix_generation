package handlers

import (
	"net/http"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/models"
	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/utils"
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
// @Success 200 {object} utils.SuccessResponse{data=[]models.Type}
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/types [get]
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
		{
			ID:          uuid.New(),
			Name:        "Natural Gas",
			Description: "Natural gas power plant",
			IsRenewable: false,
		},
	}

	utils.SuccessResponse(c, http.StatusOK, "Types retrieved successfully", types)
}

// GetTypeByID retrieves a specific energy generator type by ID
// @Summary Get energy generator type by ID
// @Description Retrieve a specific energy generator type by its ID
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID"
// @Success 200 {object} utils.SuccessResponse{data=models.Type}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/types/{id} [get]
func (h *TypeHandler) GetTypeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// TODO: Implement database query using core.get_type_by_id()
	// For now, return a mock type
	mockType := models.Type{
		ID:          id,
		Name:        "Solar",
		Description: "Solar photovoltaic panels",
		IsRenewable: true,
	}

	utils.SuccessResponse(c, http.StatusOK, "Type retrieved successfully", mockType)
}

// CreateType creates a new energy generator type
// @Summary Create a new energy generator type
// @Description Create a new energy generator type
// @Tags Types
// @Accept json
// @Produce json
// @Param type body models.TypeCreateRequest true "Type data"
// @Success 201 {object} utils.SuccessResponse{data=models.Type}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/types [post]
func (h *TypeHandler) CreateType(c *gin.Context) {
	var req models.TypeCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// TODO: Implement database insert using core.insert_type()
	newType := models.Type{
		ID:          uuid.New(),
		Name:        req.Name,
		Description: req.Description,
		IsRenewable: req.IsRenewable,
	}

	utils.SuccessResponse(c, http.StatusCreated, "Type created successfully", newType)
}

// UpdateType updates an existing energy generator type
// @Summary Update an energy generator type
// @Description Update an existing energy generator type by ID
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID"
// @Param type body models.TypeUpdateRequest true "Type data"
// @Success 200 {object} utils.SuccessResponse{data=models.Type}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/types/{id} [put]
func (h *TypeHandler) UpdateType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	var req models.TypeUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Handle the boolean pointer properly
	isRenewable := false
	if req.IsRenewable != nil {
		isRenewable = *req.IsRenewable
	}

	// TODO: Implement database update using core.update_type()
	updatedType := models.Type{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IsRenewable: isRenewable,
	}

	utils.SuccessResponse(c, http.StatusOK, "Type updated successfully", updatedType)
}

// DeleteType deletes an energy generator type
// @Summary Delete an energy generator type
// @Description Delete an existing energy generator type by ID
// @Tags Types
// @Accept json
// @Produce json
// @Param id path string true "Type ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /api/v1/types/{id} [delete]
func (h *TypeHandler) DeleteType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid ID format")
		return
	}

	// TODO: Implement database delete using core.delete_type()
	_ = id // Use the id for database operation

	utils.SuccessResponse(c, http.StatusOK, "Type deleted successfully", nil)
}
