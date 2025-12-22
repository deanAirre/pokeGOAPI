package model

import "time"

// Pokemon represents a Pokemon entity in the database
type Pokemon struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Height     int       `json:"height"`      // in decimeters
	Weight     int       `json:"weight"`      // in hectograms
	SpriteURL  string    `json:"sprite_url"`
	AnimatedFront string `json:"animated_front"`
	AnimatedBack string `json:"animated_back"`
	CreatedAt  time.Time `json:"created_at"`
}

// PokemonType represents a Pokemon's type (fire, water, etc.)
type PokemonType struct {
	ID        int    `json:"id"`
	PokemonID int    `json:"pokemon_id"`
	TypeName  string `json:"type_name"`
	Slot      int    `json:"slot"` // Primary (1) or Secondary (2)
}

// PokemonAbility represents a Pokemon's ability
type PokemonAbility struct {
	ID          int    `json:"id"`
	PokemonID   int    `json:"pokemon_id"`
	AbilityName string `json:"ability_name"`
	IsHidden    bool   `json:"is_hidden"`
	Slot        int    `json:"slot"`
}

// PokemonStats represents a Pokemon's base stats
type PokemonStats struct {
	ID             int `json:"id"`
	PokemonID      int `json:"pokemon_id"`
	HP             int `json:"hp"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	Speed          int `json:"speed"`
}