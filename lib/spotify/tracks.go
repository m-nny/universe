package spotify

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/track"
	"github.com/m-nny/universe/lib/utils/sliceutils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toTracksSaved(ctx context.Context, tracks []spotify.SavedTrack, username string) ([]*ent.Track, error) {
	return sliceutils.MapCtxErr(ctx, tracks,
		func(ctx context.Context, t spotify.SavedTrack) (*ent.Track, error) {
			album, err := s.toAlbum(ctx, t.Album)
			if err != nil {
				return nil, err
			}
			return s.toTrackWithAlbum(ctx, t.SimpleTrack, album, username)
		})
}

func (s *Service) toTracksWithAlbum(ctx context.Context, tracks []spotify.SimpleTrack, a *ent.Album) ([]*ent.Track, error) {
	return sliceutils.MapCtxErr(ctx, tracks,
		func(ctx context.Context, t spotify.SimpleTrack) (*ent.Track, error) {
			return s.toTrackWithAlbum(ctx, t, a, "")
		})
}

func (s *Service) toTrackWithAlbum(ctx context.Context, t spotify.SimpleTrack, a *ent.Album, username string) (*ent.Track, error) {
	simplifiedName := simplifiedTrackName(t, a)
	track, err := s.ent.Track.
		Query().
		Where(track.Similar(string(t.ID), simplifiedName)).
		Only(ctx)
	if err != nil {
		// Track does not exist yet
		return s._newTrack(ctx, t, a, simplifiedName, username)
	}
	if slices.Contains(track.SpotifyIds, string(t.ID)) {
		return track, nil
	}
	// log.Printf("new track version: cur: %v new: %v", track, t.ID)
	return track.Update().AppendSpotifyIds([]string{string(t.ID)}).Save(ctx)
}

func (s *Service) _newTrack(ctx context.Context, t spotify.SimpleTrack, album *ent.Album, simplifiedName string, username string) (*ent.Track, error) {
	artistIds, err := s.toArtists(ctx, t.Artists)
	if err != nil {
		return nil, err
	}
	track := s.ent.Track.
		Create().
		AddArtistIDs(artistIds...).
		SetAlbum(album).
		SetName(t.Name).
		SetSimplifiedName(simplifiedName).
		SetSpotifyIds([]string{string(t.ID)}).
		SetTrackNumber(t.TrackNumber)
	if username != "" {
		track.AddSavedByIDs(username)
	}
	return track.Save(ctx)
}

// simplifiedTrackName will return a string in form of
//
//	"<artist1>, <artist2> - <album name> [<album release year] - <track_num>. <track_name>"
func simplifiedTrackName(t spotify.SimpleTrack, a *ent.Album) string {
	msg := a.SimplifiedName
	msg += fmt.Sprintf(" %d.  %s", t.TrackNumber, t.Name)
	msg = strings.ToLower(msg)
	return msg
}
