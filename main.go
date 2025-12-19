package main

import (
	"log"
	"net/http"
	"pokeAPI/config"
	"pokeAPI/controller"
	"pokeAPI/service"
)

func main() {
	// 1. Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println(" Config loaded")

	// 2. Connect to database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println(" Database connected")

	// 3. Run migrations
	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println(" Migrations completed")

	// 4. Initialize services
	pokemonService := service.NewPokemonService(db)

	// 5. Initialize controllers
	pokemonController := controller.NewPokemonController(pokemonService)

	// 6. Setup routes
	http.HandleFunc("/health", controller.HealthCheck)
	http.HandleFunc("/api/pokemon", pokemonController.GetAllPokemon)
	http.HandleFunc("/api/pokemon/", pokemonController.GetPokemonByID)
	http.HandleFunc("/api/pokemon/sync", pokemonController.SyncGen5Pokemon)

	// 7. Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf(" Server starting on http://localhost%s", serverAddr)
	log.Println(" Available endpoints:")
	log.Println("   GET  /health              - Health check")
	log.Println("   GET  /api/pokemon         - List all Pokemon")
	log.Println("   GET  /api/pokemon/{id}    - Get Pokemon by Pokedex ID")
	log.Println("   POST /api/pokemon/sync    - Sync Gen 5 Pokemon from PokeAPI")
	
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}