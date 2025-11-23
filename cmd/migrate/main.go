package main

import (
	"flag"
	"log"

	"github.com/jayasaleh/todo-list/be/internal/config"
	"github.com/jayasaleh/todo-list/be/internal/database"
)

func main() {
	action := flag.String("action", "up", "Migration action: up or down")
	flag.Parse()

	cfg := config.LoadConfig()

	if _, err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	switch *action {
	case "up":
		if err := database.AutoMigrate(); err != nil {
			log.Fatalf("Failed to run database migration: %v", err)
		}
		log.Println("Migration up completed successfully")
	case "down":
		log.Println("Down migration not implemented. Please use SQL files manually.")
	default:
		log.Fatalf("Unknown action: %s. Use 'up' or 'down'", *action)
	}
}
