package brain

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type SpotifyTrack struct {
	ID             uint `gorm:"primarykey"`
	SpotifyId      spotify.ID
	Name           string
	SpotifyAlbumId uint
	SpotifyAlbum   *SpotifyAlbum
	Artists        []*Artist `gorm:"many2many:track_artists;"`
	MetaTrackId    uint
	MetaTrack      *MetaTrack
}

func newSpotifyTrack(sTrack spotify.SimpleTrack, bSpotifyAlbum *SpotifyAlbum, bArtists []*Artist, bMetaTrack *MetaTrack) *SpotifyTrack {
	return &SpotifyTrack{
		Name:           sTrack.Name,
		SpotifyId:      sTrack.ID,
		SpotifyAlbumId: bSpotifyAlbum.ID,
		SpotifyAlbum:   bSpotifyAlbum,
		Artists:        bArtists,
		MetaTrackId:    bMetaTrack.ID,
		MetaTrack:      bMetaTrack,
	}
}

func upsertSpotifyTracks(b *Brain, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*SpotifyTrack, error) {
	// sqlxSpotifyTracks, err := upsertSpotifyTracksSqlx(b.sqlxDb, sTracks, bi.Clone())
	// if err != nil {
	// 	return nil, err
	// }
	gormSpotifyTracks, err := upsertSpotifyTracksGorm(b.gormDb, sTracks, bi)
	if err != nil {
		return nil, err
	}
	// // TODO(m-nny): check sqlxSpotifyTracks == gormSpotifyTracks
	// if len(gormSpotifyTracks) != len(sqlxSpotifyTracks) {
	// 	return nil, fmt.Errorf("len(gormSpotifyTracks) != len(sqlxSpotifyTracks): %d != %d", len(gormSpotifyTracks), len(sqlxSpotifyTracks))
	// }
	return gormSpotifyTracks, nil

}

func upsertSpotifyTracksGorm(db *gorm.DB, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*SpotifyTrack, error) {
	var existingTracks []*SpotifyTrack
	if err := db.
		Preload("Artists").
		Where("spotify_id IN ?", sliceutils.Map(sTracks, func(item spotify.SimpleTrack) spotify.ID { return item.ID })).
		Find(&existingTracks).Error; err != nil {
		return nil, err
	}
	bi.AddSpotifyTracks(existingTracks)

	var newTracks []*SpotifyTrack
	for _, sTrack := range sTracks {
		if _, ok := bi.GetSpotifyTrack(sTrack); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sTrack.Artists)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, fmt.Errorf("could not find artist for %s, but it should be there", sTrack.Name)
		}
		bSpotifyAlbum, ok := bi.GetSpotifyAlbum(sTrack.Album)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, fmt.Errorf("could not find spotify album for %s, but it should be there", sTrack.Name)
		}
		bMetaTrack, ok := bi.GetMetaTrack(sTrack)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, fmt.Errorf("could not find meta track for %s, but it should be there", sTrack.Name)
		}
		newTracks = append(newTracks, newSpotifyTrack(sTrack, bSpotifyAlbum, bArtists, bMetaTrack))
	}
	if len(newTracks) == 0 {
		return existingTracks, nil
	}
	if err := db.Create(newTracks).Error; err != nil {
		return nil, err
	}
	bi.AddSpotifyTracks(newTracks)
	return append(existingTracks, newTracks...), nil
}
