package discsearch

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"

	"github.com/m-nny/universe/lib/brain"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/spotify"
)

type App struct {
	Brain   *brain.Brain
	Spotify *spotify.Service
	Discogs *discogs.Service
}

func New(ctx context.Context, username string, offlineMode bool) (*App, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("could not load .env: %v", err)
	}
	brain, err := getBrain()
	if err != nil {
		return nil, err
	}
	spotify, err := spotify.New(ctx, brain, username, offlineMode)
	if err != nil {
		return nil, err
	}
	discogs, err := discogs.New()
	if err != nil {
		return nil, err
	}
	return &App{
		Spotify: spotify,
		Discogs: discogs,
		Brain:   brain,
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

func getBrain() (*brain.Brain, error) {
	databasePath, err := getDbPath("gorm")
	if err != nil {
		return nil, err
	}
	return brain.New(databasePath /*enableLogging=*/, false)
}
