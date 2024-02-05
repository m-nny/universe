package spotify

import (
	"context"
	"errors"
	"log"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/user"
	"github.com/m-nny/universe/lib/jsoncache"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) GetUserTracks(ctx context.Context) ([]*ent.Track, error) {
	// tracks, err := s._GetUserTracks(ctx)
	// if err == nil && len(tracks) > 0 {
	// 	return tracks, nil
	// }
	rawTracks, err := jsoncache.CachedExec("spotify_savedTracks", func() ([]spotify.SavedTrack, error) {
		var rawTracks []spotify.SavedTrack
		resp, err := s.spotify.CurrentUsersTracks(ctx,
			spotify.Limit(50),
			// spotify.Offset(3300),
		)
		for ; err == nil; err = s.spotify.NextPage(ctx, resp) {
			log.Printf("len(resp.Tracks)=%d offest=%d total=%d", len(resp.Tracks), resp.Offset, resp.Total)
			rawTracks = append(rawTracks, resp.Tracks...)
		}
		if !errors.Is(err, spotify.ErrNoMorePages) && err != nil {
			return nil, err
		}
		return rawTracks, nil
	})
	if err != nil {
		return nil, err
	}
	return utils.SliceMapCtxErr(ctx, rawTracks, s.toTrackSaved)
}

func (s *Service) _GetUserTracks(ctx context.Context) ([]*ent.Track, error) {
	plists, err := s.ent.User.
		Query().
		Where(user.ID(s.username)).
		QuerySavedTracks().
		All(ctx)
	return plists, err
}