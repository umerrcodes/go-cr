package config

import (
	"os"
)

type Config struct {
	DatabaseDSN string
	JWTSecret   string
	Port        string
	GinMode     string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseDSN: getEnv("DB_DSN", ""),
		JWTSecret:   getEnv("JWT_SECRET", "default-secret-key"),
		Port:        getEnv("PORT", "8080"),
		GinMode:     getEnv("GIN_MODE", "debug"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
