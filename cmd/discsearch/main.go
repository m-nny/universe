package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/cmd/internal"
	"github.com/m-nny/universe/lib/spotify"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
	ctx := context.Background()
	entClient, err := internal.GetEntClient()
	if err != nil {
		log.Fatalf("failed creating ent Client: %v", err)
	}
	spotifyConfig, err := spotify.LoadConfig()
	if err != nil {
		log.Fatalf("Error getting spotify config: %v", err)
	}
	spotifyClient, err := spotify.New(ctx, spotifyConfig, entClient)
	if err != nil {
		log.Fatalf("Error getting spotify client: %v", err)
	}

	playlists, err := spotifyClient.GetAllPlaylists(ctx)
	if err != nil {
		log.Fatalf("Error getting all playlists: %v", err)
	}
	log.Printf("found total %d playlists", len(playlists))

	tracks, err := spotifyClient.GetAllTracks(ctx)
	if err != nil {
		log.Fatalf("Error getting all tracks: %v", err)
	}
	log.Printf("found total %d tracks", len(tracks))

}
