package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/lib/spotify"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
	ctx := context.Background()
	if err := spotify.GetAllTracks(ctx); err != nil {
		log.Fatalf("Error getting all tracks: %v", err)
	}
	log.Printf("discsearch finished")
}
