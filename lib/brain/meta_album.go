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
	AnyName        string
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
		AnyName:        sAlbum.Name,
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
		Preload("Artists").
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

// SaveAlbums returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveAlbums(fullAlbums []*spotify.FullAlbum) ([]*MetaAlbum, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sAlbum := range fullAlbums {
		// sAlbum.SimpleAlbum.Artists = sAlbum.Artists
		sAlbums = append(sAlbums, sAlbum.SimpleAlbum)
		sTracks = append(sTracks, sAlbum.Tracks.Tracks...)
	}
	albums, _, err := b.batchSaveAlbumTracks(sAlbums, sTracks)
	return albums, err
}

func (b *Brain) SaveSimpleAlbums(sAlbums []spotify.SimpleAlbum) ([]*MetaAlbum, error) {
	albums, _, err := b.batchSaveAlbumTracks(sAlbums, nil)
	return albums, err
}
