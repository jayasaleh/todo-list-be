package main

import (
	"fmt"
	"log"

	"github.com/jayasaleh/todo-list/be/internal/config"
	"github.com/jayasaleh/todo-list/be/internal/database"
	"github.com/jayasaleh/todo-list/be/internal/router"
)

func main() {
	cfg := config.LoadConfig()

	if _, err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := router.SetupRouter()

	port := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("API available at http://localhost%s/api", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
