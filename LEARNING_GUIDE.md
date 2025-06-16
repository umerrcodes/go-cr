# Complete Go Backend + PostgreSQL Learning Guide 📚

## 🎯 **What We Built: Complete Architecture**

### **Simple Overview:**
```
HTTP Request → Go API Server → PostgreSQL Database → Persistent Storage
```

### **Detailed Flow:**
1. **Client** (curl/browser) sends HTTP request
2. **Gorilla Mux Router** matches URL to handler function
3. **Go Handler** processes request, validates data
4. **PostgreSQL Driver** executes SQL queries
5. **Database** stores/retrieves data persistently
6. **JSON Response** sent back to client

---

## 📋 **Step-by-Step What We Did & Why**

### **Phase 1: Environment Setup**

#### **Step 1: Install Go**
```bash
brew install go
```
**What:** Programming language for backend development
**Why:** Fast compilation, great for web services, used by Google/Uber/Netflix

#### **Step 2: Install PostgreSQL**
```bash
brew install postgresql@15
brew services start postgresql@15
```
**What:** Production-grade relational database
**Why:** ACID compliance, handles millions of records, supports complex queries

#### **Step 3: Create Project Structure**
```bash
go mod init task-api
```
**What:** Initializes Go module (like package.json in Node.js)
**Why:** Manages dependencies and project metadata

### **Phase 2: Database Design**

#### **Step 4: Create Database**
```bash
createdb taskapi
```
**What:** Creates a new database named "taskapi"
**Why:** Separates our project data from other applications

#### **Step 5: Design Table Schema**
```sql
CREATE TABLE tasks (
    id SERIAL PRIMARY KEY,              -- Auto-incrementing unique ID
    title VARCHAR(255) NOT NULL,        -- Required task title
    description TEXT,                   -- Optional long description
    completed BOOLEAN DEFAULT FALSE,    -- Task status
    created_at TIMESTAMP DEFAULT NOW    -- When task was created
);
```

**Database Concepts:**
- **SERIAL:** Auto-incrementing integer (1, 2, 3, ...)
- **PRIMARY KEY:** Unique identifier for each row
- **VARCHAR(255):** Variable-length string, max 255 characters
- **TEXT:** Unlimited text length
- **BOOLEAN:** True/false values
- **TIMESTAMP:** Date and time
- **DEFAULT:** Automatic value if none provided
- **NOT NULL:** Field is required

### **Phase 3: Go Application Architecture**

#### **Step 6: Project Structure**
```
task-api/
├── go.mod          # Dependencies
├── main.go         # Main application
├── setup.sql       # Database schema
└── README.md       # Documentation
```

#### **Step 7: Import Dependencies**
```go
import (
    "database/sql"           // Go's database interface
    "encoding/json"          // JSON serialization
    "net/http"              // HTTP server
    "github.com/gorilla/mux" // URL routing
    _ "github.com/lib/pq"    // PostgreSQL driver
)
```

**Why Each Import:**
- `database/sql`: Standard Go database interface
- `encoding/json`: Convert Go structs ↔ JSON
- `net/http`: Handle HTTP requests/responses
- `gorilla/mux`: Advanced URL routing and middleware
- `lib/pq`: PostgreSQL driver (blank import registers driver)

### **Phase 4: Data Modeling**

#### **Step 8: Define Data Structure**
```go
type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}
```

**Go Concepts:**
- **struct:** Custom data type (like class in other languages)
- **json tags:** Tell Go how to serialize to JSON
- **Exported fields:** Capital letters = public access
- **time.Time:** Go's built-in date/time type

### **Phase 5: Database Connection**

#### **Step 9: Connection Management**
```go
var db *sql.DB  // Global database connection pool

func initDB() {
    // Connection string
    psqlInfo := fmt.Sprintf("host=%s port=%d user=%s dbname=%s sslmode=disable",
        "localhost", 5432, "umer", "taskapi")
    
    // Open connection
    db, err = sql.Open("postgres", psqlInfo)
    if err != nil {
        log.Fatal("Failed to connect:", err)
    }
    
    // Test connection
    if err = db.Ping(); err != nil {
        log.Fatal("Failed to ping:", err)
    }
}
```

