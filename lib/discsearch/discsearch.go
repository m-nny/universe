package discsearch

import (
	"context"
	"fmt"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/spotify"
	_ "github.com/mattn/go-sqlite3"
)

type App struct {
	Ent     *ent.Client
	Spotify *spotify.Service
	Discogs *discogs.Service
}

func New(ctx context.Context, username string) (*App, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("could not load .env: %v", err)
	}
	ent, err := getEntClient()
	if err != nil {
		return nil, err
	}
	spotify, err := spotify.New(ctx, ent, username)
	if err != nil {
		return nil, err
	}
	discogs, err := discogs.New()
	if err != nil {
		return nil, err
	}
	return &App{
		Ent:     ent,
		Spotify: spotify,
		Discogs: discogs,
	}, nil
}

func getEntClient() (*ent.Client, error) {
	client, err := ent.Open("sqlite3", "file:data/ent.db?cache=shared&_fk=1")
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}
	return client, nil
}
