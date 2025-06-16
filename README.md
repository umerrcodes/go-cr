# Go Task Management API with PostgreSQL 🚀

This is your first Go backend project with a real database! It's a REST API for managing tasks that demonstrates core backend development concepts using PostgreSQL for data persistence.

## What You'll Learn

- **HTTP Server**: How to create and run an HTTP server in Go
- **REST API**: Building RESTful endpoints (GET, POST, PUT, DELETE)
- **JSON Handling**: Parsing and returning JSON data
- **Routing**: Using Gorilla Mux for URL routing
- **Middleware**: Adding logging middleware to requests
- **Error Handling**: Proper HTTP status codes and error responses
- **Data Structures**: Working with structs and slices in Go
- **Database Integration**: Connecting to PostgreSQL and executing SQL queries
- **SQL Operations**: CREATE, INSERT, SELECT, UPDATE, DELETE operations
- **Connection Management**: Database connection pooling and error handling

## Project Structure

```
.
├── go.mod          # Go module file (dependencies)
├── main.go         # Main application code
├── setup.sql       # Database schema and sample data
└── README.md       # This file
```

## How to Run

1. **Install dependencies**:
   ```bash
   go mod tidy
   ```

2. **Run the server**:
   ```bash
   go run main.go
   ```

3. **The server will start on port 8080** and show available endpoints.

## API Endpoints

### Health Check
```bash
GET /health
```
Returns server status and version info.

### Get All Tasks
```bash
GET /tasks
```
Returns a list of all tasks.

### Get Specific Task
```bash
GET /tasks/{id}
```
Returns a specific task by ID.

### Create New Task
```bash
POST /tasks
Content-Type: application/json

{
    "title": "Learn Go",
    "description": "Study Go programming language"
}
```

### Update Task
```bash
PUT /tasks/{id}
Content-Type: application/json

{
    "title": "Updated title",
    "description": "Updated description",
    "completed": true
}
```

### Delete Task
```bash
DELETE /tasks/{id}
```

## Testing the API

You can test the API using curl commands:

### 1. Check server health
```bash
curl http://localhost:8080/health
```

### 2. Get all tasks
```bash
curl http://localhost:8080/tasks
```

### 3. Create a new task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -d '{"title": "My New Task", "description": "This is a test task"}'
```

### 4. Get a specific task
```bash
curl http://localhost:8080/tasks/1
```

### 5. Update a task
```bash
curl -X PUT http://localhost:8080/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "Updated Task", "description": "Updated description", "completed": true}'
```

### 6. Delete a task
```bash
curl -X DELETE http://localhost:8080/tasks/1
```

## Key Concepts Explained

### 1. **Structs** (`Task`)
Go uses structs to define data structures. Our `Task` struct represents a task with fields like ID, Title, Description, etc.

### 2. **JSON Tags**
The `json:"id"` tags tell Go how to serialize/deserialize the struct to/from JSON.

### 3. **HTTP Handlers**
Functions like `getTasks`, `createTask` handle different HTTP requests. They take `http.ResponseWriter` and `*http.Request` parameters.

### 4. **Routing with Gorilla Mux**
We use the Gorilla Mux router to handle different URL patterns and HTTP methods.

### 5. **Middleware**
The `loggingMiddleware` function logs every request, showing how to add cross-cutting concerns.

### 6. **Error Handling**
We return appropriate HTTP status codes (200, 201, 400, 404, etc.) and error messages.

## Next Steps

Once you understand this project, you can extend it by:
1. Adding a real database (PostgreSQL, MongoDB)
2. Adding authentication/authorization
3. Adding input validation
4. Adding unit tests
5. Adding configuration management
6. Dockerizing the application

## Common Go Patterns You'll See

- **Error handling**: `if err != nil { ... }`
- **Slices**: `[]Task` for storing multiple tasks
- **Pointers**: `*http.Request` for passing references
- **Interfaces**: `http.Handler` interface for middleware
- **Goroutines**: (not used here, but common for concurrency)

Start the server and try the API endpoints to see how everything works together! 