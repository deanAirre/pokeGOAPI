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

	var animatedFront, animatedBack string
	if apiPokemon.Sprites.Versions != nil {
		if genV, ok := apiPokemon.Sprites.Versions["generation-v"]; ok {
			if blackWhite, ok := genV["black-white"]; ok {
				if blackWhite.Animated != nil {
					animatedFront = blackWhite.Animated.FrontDefault
					animatedBack = blackWhite.Animated.BackDefault
				}
			}
		}
	}

	// Insert or update Pokemon
	var pokemonID int
	err = tx.QueryRow(`
    INSERT INTO pokemon (pokedex_id, name, height, weight, sprite_url, animated_front, animated_back)
    VALUES ($1, $2, $3, $4, $5, $6, $7)
    ON CONFLICT (pokedex_id) 
    DO UPDATE SET name = $2, height = $3, weight = $4, sprite_url = $5, 
                  animated_front = $6, animated_back = $7
    RETURNING id
	`, apiPokemon.ID, apiPokemon.Name, apiPokemon.Height, apiPokemon.Weight, 
   spriteURL, animatedFront, animatedBack).Scan(&pokemonID)
	
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
	successCount := 0
	
	for id := start; id <= end; id++ {
		log.Printf("Fetching Pokemon %d/%d (ID: %d)...", id-start+1, total, id)
		
		pokemon, err := s.pokeAPIClient.FetchPokemon(id)
		if err != nil {
			log.Printf("Warning: Failed to fetch pokemon %d: %v", id, err)
			continue
		}
		
		if err := s.SavePokemon(pokemon); err != nil {
			log.Printf("Warning: Failed to save pokemon %d: %v", id, err)
			continue
		}

		successCount++
		log.Printf(" Saved %s (#%d)", pokemon.Name, pokemon.ID)
	}
	// Update sync metadata
	if err := s.updateSyncMetaData("gen5", successCount); err != nil {
		log.Printf("Warning: Failed to update sync metadata: %v", err)
	}

	log.Printf(" Gen 5 sync complete! Saved %d/%d Pokemon", successCount, total)
	return nil
}


