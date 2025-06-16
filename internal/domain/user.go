package domain

import "time"

// User represents a user entity
type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"` // "-" excludes from JSON
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// RegisterRequest represents the request payload for user registration
type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the request payload for user login
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the response for authentication
type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}
