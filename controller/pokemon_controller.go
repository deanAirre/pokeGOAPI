package controller

import (
	"encoding/json"
	"log"
	"net/http"
	"pokeAPI/service"
	"strconv"
	"strings"
)

// PokemonController handles HTTP requests for Pokemon
type PokemonController struct {
	service *service.PokemonService
}

// NewPokemonController creates a new Pokemon controller
func NewPokemonController(service *service.PokemonService) *PokemonController {
	return &PokemonController{
		service: service,
	}
}

// GetAllPokemon handles GET /api/pokemon
func (c *PokemonController) GetAllPokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	pokemons, err := c.service.GetAllPokemon()
	if err != nil {
		log.Printf("Error getting pokemon: %v", err)
		http.Error(w, "Failed to retrieve pokemon", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"count":   len(pokemons),
		"data":    pokemons,
	})
}

// GetPokemonByID handles GET /api/pokemon/{id}
func (c *PokemonController) GetPokemonByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[2])
	if err != nil {
		http.Error(w, "Invalid pokemon ID", http.StatusBadRequest)
		return
	}

	pokemon, err := c.service.GetPokemonByID(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Pokemon not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting pokemon: %v", err)
		http.Error(w, "Failed to retrieve pokemon", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    pokemon,
	})
}

// SyncGen5Pokemon handles POST /api/pokemon/sync
func (c *PokemonController) SyncGen5Pokemon(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Println("Starting Gen 5 Pokemon sync via API...")

	// Run sync in background (this takes time!)
	go func() {
		if err := c.service.SyncGen5Pokemon(); err != nil {
			log.Printf("Sync failed: %v", err)
		}
	}()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Gen 5 Pokemon sync started. This will take a few minutes. Check logs for progress.",
	})
}

// HealthCheck handles GET /health
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "pokemon-api",
	})
}