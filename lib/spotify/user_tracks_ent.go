package spotify

import (
	"context"

	"github.com/m-nny/universe/ent"
)

func (s *Service) GetUserTracksEnt(ctx context.Context, username string) ([]*ent.Track, error) {
	savedTracks, err := s.GetUserTracks(ctx, username)
	if err != nil {
		return nil, err
	}
	return s.ToTracksSaved(ctx, savedTracks, username)
}
