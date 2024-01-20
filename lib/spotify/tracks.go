package spotify

import (
	"context"
	"errors"
	"log"

	"github.com/m-nny/universe/lib/jsoncache"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

var _tracks = make(map[spotify.ID]*Track)

type Track struct {
	Id          spotify.ID
	Name        string
	ArtistNames []string
	ArtistIds   []spotify.ID
	Album       *Album
}

func getTrack(savedTrack spotify.SavedTrack) *Track {
	id := savedTrack.ID
	if val := _tracks[id]; val != nil {
		return val
	}
	artistNames, artistIds := getArtistNames(savedTrack.Artists)
	album := getAlbum(savedTrack.Album)
	track := &Track{
		Id:          savedTrack.ID,
		Name:        savedTrack.Name,
		ArtistNames: artistNames,
		ArtistIds:   artistIds,
		Album:       album,
	}
	return track
}

func (s *SpotfyClient) GetAllTracks(ctx context.Context) ([]*Track, error) {
	return jsoncache.CachedExec("spotify_savedTracks", func() ([]*Track, error) {
		return s._GetAllTracks(ctx)
	})
}
func (s *SpotfyClient) _GetAllTracks(ctx context.Context) ([]*Track, error) {
	var allTracks []*Track
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