**Database Concepts:**
- **Connection Pool:** Reuses database connections efficiently
- **Connection String:** URL-like format for database credentials
- **Ping:** Tests if database is reachable
- **Global Variable:** One connection pool for entire application

### **Phase 6: HTTP Handlers (Business Logic)**

#### **Step 10: GET All Tasks**
```go
func getTasks(w http.ResponseWriter, r *http.Request) {
    // Execute SQL query
    rows, err := db.Query("SELECT id, title, description, completed, created_at FROM tasks ORDER BY id")
    if err != nil {
        http.Error(w, "Database error", 500)
        return
    }
    defer rows.Close()  // Important: always close rows
    
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
    
    // Return JSON
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(tasks)
}
```

**Go Concepts:**
- **http.ResponseWriter:** Where we write the response
- **http.Request:** Contains request data
- **defer:** Ensures rows.Close() runs when function exits
- **Scan:** Copies database row into Go struct
- **JSON Encoding:** Converts Go struct to JSON

#### **Step 11: POST Create Task**
```go
func createTask(w http.ResponseWriter, r *http.Request) {
    var newTask Task
    
    // Parse JSON from request body
    if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
        http.Error(w, "Invalid JSON", 400)
        return
    }
    
    // Validate data
    if newTask.Title == "" {
        http.Error(w, "Title required", 400)
        return
    }
    
    // Insert into database
    err := db.QueryRow(
        "INSERT INTO tasks (title, description, completed) VALUES ($1, $2, $3) RETURNING id, created_at",
        newTask.Title, newTask.Description, false,
    ).Scan(&newTask.ID, &newTask.CreatedAt)
    
    if err != nil {
        http.Error(w, "Database error", 500)
        return
    }
    
    // Return created task
    w.WriteHeader(http.StatusCreated)  // 201 status
    json.NewEncoder(w).Encode(newTask)
}
```

**SQL Concepts:**
- **INSERT:** Adds new row to table
- **RETURNING:** Gets back generated values (like auto-increment ID)
- **$1, $2, $3:** Parameter placeholders (prevents SQL injection)
- **QueryRow:** For queries that return exactly one row

### **Phase 7: Routing & Middleware**

#### **Step 12: URL Routing**
```go
func main() {
    initDB()
    defer db.Close()
    
    // Create router
    r := mux.NewRouter()
    
    // Add middleware
    r.Use(loggingMiddleware)
    
    // Define routes
    r.HandleFunc("/health", healthCheck).Methods("GET")
    r.HandleFunc("/tasks", getTasks).Methods("GET")
    r.HandleFunc("/tasks", createTask).Methods("POST")
    r.HandleFunc("/tasks/{id}", getTask).Methods("GET")
    r.HandleFunc("/tasks/{id}", updateTask).Methods("PUT")
    r.HandleFunc("/tasks/{id}", deleteTask).Methods("DELETE")
    
    // Start server
    log.Fatal(http.ListenAndServe(":8080", r))
}
```

**HTTP Concepts:**
- **REST:** Representational State Transfer
  - GET: Retrieve data
  - POST: Create new data
  - PUT: Update existing data
  - DELETE: Remove data
- **{id}:** URL parameter (like /tasks/123)
- **Middleware:** Code that runs before every request

#### **Step 13: Middleware Example**
```go
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        log.Printf("Started %s %s", r.Method, r.URL.Path)
        
        next.ServeHTTP(w, r)  // Call next handler
        
        log.Printf("Completed %s %s in %v", r.Method, r.URL.Path, time.Since(start))
    })
}
```

**Middleware Concepts:**
- **Cross-cutting concerns:** Logging, authentication, CORS
- **Chain of responsibility:** Each middleware calls the next
- **Timing:** Measure request duration

