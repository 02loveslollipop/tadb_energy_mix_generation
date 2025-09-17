package handlers

import (
    "database/sql"
    "net/http"

    "github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/database"
    "github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/models"
    "github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/utils"
    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
)

type GeneratorHandler struct {
    repo database.Repository
}

func NewGeneratorHandler(repo database.Repository) *GeneratorHandler {
    return &GeneratorHandler{repo: repo}
}

// CreateGenerator handles POST /generators
// @Summary Create generator
// @Description Create a new energy generator
// @Tags generators
// @Accept json
// @Produce json
// @Param body body models.CreateGeneratorRequest true "Generator data"
// @Success 201 {object} models.Generator
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /generators [post]
func (h *GeneratorHandler) CreateGenerator(c *gin.Context) {
    var req models.CreateGeneratorRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }
    gen, err := h.repo.CreateGenerator(c.Request.Context(), &req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create generator: "+err.Error())
        return
    }
    c.JSON(http.StatusCreated, gen)
}

// GetGeneratorByID handles GET /generators/:id
// @Summary Get generator by ID
// @Tags generators
// @Produce json
// @Param id path string true "Generator ID"
// @Success 200 {object} models.Generator
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /generators/{id} [get]
func (h *GeneratorHandler) GetGeneratorByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid generator ID: must be UUID")
        return
    }
    gen, err := h.repo.GetGeneratorByID(c.Request.Context(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.ErrorResponse(c, http.StatusNotFound, "Generator not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get generator: "+err.Error())
        return
    }
    c.JSON(http.StatusOK, gen)
}

// GetAllGenerators handles GET /generators
// @Summary List generators
// @Description List all generators, optionally filtered by typeId
// @Tags generators
// @Produce json
// @Param typeId query string false "Type ID (UUID)"
// @Success 200 {array} models.Generator
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /generators [get]
func (h *GeneratorHandler) GetAllGenerators(c *gin.Context) {
    var typeID *uuid.UUID
    if t := c.Query("typeId"); t != "" {
        id, err := uuid.Parse(t)
        if err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, "Invalid typeId: must be UUID")
            return
        }
        typeID = &id
    }
    list, err := h.repo.GetAllGenerators(c.Request.Context(), typeID)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list generators: "+err.Error())
        return
    }
    if list == nil { list = []*models.Generator{} }
    c.JSON(http.StatusOK, list)
}

// UpdateGenerator handles PUT /generators/:id
// @Summary Update generator
// @Tags generators
// @Accept json
// @Produce json
// @Param id path string true "Generator ID"
// @Param body body models.UpdateGeneratorRequest true "Update data"
// @Success 200 {object} models.Generator
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /generators/{id} [put]
func (h *GeneratorHandler) UpdateGenerator(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid generator ID: must be UUID")
        return
    }
    var req models.UpdateGeneratorRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }
    gen, err := h.repo.UpdateGenerator(c.Request.Context(), id, &req)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.ErrorResponse(c, http.StatusNotFound, "Generator not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update generator: "+err.Error())
        return
    }
    c.JSON(http.StatusOK, gen)
}

// DeleteGenerator handles DELETE /generators/:id
// @Summary Delete generator
// @Tags generators
// @Produce json
// @Param id path string true "Generator ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /generators/{id} [delete]
func (h *GeneratorHandler) DeleteGenerator(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid generator ID: must be UUID")
        return
    }
    if err := h.repo.DeleteGenerator(c.Request.Context(), id); err != nil {
        if err == sql.ErrNoRows {
            utils.ErrorResponse(c, http.StatusNotFound, "Generator not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete generator: "+err.Error())
        return
    }
    c.Status(http.StatusNoContent)
}

