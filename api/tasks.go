package api

import (
	"encoding/json"
	"log"
	"net/http"
)

// Tasks endpoint - handles GET (all tasks) and POST (create task)
func Tasks(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Route to appropriate handler
	switch r.Method {
	case http.MethodGet:
		getTasks(w, r)
	case http.MethodPost:
		createTask(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /api/tasks - Get all tasks
func getTasks(w http.ResponseWriter, r *http.Request) {
	// Connect to database
	db, err := ConnectDB()
	if err != nil {
		log.Printf("Database connection error: %v", err)
		http.Error(w, "Database connection failed", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Execute query
	rows, err := db.Query("SELECT id, title, description, completed, created_at FROM tasks ORDER BY id")
	if err != nil {
		log.Printf("Query error: %v", err)
		http.Error(w, "Failed to fetch tasks", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Parse results
	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt)
		if err != nil {
			log.Printf("Scan error: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}

	// Return JSON response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// POST /api/tasks - Create new task
func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task

	// Parse JSON request body
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if newTask.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
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

	// Insert new task
	err = db.QueryRow(
		"INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id, created_at",
		newTask.Title, newTask.Description, false,
	).Scan(&newTask.ID, &newTask.CreatedAt)

	if err != nil {
		log.Printf("Insert error: %v", err)
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	// Set completed to false explicitly
	newTask.Completed = false

	// Return created task
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}
