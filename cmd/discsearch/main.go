package main

import (
	"context"
	"fmt"
	"log"

	"entgo.io/ent/dialect/sql"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/discsearch"
	spotify2 "github.com/zmb3/spotify/v2"
)

const username = "m-nny"

func main() {
	ctx := context.Background()
	app, err := discsearch.New(ctx, username)
	if err != nil {
		log.Fatalf("Could not init app: %v", err)
	}

	// if err := getAlbumsById(ctx, spotify, entClient); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := getUserTracks(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := getDiscogs(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	if _, err := app.Inventory(ctx, "nezrathebeatmaker"); err != nil {
		log.Fatalf("%v", err)
	}
}

func getDiscogs(ctx context.Context, app *discsearch.App) error {
	if _, err := app.Discogs.SellerInventory(ctx, "nezrathebeatmaker"); err != nil {
		return err
	}
	return nil
}

func getAlbumsById(ctx context.Context, app *discsearch.App) error {
	albumIds := []spotify2.ID{"025WnFQfYniZWzIzFHx0mb", "1svovXeaO67ZpSgWhj0UaP", "2ArGu1xrwGla8pZfTNOBfp"}
	targetAlbums, err := app.Spotify.GetAlbumsById(ctx, albumIds)
	if err != nil {
		return err
	}
	log.Print("targetAlbums:")
	for i, album := range targetAlbums {
		log.Printf("%2d %s", i+1, album)
	}
	log.Print()

	if err := getTopAlbums(ctx, app); err != nil {
		return err
	}
	if err := getTopArtists(ctx, app); err != nil {
		return err
	}
	return nil
}

func getUserTracks(ctx context.Context, app *discsearch.App) error {
	tracks, err := app.Spotify.GetUserTracks(ctx, username)
	if err != nil {
		return fmt.Errorf("error getting all tracks: %w", err)
	}
	log.Printf("found total %d tracks", len(tracks))

	if err := getTopAlbums(ctx, app); err != nil {
		return fmt.Errorf("error getting all tracks: %w", err)
	}
	if err := getTopArtists(ctx, app); err != nil {
		return fmt.Errorf("error getting all tracks: %w", err)
	}
	return nil
}

func getTopAlbums(ctx context.Context, app *discsearch.App) error {
	const tracksNumCol = "tracks_num"
	albums, err := app.Ent.Album.
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

func getTopArtists(ctx context.Context, app *discsearch.App) error {
	const tracksNumCol = "tracks_num"
	artists, err := app.Ent.Artist.
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
