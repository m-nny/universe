package discsearch

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/spotify"
	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type App struct {
	Ent     *ent.Client
	Spotify *spotify.Service
	Discogs *discogs.Service
}

func New(ctx context.Context, username string) (*App, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
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

func getDbPath() (string, error) {
	databasePath := "data/ent.db"
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
	databasePath, err := getDbPath()
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

func getGormClient() (*gorm.DB, error) {
	databasePath, err := getDbPath()
	if err != nil {
		return nil, err
	}
	db, err := gorm.Open(sqlite.Open(databasePath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(); err != nil {
		return nil, err
	}
	return db, nil
}
