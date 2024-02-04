package spotify

import (
	"context"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) toAlbum(t spotify.SimpleAlbum) *ent.AlbumCreate {
	artistNames, artistIds := getArtistNames(t.Artists)
	return s.ent.Album.Create().
		SetID(string(t.ID)).
		SetName(string(t.Name)).
		SetArtistNames(artistNames).
		SetArtistIds(artistIds)
}

func (s *Service) saveTrackAlbums(ctx context.Context, rawPlists []spotify.SavedTrack) error {
	rawPlists = utils.SliceUniqe(rawPlists, func(item spotify.SavedTrack) string {
		return string(item.ID)
	})
	albums := utils.SliceMap(rawPlists, func(p spotify.SavedTrack) *ent.AlbumCreate {
		return s.toAlbum(p.Album)
	})
	err := s.ent.Album.
		CreateBulk(albums...).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
	return err
}
