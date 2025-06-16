# Task Management API

A clean architecture Go backend API for task management with authentication, built with Gin, GORM, and PostgreSQL.

## Features

- **Clean Architecture**: Organized into domain, repository, service, and handler layers
- **JWT Authentication**: Secure user registration and login
- **CRUD Operations**: Complete task management functionality
- **PostgreSQL Database**: Using GORM for database operations
- **Vercel Deployment**: Ready for serverless deployment
- **Neon Database**: Configured for Neon PostgreSQL

## Project Structure

```
├── cmd/api/           # Application entry point
├── internal/
│   ├── domain/        # Business entities and DTOs
│   ├── repository/    # Data access layer
│   ├── service/       # Business logic layer
│   └── handler/       # HTTP handlers (controllers)
├── pkg/
│   ├── config/        # Configuration management
│   ├── database/      # Database connection
│   └── middleware/    # HTTP middleware
├── api/               # Vercel serverless entry point
└── docs/              # Documentation
```

## Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL (or Neon database)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd dummy-backend
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables:
```bash
cp env.example .env
```

Edit `.env` with your database credentials:
```env
DB_DSN=postgresql://your_user:your_password@your_host/your_database?sslmode=require
JWT_SECRET=your-super-secret-jwt-key-change-in-production
PORT=8080
GIN_MODE=debug
```

4. Run the application:
```bash
go run cmd/api/main.go
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - Login user

### Tasks (Requires Authentication)

- `GET /api/tasks` - Get all tasks
- `POST /api/tasks` - Create a new task
- `GET /api/tasks/:id` - Get task by ID
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task

### Health Check

- `GET /health` - Health check endpoint

## Authentication

Include the JWT token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

## Example Requests

### Register User
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "user@example.com", "password": "password123"}'
```

### Create Task
```bash
curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your-token>" \
  -d '{"title": "Learn Go", "description": "Study Go programming language"}'
```

### Get All Tasks
```bash
curl -X GET http://localhost:8080/api/tasks \
  -H "Authorization: Bearer <your-token>"
```

## Database Schema

### Tasks Table
```sql
CREATE TABLE IF NOT EXISTS tasks (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    completed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

### Users Table
```sql
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

## Deployment to Vercel

1. Install Vercel CLI:
```bash
npm i -g vercel
```

2. Login to Vercel:
```bash
vercel login
```

3. Deploy:
```bash
vercel --prod
```

4. Set environment variables in Vercel dashboard:
   - `DB_DSN`: Your Neon database connection string
   - `JWT_SECRET`: Your JWT secret key
   - `GIN_MODE`: `release`

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DB_DSN` | Database connection string | Required |
| `JWT_SECRET` | JWT signing secret | `default-secret-key` |
| `PORT` | Server port | `8080` |
| `GIN_MODE` | Gin mode (debug/release) | `debug` |

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License. 