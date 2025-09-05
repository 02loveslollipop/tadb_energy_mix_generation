package handlers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/database"
	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/models"
	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// TypeHandler handles HTTP requests for energy generator types
type TypeHandler struct {
	repo database.Repository
}

// NewTypeHandler creates a new TypeHandler instance
func NewTypeHandler(repo database.Repository) *TypeHandler {
	return &TypeHandler{
		repo: repo,
	}
}

// CreateType handles POST /types
// @Summary Create a new energy generator type
// @Description Create a new energy generator type (renewable/non-renewable)
// @Tags types
// @Accept json
// @Produce json
// @Param type body models.CreateTypeRequest true "Type data"
// @Success 201 {object} models.Type
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /types [post]
func (h *TypeHandler) CreateType(c *gin.Context) {
	var req models.CreateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	typeRecord, err := h.repo.CreateType(c.Request.Context(), &req)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create type: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, typeRecord)
}

// GetTypeByID handles GET /types/:id
// @Summary Get type by ID
// @Description Get an energy generator type by its UUID
// @Tags types
// @Produce json
// @Param id path string true "Type ID (UUID)"
// @Success 200 {object} models.Type
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /types/{id} [get]
func (h *TypeHandler) GetTypeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid type ID: ID must be a valid UUID")
		return
	}

	typeRecord, err := h.repo.GetTypeByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(c, http.StatusNotFound, "Type not found: No type found with the given ID")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get type: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, typeRecord)
}

// GetAllTypes handles GET /types
// @Summary Get all types
// @Description Get all energy generator types, optionally filtered by renewable status
// @Tags types
// @Produce json
// @Param renewable query boolean false "Filter by renewable status (true/false)"
// @Success 200 {array} models.Type
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /types [get]
func (h *TypeHandler) GetAllTypes(c *gin.Context) {
	var isRenewable *bool

	if renewableParam := c.Query("renewable"); renewableParam != "" {
		renewable, err := strconv.ParseBool(renewableParam)
		if err != nil {
			utils.ErrorResponse(c, http.StatusBadRequest, "Invalid renewable parameter: renewable must be true or false")
			return
		}
		isRenewable = &renewable
	}

	types, err := h.repo.GetAllTypes(c.Request.Context(), isRenewable)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get types: "+err.Error())
		return
	}

	if types == nil {
		types = []*models.Type{}
	}

	c.JSON(http.StatusOK, types)
}

// UpdateType handles PUT /types/:id
// @Summary Update type
// @Description Update an existing energy generator type
// @Tags types
// @Accept json
// @Produce json
// @Param id path string true "Type ID (UUID)"
// @Param type body models.UpdateTypeRequest true "Updated type data"
// @Success 200 {object} models.Type
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /types/{id} [put]
func (h *TypeHandler) UpdateType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid type ID: ID must be a valid UUID")
		return
	}

	var req models.UpdateTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
		return
	}

	typeRecord, err := h.repo.UpdateType(c.Request.Context(), id, &req)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(c, http.StatusNotFound, "Type not found: No type found with the given ID")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update type: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, typeRecord)
}

// DeleteType handles DELETE /types/:id
// @Summary Delete type
// @Description Delete an energy generator type by ID
// @Tags types
// @Produce json
// @Param id path string true "Type ID (UUID)"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /types/{id} [delete]
func (h *TypeHandler) DeleteType(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid type ID: ID must be a valid UUID")
		return
	}

	err = h.repo.DeleteType(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(c, http.StatusNotFound, "Type not found: No type found with the given ID")
			return
		}
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete type: "+err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}
