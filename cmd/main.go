package main

import (
	"log"
	"net/http"
	"os"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Initialize handlers
	typeHandler := handlers.NewTypeHandler()

	// Define basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to TADB Energy Matrix API",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Health check endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "healthy",
		})
	})

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Type routes (Energy Generator Types)
		types := v1.Group("/types")
		{
			types.GET("", typeHandler.GetTypes)
			types.GET("/:id", typeHandler.GetTypeByID)
			types.POST("", typeHandler.CreateType)
			types.PUT("/:id", typeHandler.UpdateType)
			types.DELETE("/:id", typeHandler.DeleteType)
		}
	}

	// Get port from environment variable or default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start the server
	log.Printf("Starting TADB Energy Matrix API server on :%s", port)
	log.Println("Available endpoints:")
	log.Println("  GET    /")
	log.Println("  GET    /health")
	log.Println("  GET    /api/v1/types")
	log.Println("  GET    /api/v1/types/:id")
	log.Println("  POST   /api/v1/types")
	log.Println("  PUT    /api/v1/types/:id")
	log.Println("  DELETE /api/v1/types/:id")

	if err := r.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
