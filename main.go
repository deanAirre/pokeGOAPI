package main

import (
	"log"
	"net/http"
	"pokeAPI/config"
)

func main(){
	// Loading configurations
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	log.Println(" Config loaded")

	// Connect to database
	db, err := config.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to databse: %v", err)
	}
	defer db.Close()
	log.Println("Database connected")

	// Run migrations
	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed");

	http.HandleFunc("/health", healthCheckHandler)

	serverAddr := ":" + cfg.ServerPort
	log.Printf(" Server starting on http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// Health check endpoint
func healthCheckHandler(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy", "service":"pokemon-api"}`))
	w.Write([]byte("\n"))
}