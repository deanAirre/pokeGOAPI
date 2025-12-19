package service

import (
	"database/sql"
	"fmt"
	"log"
	"pokeAPI/dto"
	"pokeAPI/model"
)

// PokemonService handles Pokemon business logic
type PokemonService struct {
	db            *sql.DB
	pokeAPIClient *PokeAPIClient
}

// NewPokemonService creates a new Pokemon service
func NewPokemonService(db *sql.DB) *PokemonService {
	return &PokemonService{
		db:            db,
		pokeAPIClient: NewPokeAPIClient(),
	}
}

// SavePokemon saves a Pokemon and its related data to the database
func (s *PokemonService) SavePokemon(apiPokemon *dto.PokeAPIResponse) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Get sprite URL (prefer official artwork)
	spriteURL := apiPokemon.Sprites.Other.OfficialArtwork.FrontDefault
	if spriteURL == "" {
		spriteURL = apiPokemon.Sprites.FrontDefault
	}

	// Insert or update Pokemon
	var pokemonID int
	err = tx.QueryRow(`
		INSERT INTO pokemon (pokedex_id, name, height, weight, sprite_url)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (pokedex_id) 
		DO UPDATE SET name = $2, height = $3, weight = $4, sprite_url = $5
		RETURNING id
	`, apiPokemon.ID, apiPokemon.Name, apiPokemon.Height, apiPokemon.Weight, spriteURL).Scan(&pokemonID)
	
	if err != nil {
		return fmt.Errorf("failed to save pokemon: %w", err)
	}

	// Delete existing types, abilities, and stats (will re-insert)
	if _, err := tx.Exec("DELETE FROM pokemon_types WHERE pokemon_id = $1", pokemonID); err != nil {
		return fmt.Errorf("failed to delete old types: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM pokemon_abilities WHERE pokemon_id = $1", pokemonID); err != nil {
		return fmt.Errorf("failed to delete old abilities: %w", err)
	}
	if _, err := tx.Exec("DELETE FROM pokemon_stats WHERE pokemon_id = $1", pokemonID); err != nil {
		return fmt.Errorf("failed to delete old stats: %w", err)
	}

	// Insert types
	for _, typeSlot := range apiPokemon.Types {
		_, err := tx.Exec(`
			INSERT INTO pokemon_types (pokemon_id, type_name, slot)
			VALUES ($1, $2, $3)
		`, pokemonID, typeSlot.Type.Name, typeSlot.Slot)
		
		if err != nil {
			return fmt.Errorf("failed to save type: %w", err)
		}
	}

	// Insert abilities
	for _, abilitySlot := range apiPokemon.Abilities {
		_, err := tx.Exec(`
			INSERT INTO pokemon_abilities (pokemon_id, ability_name, is_hidden, slot)
			VALUES ($1, $2, $3, $4)
		`, pokemonID, abilitySlot.Ability.Name, abilitySlot.IsHidden, abilitySlot.Slot)
		
		if err != nil {
			return fmt.Errorf("failed to save ability: %w", err)
		}
	}

	// Insert stats
	stats := make(map[string]int)
	for _, stat := range apiPokemon.Stats {
		stats[stat.Stat.Name] = stat.BaseStat
	}

	_, err = tx.Exec(`
		INSERT INTO pokemon_stats 
		(pokemon_id, hp, attack, defense, special_attack, special_defense, speed)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, pokemonID, stats["hp"], stats["attack"], stats["defense"], 
	   stats["special-attack"], stats["special-defense"], stats["speed"])
	
	if err != nil {
		return fmt.Errorf("failed to save stats: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// SyncGen5Pokemon fetches and saves all Gen 5 Pokemon
func (s *PokemonService) SyncGen5Pokemon() error {
	log.Println("Starting Gen 5 Pokemon sync...")
	
	start, end := GetGen5Range()
	total := end - start + 1
	
	for id := start; id <= end; id++ {
		log.Printf("Fetching Pokemon %d/%d (ID: %d)...", id-start+1, total, id)
		
		pokemon, err := s.pokeAPIClient.FetchPokemon(id)
		if err != nil {
			return fmt.Errorf("failed to fetch pokemon %d: %w", id, err)
		}
		
		if err := s.SavePokemon(pokemon); err != nil {
			return fmt.Errorf("failed to save pokemon %d: %w", id, err)
		}
		
		log.Printf("✓ Saved %s (#%d)", pokemon.Name, pokemon.ID)
	}
	
	log.Printf("✓ Gen 5 sync complete! Saved %d Pokemon", total)
	return nil
}

// GetAllPokemon retrieves all Pokemon from the database
func (s *PokemonService) GetAllPokemon() ([]model.Pokemon, error) {
	rows, err := s.db.Query(`
		SELECT id, pokedex_id, name, height, weight, sprite_url, created_at
		FROM pokemon
		ORDER BY pokedex_id
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query pokemon: %w", err)
	}
	defer rows.Close()

	var pokemons []model.Pokemon
	for rows.Next() {
		var p model.Pokemon
		err := rows.Scan(&p.ID, &p.PokedexID, &p.Name, &p.Height, &p.Weight, &p.SpriteURL, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pokemon: %w", err)
		}
		pokemons = append(pokemons, p)
	}

	return pokemons, nil
}

// GetPokemonByID retrieves a single Pokemon by its Pokedex ID
func (s *PokemonService) GetPokemonByID(pokedexID int) (*model.Pokemon, error) {
	var p model.Pokemon
	err := s.db.QueryRow(`
		SELECT id, pokedex_id, name, height, weight, sprite_url, created_at
		FROM pokemon
		WHERE pokedex_id = $1
	`, pokedexID).Scan(&p.ID, &p.PokedexID, &p.Name, &p.Height, &p.Weight, &p.SpriteURL, &p.CreatedAt)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pokemon with pokedex id %d not found", pokedexID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query pokemon: %w", err)
	}

	return &p, nil
}