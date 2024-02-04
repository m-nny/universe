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

func (s *Service) GetAllTracks(ctx context.Context) ([]*ent.Track, error) {
	tracks, err := s.getUserTracks(ctx)
	if err == nil && len(tracks) > 0 {
		return tracks, nil
	}
	rawPlists, err := s._GetAllTracks(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.saveUserTracks(ctx, rawPlists); err != nil {
		return nil, err
	}
	return s.getUserTracks(ctx)

}
func (s *Service) _GetAllTracks(ctx context.Context) ([]spotify.SavedTrack, error) {
	var allTracks []spotify.SavedTrack
	var err error
	for resp, err := s.spotify.CurrentUsersTracks(ctx, spotify.Limit(50)); err == nil; err = s.spotify.NextPage(ctx, resp) {
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
		Where(user.ID(rootUserName)).
		QuerySavedTracks().
		All(ctx)
	return plists, err
}

func (s *Service) saveUserTracks(ctx context.Context, rawPlists []spotify.SavedTrack) error {
	ps := utils.SliceMap(rawPlists, s.toTrack)
	err := s.ent.Track.
		CreateBulk(ps...).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
	return err
}

func (s *Service) toTrack(t spotify.SavedTrack) *ent.TrackCreate {
	artistNames, artistIds := getArtistNames(t.Artists)
	// album := getAlbum(savedTrack.Album)
	return s.ent.Track.Create().
		SetID(string(t.ID)).
		SetName(string(t.Name)).
		SetArtistNames(artistNames).
		SetArtistIds(artistIds).
		AddSavedByIDs(rootUserName)
}
