package spotify

import (
	"context"

	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toArtist(ctx context.Context, a spotify.SimpleArtist) (int, error) {
	// Check if already have it
	artist, err := s.ent.Artist.
		Query().
		Where(artist.SpotifyId(string(a.ID))).
		Only(ctx)
	if err == nil {
		return artist.ID, nil
	}

	artist, err = s.ent.Artist.
		Create().
		SetSpotifyId(string(a.ID)).
		SetName(a.Name).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return artist.ID, nil
}

func (s *Service) toArtists(ctx context.Context, rawArtists []spotify.SimpleArtist) ([]int, error) {
	return utils.SliceMapCtxErr(ctx, rawArtists, s.toArtist)
}
