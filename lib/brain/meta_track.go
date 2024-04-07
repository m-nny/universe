package brain

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

type MetaTrack struct {
	ID             uint `gorm:"primarykey"`
	SimplifiedName string
	MetaAlbumID    uint
	MetaAlbum      *MetaAlbum
	Artists        []*Artist `gorm:"many2many:meta_track_artists;"`
}

func (s *MetaTrack) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.SimplifiedName
}

func newMetaTrack(sTrack spotify.SimpleTrack, bMetaAlbum *MetaAlbum, bArtists []*Artist) *MetaTrack {
	return &MetaTrack{
		SimplifiedName: utils.SimplifiedTrackName(sTrack, bMetaAlbum.SimplifiedName),
		MetaAlbumID:    bMetaAlbum.ID,
		MetaAlbum:      bMetaAlbum,
		Artists:        bArtists,
	}
}

func (b *Brain) MetaTrackCount() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaTrack{}).Count(&count).Error
	return int(count), err
}

func upsertMetaTracks(b *Brain, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*MetaTrack, error) {
	// sqlxBi := bi.Clone()
	// sqlxMetaTracks, err := upsertMetaTracksSqlx(b.sqlxDb, sTracks, sqlxBi)
	// if err != nil {
	// 	return nil, err
	// }
	gormMetaTracks, err := upsertMetaTracksGorm(b.gormDb, sTracks, bi)
	if err != nil {
		return nil, err
	}
	// // TODO(m-nny): check sqlxMetaTracks == gormMetaTracks
	// if len(gormMetaTracks) != len(sqlxMetaTracks) {
	// 	return nil, fmt.Errorf("len(gormMetaTracks) != len(sqlxMetaTracks): %d != %d", len(gormMetaTracks), len(sqlxMetaTracks))
	// }
	return gormMetaTracks, nil

}

func upsertMetaTracksGorm(db *gorm.DB, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*MetaTrack, error) {
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
	if err := db.
		Preload("Artists").
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
		bArtists, ok := bi.GetArtists(sTrack.Artists)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, fmt.Errorf("could not find meta album for %s, but it should be there", sTrack.Name)
		}
		newTracks = append(newTracks, newMetaTrack(sTrack, bMetaAlbum, bArtists))
	}
	if len(newTracks) == 0 {
		return existingMetaTracks, nil
	}
	if err := db.Create(newTracks).Error; err != nil {
		return nil, err
	}
	bi.AddMetaTracks(newTracks)
	return append(existingMetaTracks, newTracks...), nil
}

// SaveTracks returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveTracks(savedTracks []spotify.SavedTrack, username string) ([]*MetaTrack, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sFullTrack := range savedTracks {
		sAlbum := sFullTrack.Album
		if sAlbum.ID == "" {
			sAlbum = sFullTrack.SimpleTrack.Album
		}

		// we are using sTrack.Album to associate it with bAlbum later
		sTrack := sFullTrack.SimpleTrack
		sTrack.Album = sAlbum

		sAlbums = append(sAlbums, sAlbum)
		sTracks = append(sTracks, sTrack)
	}
	_, tracks, err := b.batchSaveAlbumTracks(sAlbums, sTracks)
	if err := b.addSavedTracks(username, tracks); err != nil {
		return nil, err
	}

	return tracks, err
}
