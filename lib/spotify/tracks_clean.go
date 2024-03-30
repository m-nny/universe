package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

func (s *Service) GetTracksById(ctx context.Context, ids []spotify.ID) ([]*spotify.SimpleTrack, error) {
	if s.offlineMode {
		return nil, ErrOffileMode
	}
	sFullTracks, err := s.spotify.GetTracks(ctx, ids)
	if err != nil {
		return nil, err
	}
	sTracks := sliceutils.Map(sFullTracks, func(item *spotify.FullTrack) *spotify.SimpleTrack { return &item.SimpleTrack })
	return sTracks, nil
}
