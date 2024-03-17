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

func newAlbum(sAlbum *spotify.SimpleAlbum, bArtists []*Artist) *Album {
	return &Album{Name: sAlbum.Name, SpotifyId: sAlbum.ID.String(), Artists: bArtists}
}

// ToArtists returns Brain representain of a spotify album
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) ToAlbum(sAlbum *spotify.SimpleAlbum) (*Album, error) {
	bArtists, err := b.ToArtists(sliceutils.MapP(sAlbum.Artists))
	if err != nil {
		return nil, err
	}
	var album Album
	if err := b.gormDb.
		Preload("Artists").
		Where(&Album{SpotifyId: sAlbum.ID.String()}).
		Attrs(newAlbum(sAlbum, bArtists)).
		FirstOrCreate(&album).Error; err != nil {
		return nil, err
	}
	return &album, nil
}

func (b *Brain) ToAlbums(sAlbums []*spotify.SimpleAlbum) ([]*Album, error) {
	// TODO(m-nny): Batch album and artist creation
	return sliceutils.MapErr(sAlbums, b.ToAlbum)
}
