package brain

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
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

func upsertMetaAlbums(b *Brain, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*MetaAlbum, error) {
	sAlbums = sliceutils.Unique(sAlbums, utils.SimplifiedAlbumName)
	var albumSimps []string
	for _, sAlbum := range sAlbums {
		albumSimps = append(albumSimps, utils.SimplifiedAlbumName(sAlbum))
	}

	var existingMetaAlbums []*MetaAlbum
	if err := b.gormDb.
		Where("simplified_name IN ?", albumSimps).
		Find(&existingMetaAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddMetaAlbums(existingMetaAlbums)

	// var newMetaAlbums []*MetaAlbum
	var newAlbums []*MetaAlbum
	for _, sAlbum := range sAlbums {
		if _, ok := bi.GetMetaAlbum(sAlbum); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sAlbum.Artists)
		if !ok {
			log.Printf("WTF sAlbum: %v", sAlbum)
			return nil, fmt.Errorf("could not get artists for %s, but it should be there", sAlbum.Name)
		}
		newAlbums = append(newAlbums, newMetaAlbum(sAlbum, bArtists))
	}
	if len(newAlbums) == 0 {
		return existingMetaAlbums, nil
	}
	if err := b.gormDb.Create(newAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddMetaAlbums(newAlbums)
	return append(existingMetaAlbums, newAlbums...), nil
}