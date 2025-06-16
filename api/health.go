package api

import (
	"encoding/json"
	"net/http"
	"time"
)

// Health check endpoint
func Health(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Test database connection
	db, err := ConnectDB()
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
	} else {
		db.Close()
	}

	// Create response
	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "2.0.0-serverless",
		"database":  dbStatus,
	}

	// Set content type and return JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
