package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Import the PostgreSQL driver
	"go.uber.org/zap"
)

const (
	EnvDBConnectionString = "DB_CONNECTION_STRING"
)

// InitDB initializes the database connection
func InitDB(log *zap.Logger) (*sql.DB, error) {
	connectionString := os.Getenv(EnvDBConnectionString)
	if connectionString == "" {
		log.Fatal("Environment variable must be set", zap.String("var", EnvDBConnectionString))
	}

	db, err := sql.Open("postgres", connectionString) // Replace "your-database-driver" with the actual driver name
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50),
		last_name VARCHAR(50)
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create user in the database: %w", err)
	}

	return db, nil
}

// GetDB returns the initialized database connection
func GetDB(log *zap.Logger) (*sql.DB, error) {
	db, err := InitDB(log)
	if err != nil {
		log.Fatal("Failed to initialize database", zap.Error(err))
		return nil, err
	}
	return db, nil
}
