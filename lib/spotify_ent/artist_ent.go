package spotify_ent

import (
	"context"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/utils/hitcounter"
	"github.com/m-nny/universe/lib/utils/maputils"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

var (
	artistNaiveHc = hitcounter.New("ArtistNaive")
	artistBatchHc = hitcounter.New("ArtistBatch")
)

func (s *Service) toArtists(ctx context.Context, rawArtists []spotify.SimpleArtist) ([]int, error) {
	res, err := s.batchToArtists(ctx, rawArtists)
	if err != nil {
		return nil, err
	}
	return maputils.Values(res), nil
}

func (s *Service) batchToArtists(ctx context.Context, rawArtists []spotify.SimpleArtist) (map[spotify.ID]int, error) {
	rawArtists = sliceutils.Unique(rawArtists, func(item spotify.SimpleArtist) string { return item.ID.String() })
	spotifyIds := sliceutils.Map(rawArtists, func(item spotify.SimpleArtist) string { return item.ID.String() })
	existingArtists, err := s.ent.Artist.Query().Where(artist.SpotifyIdIn(spotifyIds...)).All(ctx)
	if err != nil {
		return nil, err
	}
	artistBatchHc.HitN(len(existingArtists))
	res := make(map[spotify.ID]int)
	for _, artist := range existingArtists {
		res[spotify.ID(artist.SpotifyId)] = artist.ID
	}
	var createArtists []*ent.ArtistCreate
	for _, rawArtist := range rawArtists {
		// Artist already in DB -> already in res map
		if _, ok := res[rawArtist.ID]; ok {
			continue
		}
		newArtist := s.ent.Artist.Create().SetSpotifyId(rawArtist.ID.String()).SetName(rawArtist.Name)
		createArtists = append(createArtists, newArtist)
	}
	// Early return, if there are no artists to create
	if len(createArtists) == 0 {
		return res, nil
	}
	createdArtists, err := s.ent.Artist.CreateBulk(createArtists...).Save(ctx)
	if err != nil {
		return nil, err
	}
	artistBatchHc.MissN(len(existingArtists))
	for _, artist := range createdArtists {
		res[spotify.ID(artist.SpotifyId)] = artist.ID
	}
	return res, nil
}
