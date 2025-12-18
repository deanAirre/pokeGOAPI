package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver like JDBC
)

// ConnectDatabase establishes connection to PostgreSQL
func ConnectDatabase(cfg *Config) (*sql.DB, error){
	// Build connection string
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName,
	)

	// Open connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func RunMigrations(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS pokemon (
		id SERIAL PRIMARY KEY,
		pokedex_id INT UNIQUE NOT NULL,
		name VARCHAR(100) NOT NULL,
		height INT,
		weight INT,
		sprite_url TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP	
	)`

	if _, err := db.Exec(query); err != nil {
		return fmt.Errorf("failed to create pokemon table: %w", err)
	}

	return nil;
}