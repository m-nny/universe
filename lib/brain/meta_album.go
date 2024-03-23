package brain

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/spotify/utils"
)

type MetaAlbum struct {
	gorm.Model
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
