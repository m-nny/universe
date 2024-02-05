package main

import (
	"context"
	"log"

	"entgo.io/ent/dialect/sql"
	"github.com/joho/godotenv"
	"github.com/m-nny/universe/cmd/internal"
	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/spotify"
)

func main() {
	username := "m-nny"
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

	if err := getTopAlbums(ctx, entClient, username); err != nil {
		log.Fatalf("Error getting all tracks: %v", err)
	}
	if err := getTopArtists(ctx, entClient, username); err != nil {
		log.Fatalf("Error getting all tracks: %v", err)
	}
}

func getTopAlbums(ctx context.Context, entC *ent.Client, username string) error {
	const tracksNumCol = "tracks_num"
	albums, err := entC.Album.
		Query().
		Order(
			album.ByTracksCount(
				sql.OrderDesc(),
				sql.OrderSelectAs(tracksNumCol),
			),
		).
		Limit(10).
		All(ctx)
	if err != nil {
		return err
	}
	log.Print("top albums:")
	for _, album := range albums {
		tracksNum, err := album.Value(tracksNumCol)
		if err != nil {
			return err
		}
		log.Printf("name: %s tracks_num: %d", album.Name, tracksNum)
	}
	return nil
}

func getTopArtists(ctx context.Context, entC *ent.Client, username string) error {
	const tracksNumCol = "tracks_num"
	artists, err := entC.Artist.
		Query().
		Order(
			artist.ByTracksCount(
				sql.OrderDesc(),
				sql.OrderSelectAs(tracksNumCol),
			),
		).
		Limit(10).
		All(ctx)
	if err != nil {
		return err
	}
	log.Print("top artists:")
	for _, artist := range artists {
		tracksNum, err := artist.Value(tracksNumCol)
		if err != nil {
			return err
		}
		log.Printf("name: %s tracks_num: %d", artist.Name, tracksNum)
	}
	return nil
}
