package brain

import (
	"fmt"

	"github.com/zmb3/spotify/v2"

	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

type MetaTrack struct {
	ID             uint `gorm:"primarykey"`
	SimplifiedName string
	MetaAlbumID    uint
	MetaAlbum      *MetaAlbum
}

func (s *MetaTrack) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", s.SimplifiedName)
}

func newMetaTrack(sTrack spotify.SimpleTrack, bMetaAlbum *MetaAlbum) *MetaTrack {
	return &MetaTrack{
		SimplifiedName: utils.SimplifiedTrackName(sTrack, bMetaAlbum.SimplifiedName),
		MetaAlbumID:    bMetaAlbum.ID,
		MetaAlbum:      bMetaAlbum,
	}
}

func (b *Brain) MetaTrackCount() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaTrack{}).Count(&count).Error
	return int(count), err
}
