package dto

// PokeAPIResponse represents the Pokemon data from PokeAPI
type PokeAPIResponse struct {
	ID             int              `json:"id"`
	Name           string           `json:"name"`
	Height         int              `json:"height"`
	Weight         int              `json:"weight"`
	Sprites        Sprites          `json:"sprites"`
	Types          []TypeSlot       `json:"types"`
	Abilities      []AbilitySlot    `json:"abilities"`
	Stats          []StatDetail     `json:"stats"`
}

// Sprites contains Pokemon sprite URLs
type Sprites struct {
	FrontDefault string                     `json:"front_default"`
	Other        Other                      `json:"other"`
	Versions     map[string]GenVersions     `json:"versions"`
}

type GenVersions map[string]BlackWhiteSprites

type BlackWhiteSprites struct {
	Animated *AnimatedSprites `json:"animated,omitempty"`
}

type AnimatedSprites struct {
	FrontDefault string `json:"front_default"`
	BackDefault  string `json:"back_default"`
}

// Other contains additional sprite sources
type Other struct {
	OfficialArtwork OfficialArtwork `json:"official-artwork"`
}

// OfficialArtwork contains high-quality artwork URLs
type OfficialArtwork struct {
	FrontDefault string `json:"front_default"`
}

// TypeSlot represents a Pokemon type with its slot
type TypeSlot struct {
	Slot int  `json:"slot"`
	Type Type `json:"type"`
}

// Type contains type information
type Type struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// AbilitySlot represents a Pokemon ability with its slot
type AbilitySlot struct {
	IsHidden bool    `json:"is_hidden"`
	Slot     int     `json:"slot"`
	Ability  Ability `json:"ability"`
}

// Ability contains ability information
type Ability struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// StatDetail represents a single stat
type StatDetail struct {
	BaseStat int  `json:"base_stat"`
	Effort   int  `json:"effort"`
	Stat     Stat `json:"stat"`
}

// Stat contains stat name
type Stat struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}