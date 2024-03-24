package brain

import (
	"fmt"

	"github.com/zmb3/spotify/v2"

	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

type MetaAlbum struct {
	ID             uint `gorm:"primarykey"`
	SimplifiedName string
	Artists        []*Artist `gorm:"many2many:meta_album_artists;"`
}

func (s *MetaAlbum) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", s.SimplifiedName)
}

func newMetaAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist) *MetaAlbum {
	return &MetaAlbum{
		SimplifiedName: utils.SimplifiedAlbumName(sAlbum),
		Artists:        bArtists,
	}
}

func (b *Brain) MetaAlbumCount() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaAlbum{}).Count(&count).Error
	return int(count), err
}
