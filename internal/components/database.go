package components

import (
	"database/sql"
	"fmt"
	"time"

	// Import the PostgreSQL driver
	_ "github.com/lib/pq"
)

// DB is a shared database handle
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(host, port, user, password, dbname string) error {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to open database: %w", err)
	}

	// Set connection pool settings
	DB.SetMaxOpenConns(25)
	DB.SetMaxIdleConns(5)
	DB.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	return nil
}

// IsDBHealthy checks if the database connection is healthy
func IsDBHealthy() bool {
	if DB == nil {
		return false
	}
	return DB.Ping() == nil
}
