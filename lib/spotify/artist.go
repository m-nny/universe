package spotify

import (
	"context"

	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toArtist(ctx context.Context, t spotify.SimpleArtist) (string, error) {
	return s.ent.Artist.
		Create().
		SetID(string(t.ID)).
		SetName(t.Name).
		OnConflict().
		UpdateNewValues().
		ID(ctx)
}

func (s *Service) toArtists(ctx context.Context, rawArtists []spotify.SimpleArtist) ([]string, error) {
	return utils.SliceMapCtxErr(ctx, rawArtists, s.toArtist)
}
