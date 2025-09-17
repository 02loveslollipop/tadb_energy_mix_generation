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

type ProductionHandler struct {
    repo database.Repository
}

func NewProductionHandler(repo database.Repository) *ProductionHandler {
    return &ProductionHandler{repo: repo}
}

// CreateProduction handles POST /productions
// @Summary Create production record
// @Tags productions
// @Accept json
// @Produce json
// @Param body body models.CreateProductionRequest true "Production data"
// @Success 201 {object} models.Production
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /productions [post]
func (h *ProductionHandler) CreateProduction(c *gin.Context) {
    var req models.CreateProductionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }
    pr, err := h.repo.CreateProduction(c.Request.Context(), &req)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create production: "+err.Error())
        return
    }
    c.JSON(http.StatusCreated, pr)
}

// GetProductionByID handles GET /productions/:id
// @Summary Get production by ID
// @Tags productions
// @Produce json
// @Param id path string true "Production ID"
// @Success 200 {object} models.Production
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /productions/{id} [get]
func (h *ProductionHandler) GetProductionByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid production ID: must be UUID")
        return
    }
    pr, err := h.repo.GetProductionByID(c.Request.Context(), id)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.ErrorResponse(c, http.StatusNotFound, "Production not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to get production: "+err.Error())
        return
    }
    c.JSON(http.StatusOK, pr)
}

// GetAllProductions handles GET /productions with mixed search
// @Summary List productions (filter by generator/date range)
// @Description List all productions, optionally filtered by generatorId and startDate/endDate (YYYY-MM-DD)
// @Tags productions
// @Produce json
// @Param generatorId query string false "Generator ID (UUID)"
// @Param startDate query string false "Start date (YYYY-MM-DD)"
// @Param endDate query string false "End date (YYYY-MM-DD)"
// @Success 200 {array} models.Production
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /productions [get]
func (h *ProductionHandler) GetAllProductions(c *gin.Context) {
    var genID *uuid.UUID
    if g := c.Query("generatorId"); g != "" {
        id, err := uuid.Parse(g)
        if err != nil {
            utils.ErrorResponse(c, http.StatusBadRequest, "Invalid generatorId: must be UUID")
            return
        }
        genID = &id
    }
    var start, end *string
    if s := c.Query("startDate"); s != "" {
        start = &s
    }
    if e := c.Query("endDate"); e != "" {
        end = &e
    }
    list, err := h.repo.GetAllProductions(c.Request.Context(), genID, start, end)
    if err != nil {
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to list productions: "+err.Error())
        return
    }
    if list == nil { list = []*models.Production{} }
    c.JSON(http.StatusOK, list)
}

// UpdateProduction handles PUT /productions/:id
// @Summary Update production
// @Tags productions
// @Accept json
// @Produce json
// @Param id path string true "Production ID"
// @Param body body models.UpdateProductionRequest true "Update data"
// @Success 200 {object} models.Production
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /productions/{id} [put]
func (h *ProductionHandler) UpdateProduction(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid production ID: must be UUID")
        return
    }
    var req models.UpdateProductionRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request body: "+err.Error())
        return
    }
    pr, err := h.repo.UpdateProduction(c.Request.Context(), id, &req)
    if err != nil {
        if err == sql.ErrNoRows {
            utils.ErrorResponse(c, http.StatusNotFound, "Production not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to update production: "+err.Error())
        return
    }
    c.JSON(http.StatusOK, pr)
}

// DeleteProduction handles DELETE /productions/:id
// @Summary Delete production
// @Tags productions
// @Produce json
// @Param id path string true "Production ID"
// @Success 204
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /productions/{id} [delete]
func (h *ProductionHandler) DeleteProduction(c *gin.Context) {
    idStr := c.Param("id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        utils.ErrorResponse(c, http.StatusBadRequest, "Invalid production ID: must be UUID")
        return
    }
    if err := h.repo.DeleteProduction(c.Request.Context(), id); err != nil {
        if err == sql.ErrNoRows {
            utils.ErrorResponse(c, http.StatusNotFound, "Production not found")
            return
        }
        utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to delete production: "+err.Error())
        return
    }
    c.Status(http.StatusNoContent)
}

