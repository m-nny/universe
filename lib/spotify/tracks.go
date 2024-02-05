package spotify

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/track"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toTrack(ctx context.Context, t spotify.SavedTrack) (*ent.Track, error) {
	simplifiedName := simplifiedTrackName(t)
	track, err := s.ent.Track.
		Query().
		Where(track.Similar(string(t.ID), simplifiedName)).
		Only(ctx)
	if err != nil {
		// Track does not exist yet
		return s._newTrack(ctx, t, simplifiedName)
	}
	if slices.Contains(track.SpotifyIds, string(t.ID)) {
		return track, nil
	}
	log.Printf("new track version: cur: %v new: %v", track, t.ID)
	return track.Update().AppendSpotifyIds([]string{string(t.ID)}).Save(ctx)
}

func (s *Service) _newTrack(ctx context.Context, t spotify.SavedTrack, simplifiedName string) (*ent.Track, error) {
	artistIds, err := s.toArtists(ctx, t.Artists)
	if err != nil {
		return nil, err
	}
	album, err := s.toAlbum(ctx, t.Album)
	if err != nil {
		return nil, err
	}
	return s.ent.Track.
		Create().
		AddArtistIDs(artistIds...).
		AddSavedByIDs(s.username).
		SetAlbum(album).
		SetName(string(t.Name)).
		SetSimplifiedName(simplifiedName).
		SetSpotifyIds([]string{string(t.ID)}).
		SetTrackNumber(t.TrackNumber).
		Save(ctx)
}

// simplifiedTrackName will return a string in form of
//
//	"<artist1>, <artist2> - <album name> [<album release year] - <track_num>. <track_name>"
func simplifiedTrackName(t spotify.SavedTrack) string {
	msg := simplifiedAlbumName(t.Album)
	msg += fmt.Sprintf("%d.  %s", t.TrackNumber, t.Name)
	msg = strings.ToLower(msg)
	return msg
}
