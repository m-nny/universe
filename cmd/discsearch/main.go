package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"time"

	"github.com/m-nny/universe/lib/discsearch"
	"github.com/m-nny/universe/lib/spotify"
	"github.com/m-nny/universe/lib/utils/logutils"
)

var (
	offlineMode = flag.Bool("offline", false, "Offline mode")
	loggerLevel = logutils.Level("log", slog.LevelInfo, "slog level")
)

const username = "m-nny"

func main() {
	flag.Parse()
	slog.SetLogLoggerLevel(*loggerLevel)
	ctx := context.Background()
	app, err := discsearch.New(ctx, username, *offlineMode)
	if err != nil {
		logutils.Fatal("Could not init app", "err", err)
	}
	defer func() {
		if err := app.Close(); err != nil {
			slog.Error("Error closing app", "err", err)
		}
	}()

	// if err := getAlbumsById(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := benchGetUserTracks(ctx, app); err != nil {
	// 	logutils.Fatal("Could not bench GetUserTracks", "err", err)
	// }

	// if err := getDiscogs(ctx, app); err != nil {
	// 	logutils.Fatal("Could not get discogs", "err", err)
	// }

	if err := getSellerInventory(ctx, app); err != nil {
		logutils.Fatal("Could not get seller inventory", "err", err)
	}

	slog.Info("Done")
}

func getDiscogs(ctx context.Context, app *discsearch.App) error {
	if _, err := app.Discogs.SellerInventory(ctx, "nezrathebeatmaker"); err != nil {
		return err
	}
	return nil
}

func getAlbumsById(ctx context.Context, app *discsearch.App) error {
	// Hybrid Theory, Hybrid Theory (20th Edition)
	albumIds := []spotify.ID{"6PFPjumGRpZnBzqnDci6qJ", "28DUZ0itKISf2sr6hlseMy"}
	targetAlbums, err := app.Spotify.GetAlbumsById(ctx, albumIds)
	if err != nil {
		return err
	}
	slog.Info("targetAlbums", "targetAlbums", targetAlbums)

	// if err := getTopAlbums(ctx, app); err != nil {
	// 	return err
	// }
	// if err := getTopArtists(ctx, app); err != nil {
	// 	return err
	// }
	return nil
}

func benchGetUserTracks(ctx context.Context, app *discsearch.App) error {
	var start time.Time
	start = time.Now()
	userTracks, err := app.Spotify.GetUserTracks(ctx, username)
	if err != nil {
		return err
	}
	logutils.Infof("brain.GetUserTracks: finished in %s", time.Since(start))

	logutils.Infof("==========================")
	logutils.Infof("brain.SaveTracksSqlx")
	start = time.Now()
	sqlxTracks, err := app.Brain.SaveTracksSqlx(userTracks, username)
	if err != nil {
		return fmt.Errorf("error getting all tracks: %w", err)
	}
	logutils.Infof("finished in %s", time.Since(start))
	logutils.Infof("returned %d tracks", len(sqlxTracks))

	sqlxTrackCnt, err := app.Brain.MetaTrackCountSqlx()
	if err != nil {
		return err
	}
	logutils.Infof("track cnt in db: %d", sqlxTrackCnt)
	sqlxAlbumCnt, err := app.Brain.MetaAlbumCountSqlx()
	if err != nil {
		return err
	}
	logutils.Infof("album cnt in db: %d", sqlxAlbumCnt)

	return nil
}

func getSellerInventory(ctx context.Context, app *discsearch.App) error {
	sellerId := "nezrathebeatmaker"
	// sellerId := "TheRecordAlbum"
	if _, err := app.Inventory(ctx, sellerId); err != nil {
		return err
	}
	return nil
}
