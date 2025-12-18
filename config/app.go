package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBHost string;
	DBPort string;
	DBUser string;
	DBPassword string;
	DBName string;
	ServerPort string;
}

// LoadConfig read from environment variables
func LoadConfig() (*Config, error) {

	// Check required field existence for .env before autoloads
	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		return nil, fmt.Errorf("DB_PASSWORD environment variable is required")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		return nil, fmt.Errorf("DB_NAME environment variable is required")
	}

	// Env loading from system environment variables
	config := &Config {
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName: getEnv("DB_NAME", "pokemon_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	// return config and nil error if success
	return config, nil;
}

// Helper function to get env variable with default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != ""{ 
		return value
	}
	return defaultValue
}