// GetPokemonPaginated retrieves Pokemon with pagination, filtering, and sorting
func (s *PokemonService) GetPokemonPaginated(limit, offset int, sortBy, order, typeFilter string) (map[string]interface{}, error) {
	// Validate and sanitize inputs
	if limit <= 0 || limit > 100 {
		limit = 20 // Default
	}
	if offset < 0 {
		offset = 0
	}
	
	// Whitelist allowed sort columns (prevent SQL injection)
	allowedSortColumns := map[string]bool{
		"pokedex_id": true,
		"name":       true,
		"height":     true,
		"weight":     true,
		"created_at": true,
	}
	if !allowedSortColumns[sortBy] {
		sortBy = "pokedex_id" // Default
	}
	
	// Validate order
	if order != "asc" && order != "desc" {
		order = "asc" // Default
	}
	
	// Build query with optional type filter
	var query string
	var countQuery string
	var args []interface{}
	var countArgs []interface{}
	
	if typeFilter != "" {
		// Filter by type
		query = `
			SELECT p.id, p.pokedex_id, p.name, p.height, p.weight, p.sprite_url, 
           	p.animated_front, p.animated_back, p.created_at
    		FROM pokemon p
			INNER JOIN pokemon_types pt ON p.id = pt.pokemon_id
			WHERE pt.type_name = $1
			ORDER BY p.` + sortBy + ` ` + order + `
			LIMIT $2 OFFSET $3
		`
		countQuery = `
			SELECT COUNT(DISTINCT p.id)
			FROM pokemon p
			INNER JOIN pokemon_types pt ON p.id = pt.pokemon_id
			WHERE pt.type_name = $1
		`
		args = []interface{}{typeFilter, limit, offset}
		countArgs = []interface{}{typeFilter}
	} else {
		// No filter
		query = `
			SELECT p.id, p.pokedex_id, p.name, p.height, p.weight, p.sprite_url, 
           	p.animated_front, p.animated_back, p.created_at
    		FROM pokemon p
			ORDER BY p.` + sortBy + ` ` + order + `
			LIMIT $1 OFFSET $2
		`
		countQuery = `SELECT COUNT(*) FROM pokemon`
		args = []interface{}{limit, offset}
		countArgs = []interface{}{}
	}
	
	// Get total count
	var totalCount int
	err := s.db.QueryRow(countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, fmt.Errorf("failed to get count: %w", err)
	}
	
	// Get paginated data
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query pokemon: %w", err)
	}
	defer rows.Close()
	
	var pokemons []map[string]interface{}
	for rows.Next() {
		var id, pokedexID, height, weight int
		var name, spriteURL, animatedFront, animatedBack, createdAt string
		
		err := rows.Scan(&id, &pokedexID, &name, &height, &weight, &spriteURL, &animatedFront, &animatedBack, &createdAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan pokemon: %w", err)
		}
		
		// Get types for this pokemon
		typesRows, err := s.db.Query(`
			SELECT type_name FROM pokemon_types
			WHERE pokemon_id = $1
			ORDER BY slot
		`, id)
		if err != nil {
			return nil, fmt.Errorf("failed to get types: %w", err)
		}
		
		var types []string
		for typesRows.Next() {
			var typeName string
			if err := typesRows.Scan(&typeName); err != nil {
				typesRows.Close()
				return nil, err
			}
			types = append(types, typeName)
		}
		typesRows.Close()
		
		pokemons = append(pokemons, map[string]interface{}{
    	"id":             pokedexID,
    	"name":           name,
    	"height":         height,
    	"weight":         weight,
    	"sprite_url":     spriteURL,
    	"animated_front": animatedFront,
    	"animated_back":  animatedBack,
    	"types":          types,
    	"created_at":     createdAt,
		})
	}
	
	// Calculate pagination info
	totalPages := (totalCount + limit - 1) / limit
	currentPage := (offset / limit) + 1
	
	return map[string]interface{}{
		"data":         pokemons,
		"total":        totalCount,
		"page":         currentPage,
		"limit":        limit,
		"total_pages":  totalPages,
		"has_next":     offset+limit < totalCount,
		"has_previous": offset > 0,
	}, nil
}

// GetPokemonByID retrieves a single Pokemon by its Pokedex ID
func (s *PokemonService) GetPokemonByID(pokedexID int) (*model.Pokemon, error) {
	var p model.Pokemon
	var dbID int
	err := s.db.QueryRow(`
    SELECT id, pokedex_id, name, height, weight, sprite_url, animated_front, animated_back, created_at
    FROM pokemon
    WHERE pokedex_id = $1
	`, pokedexID).Scan(&dbID, &p.ID, &p.Name, &p.Height, &p.Weight, &p.SpriteURL, 
                   &p.AnimatedFront, &p.AnimatedBack, &p.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("pokemon with pokedex id %d not found", pokedexID)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query pokemon: %w", err)
	}

	return &p, nil
}

// Log Update Sync to Database
func (s *PokemonService) updateSyncMetaData(syncType string, totalSynced int) error{
	_, err := s.db.Exec(`
	
	INSERT INTO sync_metadata (sync_type, last_sync_at, total_synced, status)
		VALUES ($1, NOW(), $2, 'completed')
		ON CONFLICT (sync_type)
		DO UPDATE SET 
			last_sync_at = NOW(),
			total_synced = $2,
			status = 'completed'
	`, syncType, totalSynced)

	return err;
}

func (s *PokemonService) GetLastSyncInfo(syncType string) (map[string]interface{}, error){
	var lastSyncAt string
	var totalSynced int
	var status string

	err := s.db.QueryRow(`
	SELECT last_sync_at, total_synced, status
	FROM sync_metadata
	WHERE sync_type = $1`, syncType).Scan(&lastSyncAt, &totalSynced, &status)

	if err == sql.ErrNoRows{
		return map[string]interface{}{
			"sync_type": syncType,
			"last_sync_at": nil,
			"total_synced": 0,
			"status": 	"never_synced",
		}, nil
	}

	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"sync_type": syncType,
		"last_sync_at": lastSyncAt,
		"total_synced": totalSynced,
		"status": status,
	}, nil
}

