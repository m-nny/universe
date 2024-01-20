package spotify

import (
	"context"
	"errors"
	"log"

	"github.com/m-nny/universe/lib/jsoncache"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

type Track struct {
	Id          spotify.ID
	Name        string
	ArtistNames []string
	ArtistIds   []spotify.ID
	AlbumName   string
	AlbumId     spotify.ID
}

func getTrack(track spotify.SavedTrack) Track {
	var artistNames []string
	var artistIds []spotify.ID
	for _, artist := range track.Artists {
		artistNames = append(artistNames, artist.Name)
		artistIds = append(artistIds, artist.ID)
	}
	return Track{
		Id:          track.ID,
		Name:        track.Name,
		ArtistNames: artistNames,
		ArtistIds:   artistIds,
		AlbumName:   track.Album.Name,
		AlbumId:     track.Album.ID,
	}
}

func (s *SpotfyClient) GetAllTracks(ctx context.Context) ([]Track, error) {
	return jsoncache.CachedExec("spotify_savedTracks", func() ([]Track, error) {
		return s._GetAllTracks(ctx)
	})
}
func (s *SpotfyClient) _GetAllTracks(ctx context.Context) ([]Track, error) {
	var allTracks []Track
	var err error
	for resp, err := s.client.CurrentUsersTracks(ctx, spotify.Limit(50)); err == nil; err = s.client.NextPage(ctx, resp) {
		log.Printf("len(resp.Tracks)=%d", len(resp.Tracks))
		allTracks = append(allTracks, utils.SliceMap(resp.Tracks, getTrack)...)
	}
	if !errors.Is(err, spotify.ErrNoMorePages) && err != nil {
		return nil, err
	}
	return allTracks, nil
}
