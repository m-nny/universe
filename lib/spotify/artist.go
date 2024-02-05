package spotify

import (
	"context"

	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toArtist(ctx context.Context, t spotify.SimpleArtist) (int, error) {
	// Check if already have it
	artistId, err := s.ent.Artist.
		Query().
		Where(artist.SpotifyId(string(t.ID))).
		OnlyID(ctx)
	if err == nil {
		return artistId, nil
	}

	artist, err := s.ent.Artist.
		Create().
		SetSpotifyId(string(t.ID)).
		SetName(t.Name).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return artist.ID, nil
}

func (s *Service) toArtists(ctx context.Context, rawArtists []spotify.SimpleArtist) ([]int, error) {
	return utils.SliceMapCtxErr(ctx, rawArtists, s.toArtist)
}
