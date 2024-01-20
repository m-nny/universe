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
	spotifyClient, err := spotify.New(ctx)
	if err != nil {
		log.Fatalf("Error getting spotify client: %v", err)
	}
	tracks, err := spotifyClient.GetAllTracks(ctx)
	if err != nil {
		log.Fatalf("Error getting all Tracks: %v", err)
	}
	log.Printf("found total %d Tracks", len(tracks))
}
