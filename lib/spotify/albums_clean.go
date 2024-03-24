package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func (s *Service) GetAlbumsById(ctx context.Context, ids []spotify.ID) ([]*spotify.FullAlbum, error) {
	sAlbums, err := s.spotify.GetAlbums(ctx, ids)
	if err != nil {
		return nil, err
	}
	return sAlbums, nil
}
