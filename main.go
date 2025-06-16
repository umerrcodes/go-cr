package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq" // PostgreSQL driver
)

// Task represents a task in our system
type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

// Database connection
var db *sql.DB

// Database configuration
const (
	host     = "localhost"
	port     = 5432
	user     = "umer" // your system username
	password = ""     // no password for local development
	dbname   = "taskapi"
)

// Initialize database connection
func initDB() {
	// Build connection string without password since it's empty
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
		host, port, user, dbname)

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Test the connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	log.Println("✅ Successfully connected to PostgreSQL database!")
}

// Middleware for logging requests
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("Started %s %s", r.Method, r.URL.Path)

		next.ServeHTTP(w, r)

		log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// Helper function to set JSON content type
func setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// GET /tasks - Get all tasks from database
func getTasks(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	// SQL query to get all tasks
	rows, err := db.Query("SELECT id, title, description, completed, created_at FROM tasks ORDER BY id")
	if err != nil {
		log.Printf("Error querying tasks: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch tasks"})
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt)
		if err != nil {
			log.Printf("Error scanning task: %v", err)
			continue
		}
		tasks = append(tasks, task)
	}

	json.NewEncoder(w).Encode(tasks)
}

// GET /tasks/{id} - Get a specific task from database
func getTask(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	var task Task
	err = db.QueryRow("SELECT id, title, description, completed, created_at FROM tasks WHERE id = $1", id).
		Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt)

	if err == sql.ErrNoRows {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	} else if err != nil {
		log.Printf("Error querying task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch task"})
		return
	}

	json.NewEncoder(w).Encode(task)
}

// POST /tasks - Create a new task in database
func createTask(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)
	var newTask Task

	// Decode the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	// Validate required fields
	if newTask.Title == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Title is required"})
		return
	}

	// Insert into database and get the generated ID
	err := db.QueryRow(
		"INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id, created_at",
		newTask.Title, newTask.Description, false,
	).Scan(&newTask.ID, &newTask.CreatedAt)

	if err != nil {
		log.Printf("Error creating task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to create task"})
		return
	}

	newTask.Completed = false
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

// PUT /tasks/{id} - Update a task in database
func updateTask(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	var updatedTask Task
	if err := json.NewDecoder(r.Body).Decode(&updatedTask); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid JSON"})
		return
	}

	// Update the task in database
	result, err := db.Exec(
		"UPDATE tasks SET title = $1, description = $2, completed = $3 WHERE id = $4",
		updatedTask.Title, updatedTask.Description, updatedTask.Completed, id,
	)

	if err != nil {
		log.Printf("Error updating task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update task"})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to update task"})
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	// Fetch the updated task to return it
	err = db.QueryRow("SELECT id, title, description, completed, created_at FROM tasks WHERE id = $1", id).
		Scan(&updatedTask.ID, &updatedTask.Title, &updatedTask.Description, &updatedTask.Completed, &updatedTask.CreatedAt)

	if err != nil {
		log.Printf("Error fetching updated task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch updated task"})
		return
	}

	json.NewEncoder(w).Encode(updatedTask)
}

// DELETE /tasks/{id} - Delete a task from database
func deleteTask(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid task ID"})
		return
	}

	// Delete the task from database
	result, err := db.Exec("DELETE FROM tasks WHERE id = $1", id)
	if err != nil {
		log.Printf("Error deleting task: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete task"})
		return
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error checking rows affected: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to delete task"})
		return
	}

	if rowsAffected == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	setJSONContentType(w)

	// Test database connection
	err := db.Ping()
	dbStatus := "connected"
	if err != nil {
		dbStatus = "disconnected"
	}

	response := map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now(),
		"version":   "2.0.0",
		"database":  dbStatus,
	}
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Initialize database connection
	initDB()
	defer db.Close()

	// Create router
	r := mux.NewRouter()

	// Add logging middleware
	r.Use(loggingMiddleware)

	// Define routes
	r.HandleFunc("/health", healthCheck).Methods("GET")
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/tasks", createTask).Methods("POST")
	r.HandleFunc("/tasks/{id}", getTask).Methods("GET")
	r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")

	// Start server
	port := ":8080"
	fmt.Printf("🚀 Server starting on port %s\n", port)
	fmt.Printf("🗄️  Connected to PostgreSQL database: %s\n", dbname)
	fmt.Println("📋 Available endpoints:")
	fmt.Println("  GET    /health       - Health check")
	fmt.Println("  GET    /tasks        - Get all tasks")
	fmt.Println("  POST   /tasks        - Create new task")
	fmt.Println("  GET    /tasks/{id}   - Get specific task")
	fmt.Println("  PUT    /tasks/{id}   - Update task")
	fmt.Println("  DELETE /tasks/{id}   - Delete task")

	log.Fatal(http.ListenAndServe(port, r))
}
