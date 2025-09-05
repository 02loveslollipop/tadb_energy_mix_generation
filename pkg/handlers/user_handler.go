package handlers

import (
	"net/http"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/database"
	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/utils"
	"github.com/gin-gonic/gin"
)

// UserHandler handles HTTP requests for users
type UserHandler struct {
	repo database.Repository
}

// NewUserHandler creates a new UserHandler instance
func NewUserHandler(repo database.Repository) *UserHandler {
	return &UserHandler{
		repo: repo,
	}
}

// GetUserProfile handles GET /user/profile
// @Summary Get user profile
// @Description Get the current user's profile information
// @Tags users
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} models.ErrorResponse
// @Router /user/profile [get]
func (h *UserHandler) GetUserProfile(c *gin.Context) {
	// TODO: Implement user profile retrieval
	// This is a placeholder implementation
	utils.ErrorResponse(c, http.StatusNotImplemented, "Not implemented: User operations not yet implemented")
}

// HealthCheck handles GET /health
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Produce json
// @Success 200 {object} models.SuccessResponse
// @Router /health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API is running",
	})
}
