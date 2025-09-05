package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// DB holds the database connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// Config represents database configuration
type Config struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxConnections  int32
	MinConnections  int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
}

// LoadConfig loads database configuration from environment variables
func LoadConfig() (*Config, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	config := &Config{
		Host:            getEnvWithDefault("DB_HOST", "localhost"),
		Port:            getEnvAsIntWithDefault("DB_PORT", 5432),
		User:            getEnvWithDefault("DB_USER", "postgres"),
		Password:        getEnvWithDefault("DB_PASSWORD", ""),
		Database:        getEnvWithDefault("DB_NAME", "tadb"),
		SSLMode:         getEnvWithDefault("DB_SSL_MODE", "disable"),
		MaxConnections:  int32(getEnvAsIntWithDefault("DB_MAX_CONNECTIONS", 25)),
		MinConnections:  int32(getEnvAsIntWithDefault("DB_MIN_CONNECTIONS", 5)),
		MaxConnLifetime: time.Duration(getEnvAsIntWithDefault("DB_MAX_CONN_LIFETIME", 60)) * time.Minute,
		MaxConnIdleTime: time.Duration(getEnvAsIntWithDefault("DB_MAX_CONN_IDLE_TIME", 30)) * time.Minute,
	}

	// Validate required configuration
	if config.Password == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	return config, nil
}

// NewConnection creates a new database connection pool
func NewConnection(ctx context.Context) (*DB, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load database configuration: %w", err)
	}

	// Build connection string
	dsn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Host,
		config.Port,
		config.User,
		config.Password,
		config.Database,
		config.SSLMode,
	)

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database configuration: %w", err)
	}

	// Set pool configuration
	poolConfig.MaxConns = config.MaxConnections
	poolConfig.MinConns = config.MinConnections
	poolConfig.MaxConnLifetime = config.MaxConnLifetime
	poolConfig.MaxConnIdleTime = config.MaxConnIdleTime

	// Set connection parameters for better performance
	poolConfig.ConnConfig.RuntimeParams["application_name"] = "tadb-api"
	poolConfig.ConnConfig.RuntimeParams["search_path"] = "core,public"

	// Create connection pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database: %s@%s:%d/%s",
		config.User, config.Host, config.Port, config.Database)
	log.Printf("Connection pool configured - Min: %d, Max: %d",
		config.MinConnections, config.MaxConnections)

	return &DB{Pool: pool}, nil
}

// Close closes the database connection pool
func (db *DB) Close() {
	if db.Pool != nil {
		db.Pool.Close()
		log.Println("Database connection pool closed")
	}
}

// GetStats returns connection pool statistics
func (db *DB) GetStats() *pgxpool.Stat {
	if db.Pool == nil {
		return nil
	}
	return db.Pool.Stat()
}

// Health checks the database connection health
func (db *DB) Health(ctx context.Context) error {
	if db.Pool == nil {
		return fmt.Errorf("database connection pool is nil")
	}

	// Test connection with a simple query
	var result int
	err := db.Pool.QueryRow(ctx, "SELECT 1").Scan(&result)
	if err != nil {
		return fmt.Errorf("database health check failed: %w", err)
	}

	return nil
}

func (db *DB) BeginTransaction(ctx context.Context) (pgx.Tx, error) {
	if db.Pool == nil {
		return nil, fmt.Errorf("database connection pool is nil")
	}

	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}

	return tx, nil
}


// Helper functions for environment variables
func getEnvWithDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsIntWithDefault(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// Global database instance (singleton pattern)
var globalDB *DB

// Initialize initializes the global database connection
func Initialize(ctx context.Context) error {
	db, err := NewConnection(ctx)
	if err != nil {
		return err
	}
	globalDB = db
	return nil
}

// GetDB returns the global database instance
func GetDB() *DB {
	return globalDB
}

// CloseDB closes the global database connection
func CloseDB() {
	if globalDB != nil {
		globalDB.Close()
		globalDB = nil
	}
}
