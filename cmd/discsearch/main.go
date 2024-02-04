package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/cmd/internal"
	"github.com/m-nny/universe/lib/spotify"
	"github.com/m-nny/universe/lib/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
	ctx := context.Background()
	entClient, err := internal.GetEntClient()
	if err != nil {
		log.Fatalf("failed creating ent Client: %v", err)
	}
	spotifyConfig, err := spotify.LoadConfig()
	if err != nil {
		log.Fatalf("Error getting spotify config: %v", err)
	}
	spotifyClient, err := spotify.New(ctx, spotifyConfig, entClient)
	if err != nil {
		log.Fatalf("Error getting spotify client: %v", err)
	}

	playlists, err := spotifyClient.GetAllPlaylists(ctx)
	if err != nil {
		log.Fatalf("Error getting all playlists: %v", err)
	}
	log.Printf("found total %d playlists", len(playlists))

	// tracks, err := spotifyClient.GetAllTracks(ctx)
	// if err != nil {
	// 	log.Fatalf("Error getting all Tracks: %v", err)
	// }
	// log.Printf("found total %d Tracks", len(tracks))
	// GetTopArtistsAndAlbums(tracks)
}

func GetTopArtistsAndAlbums(tracks []*spotify.Track) {
	albums := make(map[spotify.ID]int)
	for _, track := range tracks {
		albums[track.Album.Id]++
	}
	topAlbums := utils.TopValuesMap(albums)
	log.Print("Top albums")
	for idx, item := range topAlbums {
		cnt := albums[item]
		if cnt <= 1 {
			break
		}
		log.Printf("%3d. %s (%d)", idx+1, item, cnt)
	}
}
