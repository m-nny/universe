package spotify

import (
	"context"

	"github.com/zmb3/spotify/v2"
)

func (s *Service) toAlbum(ctx context.Context, t spotify.SimpleAlbum) (string, error) {
	artistIds, err := s.toArtists(ctx, t.Artists)
	if err != nil {
		return "", err
	}
	return s.ent.Album.Create().
		SetID(string(t.ID)).
		SetName(string(t.Name)).
		AddArtistIDs(artistIds...).
		OnConflict().
		UpdateNewValues().
		ID(ctx)
}

// func (s *Service) toAlbums(ctx context.Context, rawAlbums []spotify.SimpleAlbum) ([]string, error) {
// 	return utils.SliceMapCtxErr(ctx, rawAlbums, s.toAlbum)
// }
