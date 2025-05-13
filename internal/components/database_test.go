package components

import (
	"os"
	"testing"
)

func TestDatabaseConnection(t *testing.T) {
	// Skip test if explicitly disabled
	if os.Getenv("SKIP_DB_TEST") == "true" {
		t.Skip("Skipping database test")
	}

	// Configure test database connection
	host := os.Getenv("TEST_DB_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("TEST_DB_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("TEST_DB_USER")
	if user == "" {
		user = "dashboard"
	}

	password := os.Getenv("TEST_DB_PASSWORD")
	if password == "" {
		password = "changeme"
	}

	dbname := os.Getenv("TEST_DB_NAME")
	if dbname == "" {
		dbname = "dashboard"
	}

	// Initialize database connection
	err := InitDB(host, port, user, password, dbname)
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}
	defer DB.Close()

	// Test database connection with a ping
	if err := DB.Ping(); err != nil {
		t.Fatalf("Database ping failed: %v", err)
	}

	// Verify we can execute a simple query
	var one int
	err = DB.QueryRow("SELECT 1").Scan(&one)
	if err != nil {
		t.Fatalf("Failed to execute simple query: %v", err)
	}

	if one != 1 {
		t.Errorf("Expected 1, got %d", one)
	}
}
