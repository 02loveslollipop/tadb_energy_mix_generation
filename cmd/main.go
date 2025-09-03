package main

import (
	"log"
	"net/http"

	"github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/handlers"
	"github.com/gin-gonic/gin"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Initialize handlers
	userHandler := handlers.NewUserHandler()

	// Define basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to TADB API",
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
		// User routes
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
		}
	}

	// Start the server on port 8080
	log.Println("Starting TADB API server on :8080")
	log.Println("Available endpoints:")
	log.Println("  GET    /")
	log.Println("  GET    /health")
	log.Println("  GET    /api/v1/users")
	log.Println("  GET    /api/v1/users/:id")
	log.Println("  POST   /api/v1/users")
	log.Println("  PUT    /api/v1/users/:id")
	log.Println("  DELETE /api/v1/users/:id")

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
