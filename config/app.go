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

	// Env loading from system environment variables
	config := &Config {
		DBHost: getEnv("DB_HOST", "localhost"),
		DBPort: getEnv("DB_PORT", "5432"),
		DBUser: getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName: getEnv("DB_NAME", "pokemon_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	if config.DBPassword == "postgres" {
		fmt.Println("This is default password for test, change it later");
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