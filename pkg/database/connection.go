package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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

	// Check if DB_URI is provided (preferred method)
	dbURI := strings.TrimSpace(os.Getenv("DB_URI"))
	if dbURI != "" {
		// Parse the URI to get connection details
		return parseDBURI(dbURI)
	}

	// Fallback to individual environment variables
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

	return config, nil
}

// parseDBURI parses a PostgreSQL connection URI and returns a Config
func parseDBURI(dbURI string) (*Config, error) {
	// Parse the connection URI using pgx
	pgxConfig, err := pgx.ParseConfig(dbURI)
	if err != nil {
		return nil, fmt.Errorf("failed to parse DB_URI: %w", err)
	}

	config := &Config{
		Host:            pgxConfig.Host,
		Port:            int(pgxConfig.Port),
		User:            pgxConfig.User,
		Password:        pgxConfig.Password,
		Database:        pgxConfig.Database,
		SSLMode:         "require", // Default for cloud databases
		MaxConnections:  int32(getEnvAsIntWithDefault("DB_MAX_CONNECTIONS", 25)),
		MinConnections:  int32(getEnvAsIntWithDefault("DB_MIN_CONNECTIONS", 5)),
		MaxConnLifetime: time.Duration(getEnvAsIntWithDefault("DB_MAX_CONN_LIFETIME", 60)) * time.Minute,
		MaxConnIdleTime: time.Duration(getEnvAsIntWithDefault("DB_MAX_CONN_IDLE_TIME", 30)) * time.Minute,
	}

	// Override SSL mode if specified in URI
	if sslMode, exists := pgxConfig.RuntimeParams["sslmode"]; exists {
		config.SSLMode = sslMode
	}

	return config, nil
}

// NewConnection creates a new database connection pool
func NewConnection(ctx context.Context) (*DB, error) {
	// Try to load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	var poolConfig *pgxpool.Config
	var err error

	// Check if DB_URI is provided (preferred method)
	dbURI := strings.TrimSpace(os.Getenv("DB_URI"))
	if dbURI != "" {
		// Use DB_URI directly for connection pool
		poolConfig, err = pgxpool.ParseConfig(dbURI)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DB_URI: %w", err)
		}
		// Log chosen path (sanitized)
		if poolConfig != nil && poolConfig.ConnConfig != nil {
			log.Printf("Using DB_URI for connection (user=%s host=%s port=%d db=%s)",
				poolConfig.ConnConfig.User, poolConfig.ConnConfig.Host, poolConfig.ConnConfig.Port, poolConfig.ConnConfig.Database)
		}
	} else {
		// Fallback to individual environment variables
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
		poolConfig, err = pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("failed to parse database configuration: %w", err)
		}
		log.Printf("Using discrete DB_* env vars for connection (user=%s host=%s port=%d db=%s)",
			config.User, config.Host, config.Port, config.Database)
	}

	// Set pool configuration with defaults
	if poolConfig.MaxConns == 0 {
		poolConfig.MaxConns = int32(getEnvAsIntWithDefault("DB_MAX_CONNECTIONS", 25))
	}
	if poolConfig.MinConns == 0 {
		poolConfig.MinConns = int32(getEnvAsIntWithDefault("DB_MIN_CONNECTIONS", 5))
	}
	if poolConfig.MaxConnLifetime == 0 {
		poolConfig.MaxConnLifetime = time.Duration(getEnvAsIntWithDefault("DB_MAX_CONN_LIFETIME", 60)) * time.Minute
	}
	if poolConfig.MaxConnIdleTime == 0 {
		poolConfig.MaxConnIdleTime = time.Duration(getEnvAsIntWithDefault("DB_MAX_CONN_IDLE_TIME", 30)) * time.Minute
	}

	// Set connection parameters for better performance
	if poolConfig.ConnConfig.RuntimeParams == nil {
		poolConfig.ConnConfig.RuntimeParams = make(map[string]string)
	}
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

	// Log connection success
	log.Printf("Successfully connected to PostgreSQL database: %s@%s:%d/%s",
		poolConfig.ConnConfig.User, poolConfig.ConnConfig.Host, poolConfig.ConnConfig.Port, poolConfig.ConnConfig.Database)
	log.Printf("Connection pool configured - Min: %d, Max: %d",
		poolConfig.MinConns, poolConfig.MaxConns)

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
