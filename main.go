package main

import (
	"log"
	"net/http"
	"pokeAPI/config"
	"pokeAPI/controller"
	"pokeAPI/service"
)

// CORS Middleware
func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS"{
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}



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
	defer db.Close() // do this after function returns
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
	http.HandleFunc("/health", enableCORS(controller.HealthCheck))
	http.HandleFunc("/api/pokemon", enableCORS(pokemonController.GetAllPokemon))
	http.HandleFunc("/api/pokemon/", enableCORS(pokemonController.GetPokemonByID))
	http.HandleFunc("/api/pokemon/sync", enableCORS(pokemonController.SyncGen5Pokemon))
	http.HandleFunc("/api/pokemon/sync/status", enableCORS(pokemonController.GetSyncStatus))


	// 7. Start server
	serverAddr := ":" + cfg.ServerPort
	log.Printf(" Server starting on http://localhost%s", serverAddr)
	log.Println(" Available endpoints:")
	log.Println("   GET  /health              		- Health check")
	log.Println("   GET  /api/pokemon         		- List all Pokemon")
	log.Println("   GET  /api/pokemon/{id}    		- Get Pokemon by Pokedex ID")
	log.Println("   POST /api/pokemon/sync    		- Sync Gen 5 Pokemon from PokeAPI")
	log.Println("	GET /api/pokemon/sync/status	- Get last sync information")
	
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}