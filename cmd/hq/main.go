package main

import (
	"context"
	"log"
	"os"
	"sentinel/internal/hq"
)

func main() {
	// 1. Configuration (Environment variables or flags preferred in prod)
	// Default to local postgres if not set.
	dbConnString := os.Getenv("DATABASE_URL")
	if dbConnString == "" {
		dbConnString = "postgres://postgres:root@localhost:5432/sentinel?sslmode=disable"
		log.Printf("DATABASE_URL not set, using default: %s", dbConnString)
	}

	// 2. Initialize Database
	ctx := context.Background()
	store, err := hq.NewDBStore(ctx, dbConnString)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer store.Close()

	if err := store.Init(ctx); err != nil {
		log.Fatalf("Failed to init DB schema: %v", err)
	}
	log.Println("Database connection established and schema initialized.")
}
