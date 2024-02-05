package spotify

import (
	"context"
	"errors"
	"log"

	"github.com/m-nny/universe/ent"
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
		spotify.Offset(3300),
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

func (s *Service) saveUserTracks(ctx context.Context, rawTracks []spotify.SavedTrack) ([]string, error) {
	return utils.SliceMapCtxErr(ctx, rawTracks, s.toTrack)
}

func (s *Service) toTrack(ctx context.Context, t spotify.SavedTrack) (string, error) {
	artistIds, err := s.toArtists(ctx, t.Artists)
	if err != nil {
		return "", err
	}
	albumId, err := s.toAlbum(ctx, t.Album)
	if err != nil {
		return "", err
	}
	return s.ent.Track.
		Create().
		SetID(string(t.ID)).
		SetName(string(t.Name)).
		AddArtistIDs(artistIds...).
		SetAlbumID(albumId).
		AddSavedByIDs(s.username).
		OnConflict().
		UpdateNewValues().
		ID(ctx)
}
