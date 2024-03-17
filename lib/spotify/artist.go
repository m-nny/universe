package spotify

import (
	"context"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/artist"
	"github.com/m-nny/universe/lib/utils/slices"
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
	return slices.MapCtxErr(ctx, rawArtists, s.toArtist)
}

func (s *Service) batchToArtists(ctx context.Context, rawArtists []spotify.SimpleArtist) (map[spotify.ID]int, error) {
	rawArtists = slices.Uniqe(rawArtists, func(item spotify.SimpleArtist) string { return item.ID.String() })
	spotifyIds := slices.Map(rawArtists, func(item spotify.SimpleArtist) string { return item.ID.String() })
	existingArtists, err := s.ent.Artist.Query().Where(artist.SpotifyIdIn(spotifyIds...)).All(ctx)
	if err != nil {
		return nil, err
	}
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
	for _, artist := range createdArtists {
		res[spotify.ID(artist.SpotifyId)] = artist.ID
	}
	return res, nil
}

func ArtistsString(artists []*ent.Artist) string {
	var s []string
	for _, a := range artists {
		s = append(s, a.Name)
	}
	return strings.Join(s, " ")
}
