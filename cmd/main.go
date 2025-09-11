package main

//
//  Swagger/OpenAPI meta
//  @title           TADB Energy Matrix API
//  @version         1.0.0
//  @description     The TADB (Tabla de Análisis de Datos de Boletines) Energy Matrix API provides endpoints for managing energy generator types and basic health checks.
//  @contact.name    TADB API Support
//  @contact.url     https://github.com/02loveslollipop/api_matriz_enegertica_tadb
//  @license.name    MIT
//  @license.url     https://opensource.org/licenses/MIT
//  @schemes         http https
//  @BasePath        /api/v1

import (
    "context"
    "log"
    "net/http"

    "github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/database"
    "github.com/02loveslollipop/api_matriz_enegertica_tadb/pkg/handlers"
    "github.com/gin-gonic/gin"

    // Swagger UI
    swaggerFiles "github.com/swaggo/files"
    ginSwagger "github.com/swaggo/gin-swagger"
    // Generated docs (created by `swag init`)
    _ "github.com/02loveslollipop/api_matriz_enegertica_tadb/docs"
)

func main() {
	// Initialize database connection
	ctx := context.Background()
	db, err := database.NewConnection(ctx)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create repository
	repo := database.NewRepository(db.Pool)

	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Initialize handlers
	userHandler := handlers.NewUserHandler(repo)
	typeHandler := handlers.NewTypeHandler(repo)

	// Define basic routes
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Welcome to TADB API",
			"status":  "running",
			"version": "1.0.0",
		})
	})

	// Health check endpoint
	r.GET("/health", userHandler.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Type routes
		types := v1.Group("/types")
		{
			types.GET("", typeHandler.GetAllTypes)
			types.GET("/:id", typeHandler.GetTypeByID)
			types.POST("", typeHandler.CreateType)
			types.PUT("/:id", typeHandler.UpdateType)
			types.DELETE("/:id", typeHandler.DeleteType)
		}

		// User routes (placeholder)
		users := v1.Group("/users")
		{
			users.GET("/profile", userHandler.GetUserProfile)
		}
	}

	// Start the server on port 8080
	log.Println("Starting TADB API server on :8080")
	log.Println("Available endpoints:")
	log.Println("  GET  /")
	log.Println("  GET  /health")
	log.Println("  GET  /api/v1/types")
	log.Println("  POST /api/v1/types")
	log.Println("  GET  /api/v1/types/:id")
	log.Println("  PUT  /api/v1/types/:id")
	log.Println("  DELETE /api/v1/types/:id")
	log.Println("  GET  /api/v1/users/profile")

    // Swagger UI endpoint
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