---

## 🗄️ **Database Exploration (Your "Supabase Dashboard")**

### **Where is Your Database?**
**Physical Location:** `/opt/homebrew/var/postgresql@15/`
**Process:** Running as background service on port 5432
**Access:** Local only (not on internet like Supabase)

### **Command Line Interface (Like SQL Editor)**
```bash
# Connect to database
psql -d taskapi

# View all tables
\d

# View specific table structure
\d tasks

# Run queries
SELECT * FROM tasks;
SELECT COUNT(*) FROM tasks;
SELECT * FROM tasks WHERE completed = true;
```

### **Visual Interface (Like Supabase Dashboard)**
**pgAdmin 4** (now installed):
1. Open "pgAdmin 4" from Applications
2. Create new server connection:
   - Host: localhost
   - Port: 5432
   - Username: umer
   - Database: taskapi
3. Browse tables, run queries, view data

### **Common Database Operations**
```sql
-- View all data
SELECT * FROM tasks ORDER BY created_at DESC;

-- Filter data
SELECT * FROM tasks WHERE completed = false;

-- Count records
SELECT COUNT(*) FROM tasks;

-- Update data
UPDATE tasks SET completed = true WHERE id = 1;

-- Delete data
DELETE FROM tasks WHERE id = 1;

-- Add new column
ALTER TABLE tasks ADD COLUMN priority VARCHAR(20) DEFAULT 'medium';
```

---

## 🔐 **Security & Best Practices**

### **SQL Injection Prevention**
```go
// SECURE ✅
db.Query("SELECT * FROM tasks WHERE id = $1", userID)

// VULNERABLE ❌
db.Query("SELECT * FROM tasks WHERE id = " + userID)
```

### **Error Handling**
```go
if err != nil {
    log.Printf("Database error: %v", err)  // Log for debugging
    http.Error(w, "Internal server error", 500)  // Generic user message
    return
}
```

### **Resource Management**
```go
rows, err := db.Query("SELECT ...")
if err != nil {
    return err
}
defer rows.Close()  // Always close to prevent memory leaks
```

---

## 📊 **Performance & Scaling**

### **Connection Pooling**
```go
db.SetMaxOpenConns(25)        // Maximum concurrent connections
db.SetMaxIdleConns(25)        // Keep connections alive
db.SetConnMaxLifetime(5 * time.Minute)  // Connection lifespan
```

### **Database Indexing**
```sql
-- Speed up queries on commonly searched fields
CREATE INDEX idx_tasks_completed ON tasks(completed);
CREATE INDEX idx_tasks_created_at ON tasks(created_at);
```

### **Query Optimization**
```sql
-- Good: Specific fields
SELECT id, title FROM tasks WHERE completed = false;

-- Avoid: Select everything
SELECT * FROM tasks;
```

---

## 🚀 **What You've Learned**

### **Go Backend Concepts**
✅ HTTP servers and routing  
✅ JSON serialization/deserialization  
✅ Error handling patterns  
✅ Middleware and request processing  
✅ Database integration  
✅ REST API design  

### **Database Concepts**
✅ Relational database design  
✅ SQL queries (CRUD operations)  
✅ Primary keys and auto-increment  
✅ Data types and constraints  
✅ Connection pooling  
✅ SQL injection prevention  

### **Production Concepts**
✅ Data persistence  
✅ Concurrent request handling  
✅ Resource management  
✅ Error logging  
✅ API versioning  

---

## 🎯 **Next Steps**

1. **Authentication:** JWT tokens, user sessions
2. **Validation:** Input sanitization, data validation
3. **Testing:** Unit tests, integration tests
4. **Deployment:** Docker, cloud deployment
5. **Monitoring:** Logging, metrics, health checks
6. **Advanced SQL:** Joins, indexes, transactions

You now have a **production-ready foundation** for backend development! 🎉 