package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/m-nny/universe/lib/discsearch"
	"github.com/m-nny/universe/lib/spotify"
	"github.com/m-nny/universe/lib/utils/spotifyutils"
)

var (
	offlineMode = flag.Bool("offline", false, "Offline mode")
)

const username = "m-nny"

func main() {
	flag.Parse()
	ctx := context.Background()
	app, err := discsearch.New(ctx, username, *offlineMode)
	if err != nil {
		log.Fatalf("Could not init app: %v", err)
	}

	// if err := demoGormTracks(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := getAlbumsById(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := gormGetUserTracks(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	if err := benchGetUserTracks(ctx, app); err != nil {
		log.Fatalf("%v", err)
	}

	// if err := benchGetUserTracks(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := getDiscogs(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }

	// if err := getSellerInventory(ctx, app); err != nil {
	// 	log.Fatalf("%v", err)
	// }
	log.Printf("Done")
}

func demoGormArtists(ctx context.Context, app *discsearch.App) error {
	// Porter Robinson, Linkin Park, Linkin Park
	artistIds := []spotify.ID{"3dz0NnIZhtKKeXZxLOxCam", "6XyY86QOPPrYVGvF9ch6wz", "6XyY86QOPPrYVGvF9ch6wz"}
	sArtists, err := app.Spotify.GetArtistById(ctx, artistIds)
	if err != nil {
		return err
	}
	for idx, artist := range sArtists {
		log.Printf("[%d/%d] sArtist: %+v", idx+1, len(sArtists), artist)
	}
	bArtists, err := app.Brain.SaveArtists(sArtists)
	if err != nil {
		return err
	}
	for idx, artist := range bArtists {
		log.Printf("[%d/%d] bArtist %+v", idx+1, len(sArtists), artist)
	}
	return nil
}

func demoGormAlbums(ctx context.Context, app *discsearch.App) error {
	// Hybrid Theory, Hybrid Theory (20th Edition)
	albumIds := []spotify.ID{"6PFPjumGRpZnBzqnDci6qJ", "28DUZ0itKISf2sr6hlseMy", "28DUZ0itKISf2sr6hlseMy"}
	sAlbums, err := app.Spotify.GetAlbumsById(ctx, albumIds)
	if err != nil {
		return err
	}
	for idx, sAlbum := range sAlbums {
		log.Printf("[%d/%d] sAlbum: %+v - %+v", idx+1, len(sAlbums), spotifyutils.SArtistsString(sAlbum.Artists), sAlbum.Name)
	}
	bAlbums, err := app.Brain.SaveAlbumsSqlx(sAlbums)
	if err != nil {
		return err
	}
	for idx, bAlbum := range bAlbums {
		log.Printf("[%d/%d] bAlbum %+v", idx+1, len(sAlbums), bAlbum)
	}
	return nil
}

func demoGormTracks(ctx context.Context, app *discsearch.App) error {
	// Hybrid Theory, Hybrid Theory (20th Edition)
	sTracks, err := app.Spotify.GetUserTracks(ctx, username)
	sTracks = sTracks[:10]
	if err != nil {
		return err
	}
	for idx, sTrack := range sTracks {
		log.Printf("[%d/%d] sTrack: %s - %s - %+s", idx+1, len(sTracks), spotifyutils.SArtistsString(sTrack.Artists), sTrack.Album.Name, sTrack.Name)
	}
	bTracks, err := app.Brain.SaveTracksGorm(sTracks, username)
	if err != nil {
		return err
	}
	for idx, bTrack := range bTracks {
		log.Printf("[%d/%d] bTrack: %+v - %+v - %+v", idx+1, len(sTracks), bTrack.Artists, bTrack.MetaAlbum, bTrack.SimplifiedName)
	}
	return nil
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
	log.Print("targetAlbums:")
	for i, album := range targetAlbums {
		log.Printf("%2d %+v", i+1, album)
	}
	log.Print()

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
	log.Printf("GetUserTracks: finished in %s", time.Since(start))

	// log.Printf("==========================")
	// log.Printf("ent.ToTracksSaved")
	// start = time.Now()
	// entTracks, err := app.SpotifyEnt.ToTracksSaved(ctx, userTracks, username)
	// if err != nil {
	// 	return fmt.Errorf("error getting all tracks: %w", err)
	// }
	// log.Printf("finished in %s", time.Since(start))
	// log.Printf("returned %d tracks", len(entTracks))
	// entTrackCnt, err := app.SpotifyEnt.EntTrackCount(ctx)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("track cnt in db: %d", entTrackCnt)
	// entAlbumCnt, err := app.SpotifyEnt.EntAlbumCount(ctx)
	// if err != nil {
	// 	return err
	// }
	// log.Printf("album cnt in db: %d", entAlbumCnt)

	log.Printf("==========================")
	log.Printf("brain.SaveTracks")
	start = time.Now()
	brainTracks, err := app.Brain.SaveTracksGorm(userTracks, username)
	if err != nil {
		return fmt.Errorf("error getting all tracks: %w", err)
	}
	log.Printf("finished in %s", time.Since(start))
	log.Printf("returned %d tracks", len(brainTracks))

	brainTrackCnt, err := app.Brain.MetaTrackCountGorm()
	if err != nil {
		return err
	}
	log.Printf("track cnt in db: %d", brainTrackCnt)
	brainAlbumCnt, err := app.Brain.MetaAlbumCount()
	if err != nil {
		return err
	}
	log.Printf("album cnt in db: %d", brainAlbumCnt)

	// if entTrackCnt != brainTrackCnt {
	// 	return fmt.Errorf("Different ent and brain track counts: ent %d brain %d", brainTrackCnt, entTrackCnt)
	// }

	// if entAlbumCnt != brainAlbumCnt {
	// 	return fmt.Errorf("Different ent and brain album counts: ent %d brain %d", brainAlbumCnt, entAlbumCnt)
	// }

	return nil
}

func getSellerInventory(ctx context.Context, app *discsearch.App) error {
	// sellerId := "nezrathebeatmaker"
	sellerId := "TheRecordAlbum"
	if _, err := app.Inventory(ctx, sellerId); err != nil {
		return err
	}
	return nil
}
