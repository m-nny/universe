package spotify

import (
	"context"
	"errors"
	"log"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/jsoncache"
)

func (s *Service) GetUserTracks(ctx context.Context, username string) ([]spotify.SavedTrack, error) {
	rawTracks, err := jsoncache.CachedExec("spotify/savedTracks/"+username, func() ([]spotify.SavedTrack, error) {
		if s.offlineMode {
			return nil, ErrOffileMode
		}
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
	return rawTracks, nil
}
