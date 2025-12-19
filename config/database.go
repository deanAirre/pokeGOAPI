package config

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // PostgreSQL driver like JDBC
)

// ConnectDatabase establishes connection to PostgreSQL
func ConnectDatabase(cfg *Config) (*sql.DB, error) {
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

// RunMigrations creates necessary database tables
func RunMigrations(db *sql.DB) error {
	migrations := []string{
		// Pokemon table
		`CREATE TABLE IF NOT EXISTS pokemon (
			id SERIAL PRIMARY KEY,
			pokedex_id INT UNIQUE NOT NULL,
			name VARCHAR(100) NOT NULL,
			height INT,
			weight INT,
			sprite_url TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)`,
		
		// Pokemon types table
		`CREATE TABLE IF NOT EXISTS pokemon_types (
			id SERIAL PRIMARY KEY,
			pokemon_id INT NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
			type_name VARCHAR(50) NOT NULL,
			slot INT NOT NULL,
			UNIQUE(pokemon_id, slot)
		)`,
		
		// Pokemon abilities table
		`CREATE TABLE IF NOT EXISTS pokemon_abilities (
			id SERIAL PRIMARY KEY,
			pokemon_id INT NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
			ability_name VARCHAR(100) NOT NULL,
			is_hidden BOOLEAN DEFAULT FALSE,
			slot INT NOT NULL,
			UNIQUE(pokemon_id, slot)
		)`,
		
		// Pokemon stats table
		`CREATE TABLE IF NOT EXISTS pokemon_stats (
			id SERIAL PRIMARY KEY,
			pokemon_id INT NOT NULL REFERENCES pokemon(id) ON DELETE CASCADE,
			hp INT NOT NULL,
			attack INT NOT NULL,
			defense INT NOT NULL,
			special_attack INT NOT NULL,
			special_defense INT NOT NULL,
			speed INT NOT NULL,
			UNIQUE(pokemon_id)
		)`,
		
		// Indexes for better performance
		`CREATE INDEX IF NOT EXISTS idx_pokemon_pokedex_id ON pokemon(pokedex_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pokemon_name ON pokemon(name)`,
		`CREATE INDEX IF NOT EXISTS idx_pokemon_types_pokemon_id ON pokemon_types(pokemon_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pokemon_abilities_pokemon_id ON pokemon_abilities(pokemon_id)`,
		`CREATE INDEX IF NOT EXISTS idx_pokemon_stats_pokemon_id ON pokemon_stats(pokemon_id)`,
	}

	// Execute each migration
	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}