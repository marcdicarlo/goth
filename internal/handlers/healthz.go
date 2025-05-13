package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"goth/internal/components"
)

// HealthStatus represents the application's health status
type HealthStatus struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Database  DBStatus  `json:"database"`
}

// DBStatus represents the database connection status
type DBStatus struct {
	Connected bool   `json:"connected"`
	Message   string `json:"message,omitempty"`
}

// HealthzHandler returns the health status of the application
func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	// Create health status response
	health := HealthStatus{
		Status:    "OK",
		Timestamp: time.Now(),
		Database: DBStatus{
			Connected: false,
			Message:   "Database not initialized",
		},
	}

	// Check database connection
	if components.DB != nil {
		if err := components.DB.Ping(); err != nil {
			health.Database.Connected = false
			health.Database.Message = "Database connection failed: " + err.Error()
			health.Status = "WARNING" // Service is running but DB is not available
		} else {
			health.Database.Connected = true
			health.Database.Message = "Connected"
		}
	}

	// Set content type and encode response as JSON
	w.Header().Set("Content-Type", "application/json")

	// Determine HTTP status code based on health status
	if health.Status != "OK" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	// Encode and send response
	json.NewEncoder(w).Encode(health)
}
