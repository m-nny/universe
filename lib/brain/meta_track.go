package brain

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
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

func upsertMetaTracks(b *Brain, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*MetaTrack, error) {
	sTracks = sliceutils.Unique(sTracks, bi.MustTrackSimplifiedName)
	var trackSimps []string
	for _, sTrack := range sTracks {
		simpName, ok := bi.TrackSimplifiedName(sTrack)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, fmt.Errorf("could not get simplified name for %s, but it should be there", sTrack.Name)
		}
		trackSimps = append(trackSimps, simpName)
	}

	var existingMetaTracks []*MetaTrack
	if err := b.gormDb.
		Where("simplified_name IN ?", trackSimps).
		Find(&existingMetaTracks).Error; err != nil {
		return nil, err
	}
	bi.AddMetaTracks(existingMetaTracks)

	// var newMetaTracks []*MetaTrack
	var newTracks []*MetaTrack
	for _, sTrack := range sTracks {
		if _, ok := bi.GetMetaTrack(sTrack); ok {
			continue
		}
		bMetaAlbum, ok := bi.GetMetaAlbum(sTrack.Album)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, fmt.Errorf("could not find meta album for %s, but it should be there", sTrack.Name)
		}
		newTracks = append(newTracks, newMetaTrack(sTrack, bMetaAlbum))
	}
	if len(newTracks) == 0 {
		return existingMetaTracks, nil
	}
	if err := b.gormDb.Create(newTracks).Error; err != nil {
		return nil, err
	}
	bi.AddMetaTracks(newTracks)
	return append(existingMetaTracks, newTracks...), nil
}