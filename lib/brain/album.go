package brain

import (
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Album struct {
	gorm.Model
	SpotifyId string
	Name      string
	Artists   []*Artist `gorm:"many2many:album_artists;"`
}

func newAlbum(sAlbum *spotify.FullAlbum) *Album {
	return &Album{Name: sAlbum.Name, SpotifyId: sAlbum.ID.String()}
}

// ToArtists returns Brain representain of a spotify album
//   - NOTE: Does not debupe based on simplified name
//   - NOTE: Does not associate Album with Artist
func (b *Brain) ToAlbum(sAlbum *spotify.FullAlbum) (*Album, error) {
	var artist Album
	if err := b.gormDb.
		Where(&Album{SpotifyId: sAlbum.ID.String()}).
		Attrs(newAlbum(sAlbum)).
		FirstOrCreate(&artist).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

func (b *Brain) ToAlbums(sAlbums []*spotify.FullAlbum) ([]*Album, error) {
	return sliceutils.MapErr(sAlbums, b.ToAlbum)
}
