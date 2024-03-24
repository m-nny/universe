package brain

import (
	"fmt"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/spotify/utils"
)

type MetaAlbum struct {
	ID             uint `gorm:"primarykey"`
	SimplifiedName string
}

func (s *MetaAlbum) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", s.SimplifiedName)
}

func newMetaAlbum(sAlbum spotify.SimpleAlbum) *MetaAlbum {
	return &MetaAlbum{SimplifiedName: utils.SimplifiedAlbumName(sAlbum)}
}

func (b *Brain) MetaAlbumCount() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaAlbum{}).Count(&count).Error
	return int(count), err
}
