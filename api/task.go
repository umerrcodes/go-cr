package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// TaskByID endpoint by ID - handles GET, PUT, DELETE /api/task?id=123
func TaskByID(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Get ID from query parameter or URL path
	var idStr string

	// Try query parameter first: /api/task?id=123
	idStr = r.URL.Query().Get("id")

	// If no query parameter, try URL path: /api/task/123
	if idStr == "" {
		path := r.URL.Path
		parts := strings.Split(path, "/")
		if len(parts) >= 3 && parts[2] != "" {
			idStr = parts[2]
		}
	}

	if idStr == "" {
		http.Error(w, "Task ID is required", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Route to appropriate handler
	switch r.Method {
	case http.MethodGet:
		getTaskByID(w, r, id)
	case http.MethodPut:
		updateTaskByID(w, r, id)
	case http.MethodDelete:
		deleteTaskByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /api/task?id=123 - Get specific task
func getTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	// Connect to database
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var task Task
	err = db.QueryRow("SELECT id, title, description, completed, created_at FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch task", http.StatusInternalServerError)
		return
	}

	// Return task as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// PUT /api/task?id=123 - Update task
func updateTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	var updatedTask Task

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Connect to database
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Update task in database
	result, err := db.Exec(
		"UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4",
		updatedTask.Title, updatedTask.Description, updatedTask.Completed, id,
	)

	if err != nil {
		log.Printf("Update error: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
		http.Error(w, "Failed to update task", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Fetch updated task to return it
	err = db.QueryRow("SELECT id, title, description, completed, created_at FROM tasks WHERE id = $1", id).
		Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Completed, &updatedTask.CreatedAt)

	if err != nil {
		log.Printf("Fetch updated task error: %v", err)
		http.Error(w, "Failed to fetch updated task", http.StatusInternalServerError)
		return
	}

	// Return updated task
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTask)
}

// DELETE /api/task?id=123 - Delete task
func deleteTaskByID(w http.ResponseWriter, r *http.Request, id int) {
	// Connect to database
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Delete task from database
	result, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Delete error: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("RowsAffected error: %v", err)
		http.Error(w, "Failed to delete task", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	}

	// Return 204 No Content for successful deletion
	w.WriteHeader(http.StatusNoContent)
}
