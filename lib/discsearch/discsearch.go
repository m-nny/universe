package discsearch

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/brain"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/spotify"
	"github.com/m-nny/universe/lib/spotify_ent"
)

type App struct {
	Ent        *ent.Client
	Brain      *brain.Brain
	Spotify    *spotify.Service
	SpotifyEnt *spotify_ent.Service
	Discogs    *discogs.Service
}

func New(ctx context.Context, username string) (*App, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("could not load .env: %v", err)
	}
	ent, err := getEntClient()
	if err != nil {
		return nil, err
	}
	brain, err := getBrain()
	if err != nil {
		return nil, err
	}
	spotify, err := spotify.New(ctx, brain, username)
	if err != nil {
		return nil, err
	}
	spotifyEnt := spotify_ent.New(ctx, ent, username)
	discogs, err := discogs.New()
	if err != nil {
		return nil, err
	}
	return &App{
		Ent:        ent,
		Spotify:    spotify,
		SpotifyEnt: spotifyEnt,
		Discogs:    discogs,
		Brain:      brain,
	}, nil
}

func getDbPath(db string) (string, error) {
	databasePath := fmt.Sprintf("data/%s.db", db)
	databasePath, err := filepath.Abs(databasePath)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(databasePath); err != nil {
		log.Printf("err: %v", err)
		return "", err
	}
	return databasePath, nil
}

func getEntClient() (*ent.Client, error) {
	databasePath, err := getDbPath("ent")
	if err != nil {
		return nil, err
	}
	client, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&_fk=1", databasePath))
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %w", err)
	}
	// Run the auto migration tool.
	if err := client.Schema.Create(context.Background()); err != nil {
		return nil, fmt.Errorf("failed creating schema resources: %w", err)
	}
	return client, nil
}

func getBrain() (*brain.Brain, error) {
	databasePath, err := getDbPath("gorm")
	if err != nil {
		return nil, err
	}
	return brain.New(databasePath /*enableLogging=*/, false)
}
