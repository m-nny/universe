package main

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect/sql"
	"github.com/joho/godotenv"
	"github.com/m-nny/universe/cmd/internal"
	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/spotify"
	spotify2 "github.com/zmb3/spotify/v2"
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

	if err := getAlbumsById(ctx, spotify, entClient); err != nil {
		log.Fatalf("%v", err)
	}

	// if err := getUserTracks(ctx, spotify, entClient); err != nil {
	// 	log.Fatalf("%v", err)
	// }
}

func getAlbumsById(ctx context.Context, spotify *spotify.Service, entClient *ent.Client) error {
	albumIds := []spotify2.ID{"025WnFQfYniZWzIzFHx0mb", "1svovXeaO67ZpSgWhj0UaP", "2ArGu1xrwGla8pZfTNOBfp"}
	targetAlbums, err := spotify.GetAlbumsById(ctx, albumIds)
	if err != nil {
		return err
	}
	log.Print("targetAlbums:")
	for i, album := range targetAlbums {
		log.Printf("%2d %s", i+1, album)
	}
	log.Print()

	if err := getTopAlbums(ctx, entClient); err != nil {
		return err
	}
	if err := getTopArtists(ctx, entClient); err != nil {
		return err
	}
	return nil
}

func getPlaylists(ctx context.Context, spotify *spotify.Service) error {
	playlists, err := spotify.GetAllPlaylists(ctx)
	if err != nil {
		return fmt.Errorf("Error getting all playlists: %w", err)
	}
	log.Printf("found total %d playlists", len(playlists))
	return nil
}

func getUserTracks(ctx context.Context, spotify *spotify.Service, entClient *ent.Client) error {
	tracks, err := spotify.GetUserTracks(ctx)
	if err != nil {
		return fmt.Errorf("Error getting all tracks: %w", err)
	}
	log.Printf("found total %d tracks", len(tracks))

	if err := getTopAlbums(ctx, entClient); err != nil {
		return fmt.Errorf("Error getting all tracks: %w", err)
	}
	if err := getTopArtists(ctx, entClient); err != nil {
		return fmt.Errorf("Error getting all tracks: %w", err)
	}
	return nil
}

func getTopAlbums(ctx context.Context, entC *ent.Client) error {
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

func getTopArtists(ctx context.Context, entC *ent.Client) error {
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
