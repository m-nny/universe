package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/cmd/internal"
	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/track"
	"github.com/m-nny/universe/ent/user"
	"github.com/m-nny/universe/lib/spotify"
)

const username = "m-nny"

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
	ctx := context.Background()
	entClient, err := internal.GetEntClient()
	if err != nil {
		log.Fatalf("failed creating ent Client: %v", err)
	}
	spotify, err := spotify.New(ctx, entClient, username)
	if err != nil {
		log.Fatalf("Error getting spotify client: %v", err)
	}

	playlists, err := spotify.GetAllPlaylists(ctx)
	if err != nil {
		log.Fatalf("Error getting all playlists: %v", err)
	}
	log.Printf("found total %d playlists", len(playlists))

	tracks, err := spotify.GetAllTracks(ctx)
	if err != nil {
		log.Fatalf("Error getting all tracks: %v", err)
	}
	log.Printf("found total %d tracks", len(tracks))

	savedTracks, err := entClient.Track.
		Query().
		WithAlbum().
		WithArtists().
		First(ctx)
	if err != nil {
		log.Fatalf("Error getting all tracks: %v", err)
	}
	log.Printf("tracks: %+v", savedTracks)
}

func getTracks(ctx context.Context, ent *ent.Client) error {
	tracks, err := ent.Track.
		Query().
		Where(track.HasSavedByWith(user.ID(username))).
		All(ctx)
	if err != nil {
		return err
	}
	log.Printf("tracks: %+v", tracks)
	return nil
}
