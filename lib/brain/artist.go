package brain

import (
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type Artist struct {
	gorm.Model
	SpotifyId string
	Name      string
}

func (b *Brain) ToArtist(sArtist *spotify.FullArtist) (*Artist, error) {
	var artist Artist
	if err := b.gormDb.
		Where(&Artist{SpotifyId: sArtist.ID.String()}).
		Attrs(Artist{Name: sArtist.Name, SpotifyId: sArtist.ID.String()}).
		FirstOrCreate(&artist).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}
