package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func (s *Service) GetArtistById(ctx context.Context, ids []spotify.ID) ([]*spotify.FullArtist, error) {
	if s.offlineMode {
		return nil, ErrOffileMode
	}
	sArtists, err := s.spotify.GetArtists(ctx, ids...)
	if err != nil {
		return nil, err
	}
	return sArtists, nil
}
