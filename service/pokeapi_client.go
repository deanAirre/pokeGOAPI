package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"pokeAPI/dto"
	"time"
)

const (
	pokeAPIBaseURL = "https://pokeapi.co/api/v2"
	gen5StartID    = 494 // Victini
	gen5EndID      = 649 // Genesect
)

// PokeAPIClient handles requests to PokeAPI
type PokeAPIClient struct {
	httpClient *http.Client
	baseURL    string
}

// NewPokeAPIClient creates a new PokeAPI client
func NewPokeAPIClient() *PokeAPIClient {
	return &PokeAPIClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: pokeAPIBaseURL,
	}
}

// FetchPokemon fetches a single Pokemon by ID from PokeAPI
func (c *PokeAPIClient) FetchPokemon(id int) (*dto.PokeAPIResponse, error) {
	url := fmt.Sprintf("%s/pokemon/%d", c.baseURL, id)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pokemon %d: %w", id, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("pokeapi returned status %d for pokemon %d", resp.StatusCode, id)
	}

	var pokemon dto.PokeAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&pokemon); err != nil {
		return nil, fmt.Errorf("failed to decode pokemon %d: %w", id, err)
	}

	return &pokemon, nil
}

// FetchGen5Pokemon fetches all Gen 5 Pokemon (IDs 494-649)
func (c *PokeAPIClient) FetchGen5Pokemon() ([]*dto.PokeAPIResponse, error) {
	var pokemons []*dto.PokeAPIResponse
	
	for id := gen5StartID; id <= gen5EndID; id++ {
		pokemon, err := c.FetchPokemon(id)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch gen 5 pokemon: %w", err)
		}
		
		pokemons = append(pokemons, pokemon)
		
		// Be nice to PokeAPI - small delay between requests
		time.Sleep(100 * time.Millisecond)
	}
	
	return pokemons, nil
}

// GetGen5Range returns the ID range for Gen 5 Pokemon
func GetGen5Range() (start, end int) {
	return gen5StartID, gen5EndID
}