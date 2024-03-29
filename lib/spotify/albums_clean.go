package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/jsoncache"
)

func (s *Service) GetAlbumsById(ctx context.Context, ids []spotify.ID) ([]*spotify.FullAlbum, error) {
	sAlbums, err := s.spotify.GetAlbums(ctx, ids)
	if err != nil {
		return nil, err
	}
	return sAlbums, nil
}

func (s *Service) SearchAlbum(ctx context.Context, q string, searchSize int) ([]spotify.SimpleAlbum, error) {
	results, err := jsoncache.CachedExec("spotify_search/"+q, func() (*spotify.SearchResult, error) {
		return s.spotify.Search(ctx, q, spotify.SearchTypeAlbum, spotify.Limit(searchSize))
	})
	if err != nil {
		return nil, err
	}
	return results.Albums.Albums, nil
}
