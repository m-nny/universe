package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/track"
	"github.com/m-nny/universe/ent/user"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) GetUserTracks(ctx context.Context) ([]*ent.Track, error) {
	tracks, err := s.getUserTracks(ctx)
	if err == nil && len(tracks) > 0 {
		return tracks, nil
	}
	rawTracks, err := s._GetUserTracks(ctx)
	if err != nil {
		return nil, err
	}
	if _, err := s.saveUserTracks(ctx, rawTracks); err != nil {
		return nil, err
	}
	return s.getUserTracks(ctx)
}

func (s *Service) _GetUserTracks(ctx context.Context) ([]spotify.SavedTrack, error) {
	var allTracks []spotify.SavedTrack
	resp, err := s.spotify.CurrentUsersTracks(ctx,
		spotify.Limit(50),
		// spotify.Offset(3300),
	)
	for ; err == nil; err = s.spotify.NextPage(ctx, resp) {
		log.Printf("len(resp.Tracks)=%d offest=%d total=%d", len(resp.Tracks), resp.Offset, resp.Total)
		allTracks = append(allTracks, resp.Tracks...)
	}
	if !errors.Is(err, spotify.ErrNoMorePages) && err != nil {
		return nil, err
	}
	return allTracks, nil
}

func (s *Service) getUserTracks(ctx context.Context) ([]*ent.Track, error) {
	plists, err := s.ent.User.
		Query().
		Where(user.ID(s.username)).
		QuerySavedTracks().
		All(ctx)
	return plists, err
}

func (s *Service) saveUserTracks(ctx context.Context, rawTracks []spotify.SavedTrack) ([]*ent.Track, error) {
	return utils.SliceMapCtxErr(ctx, rawTracks, s.toTrack)
}

func (s *Service) getTrackFromSpotify(ctx context.Context, t spotify.SavedTrack) (*ent.Track, error) {
	simplifiedName := simplifiedTrackName(t)
	return s.ent.Track.
		Query().
		Where(
			track.Or(
				track.SpotifyIdContains(string(t.ID)),
				track.SimplifiedName(simplifiedName),
			),
		).
		Only(ctx)
}

func (s *Service) toTrack(ctx context.Context, t spotify.SavedTrack) (*ent.Track, error) {
	track, err := s.getTrackFromSpotify(ctx, t)
	if err != nil {
		// Track does not exist yet
		return s._newTrack(ctx, t)
	}
	if slices.Contains(track.SpotifyIds, string(t.ID)) {
		return track, nil
	}
	log.Printf("new track version: cur: %v new: %v", track, t.ID)
	return track.Update().AppendSpotifyIds([]string{string(t.ID)}).Save(ctx)
}

func (s *Service) _newTrack(ctx context.Context, t spotify.SavedTrack) (*ent.Track, error) {
	artistIds, err := s.toArtists(ctx, t.Artists)
	if err != nil {
		return nil, err
	}
	album, err := s.toAlbum(ctx, t.Album)
	if err != nil {
		return nil, err
	}
	simplfiedName := simplifiedTrackName(t)
	return s.ent.Track.
		Create().
		AddArtistIDs(artistIds...).
		AddSavedByIDs(s.username).
		SetAlbum(album).
		SetName(string(t.Name)).
		SetSimplifiedName(simplfiedName).
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
