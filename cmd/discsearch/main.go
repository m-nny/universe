package main

import (
	"context"
	"log"

	"github.com/joho/godotenv"
	"github.com/m-nny/universe/lib/spotify"
	"github.com/m-nny/universe/lib/utils"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Could not load .env: %v", err)
	}
	ctx := context.Background()
	spotifyClient, err := spotify.New(ctx)
	if err != nil {
		log.Fatalf("Error getting spotify client: %v", err)
	}
	tracks, err := spotifyClient.GetAllTracks(ctx)
	if err != nil {
		log.Fatalf("Error getting all Tracks: %v", err)
	}
	log.Printf("found total %d Tracks", len(tracks))
	GetTopArtistsAndAlbums(tracks)
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