// GetPokemonWithDetails retrieves a Pokemon with all related data
func (s *PokemonService) GetPokemonWithDetails(pokedexID int) (map[string]interface{}, error) {
	// Get basic Pokemon info
	pokemon, err := s.GetPokemonByID(pokedexID)
	if err != nil {
		return nil, err
	}

	// Get types
	typesRows, err := s.db.Query(`
		SELECT type_name, slot
		FROM pokemon_types
		WHERE pokemon_id = (SELECT id FROM pokemon WHERE pokedex_id = $1)
		ORDER BY slot
	`, pokedexID)
	if err != nil {
		return nil, fmt.Errorf("failed to get types: %w", err)
	}
	defer typesRows.Close()

	var types []map[string]interface{}
	for typesRows.Next() {
		var typeName string
		var slot int
		if err := typesRows.Scan(&typeName, &slot); err != nil {
			return nil, err
		}
		types = append(types, map[string]interface{}{
			"slot": slot,
			"type": map[string]string{
				"name": typeName,
			},
		})
	}

	// Get abilities
	abilitiesRows, err := s.db.Query(`
		SELECT ability_name, is_hidden, slot
		FROM pokemon_abilities
		WHERE pokemon_id = (SELECT id FROM pokemon WHERE pokedex_id = $1)
		ORDER BY slot
	`, pokedexID)
	if err != nil {
		return nil, fmt.Errorf("failed to get abilities: %w", err)
	}
	defer abilitiesRows.Close()

	var abilities []map[string]interface{}
	for abilitiesRows.Next() {
		var abilityName string
		var isHidden bool
		var slot int
		if err := abilitiesRows.Scan(&abilityName, &isHidden, &slot); err != nil {
			return nil, err
		}
		abilities = append(abilities, map[string]interface{}{
			"is_hidden": isHidden,
			"slot":      slot,
			"ability": map[string]string{
				"name": abilityName,
			},
		})
	}

	// Get stats
	var stats []map[string]interface{}
	var hp, attack, defense, specialAttack, specialDefense, speed int
	err = s.db.QueryRow(`
		SELECT hp, attack, defense, special_attack, special_defense, speed
		FROM pokemon_stats
		WHERE pokemon_id = (SELECT id FROM pokemon WHERE pokedex_id = $1)
	`, pokedexID).Scan(&hp, &attack, &defense, &specialAttack, &specialDefense, &speed)

	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get stats: %w", err)
	}

	if err == nil {
		statNames := []string{"hp", "attack", "defense", "special-attack", "special-defense", "speed"}
		statValues := []int{hp, attack, defense, specialAttack, specialDefense, speed}

		for i, name := range statNames {
			stats = append(stats, map[string]interface{}{
				"base_stat": statValues[i],
				"effort":    0,
				"stat": map[string]string{
					"name": name,
				},
			})
		}
	}

	// Build response matching PokeAPI structure
	response := map[string]interface{}{
		"id":     pokemon.ID,  
		"name":   pokemon.Name,
		"height": pokemon.Height,
		"weight": pokemon.Weight,
		"sprites": map[string]interface{}{
    		"front_default": pokemon.SpriteURL,
    			"other": map[string]interface{}{
        		"official-artwork": map[string]interface{}{
            "front_default": pokemon.SpriteURL,
        	},
    	},
    	"versions": map[string]interface{}{
        	"generation-v": map[string]interface{}{
            	"black-white": map[string]interface{}{
                	"animated": map[string]interface{}{
                    	"front_default": pokemon.AnimatedFront,
                    	"back_default":  pokemon.AnimatedBack,
                		},
            		},
        		},
    		},
		},
		"types":      types,
		"abilities":  abilities,
		"stats":      stats,
		"created_at": pokemon.CreatedAt,
	}

	return response, nil
}