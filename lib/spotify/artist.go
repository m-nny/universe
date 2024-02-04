package spotify

import (
	"context"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toArtist(t spotify.SimpleArtist) *ent.ArtistCreate {
	return s.ent.Artist.Create().
		SetID(string(t.ID)).
		SetName(string(t.Name))
}

func (s *Service) saveTrackArtists(ctx context.Context, rawTracks []spotify.SavedTrack) error {
	rawArtists := utils.SliceFlatMap(rawTracks, func(p spotify.SavedTrack) []spotify.SimpleArtist {
		return append(p.Artists, p.Album.Artists...)
	})
	rawArtists = utils.SliceUniqe(rawArtists, func(item spotify.SimpleArtist) string { return string(item.ID) })
	artists := utils.SliceMap(rawArtists, s.toArtist)
	err := s.ent.Artist.
		CreateBulk(artists...).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
	return err
}
