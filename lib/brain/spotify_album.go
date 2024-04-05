package brain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type SpotifyAlbum struct {
	ID          uint       `gorm:"primarykey"`
	SpotifyId   spotify.ID `db:"spotify_id"`
	Name        string
	Artists     []*Artist `gorm:"many2many:spotify_album_artists;"`
	MetaAlbumId uint      `db:"meta_album_id"`
	MetaAlbum   *MetaAlbum
}

func (s *SpotifyAlbum) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.Name
}

func newSpotifyAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist, bMetaAlbum *MetaAlbum) *SpotifyAlbum {
	return &SpotifyAlbum{
		Name:        sAlbum.Name,
		SpotifyId:   sAlbum.ID,
		Artists:     bArtists,
		MetaAlbumId: bMetaAlbum.ID,
		MetaAlbum:   bMetaAlbum,
	}
}

func upsertSpotifyAlbums(b *Brain, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*SpotifyAlbum, error) {
	// sqlxBi := bi.Clone()
	// sqlxSpotifyAlbums, err := upsertSpotifyAlbumsSqlx(b.sqlxDb, sAlbums, sqlxBi)
	// if err != nil {
	// 	return nil, err
	// }
	gormSpotifyAlbums, err := upsertSpotifyAlbumsGorm(b.gormDb, sAlbums, bi)
	if err != nil {
		return nil, err
	}
	// // TODO(m-nny): check sqlxSpotifyAlbums == gormSpotifyAlbums
	// if len(gormSpotifyAlbums) != len(sqlxSpotifyAlbums) {
	// 	return nil, fmt.Errorf("len(gormSpotifyAlbums) != len(sqlxSpotifyAlbums): %d != %d", len(gormSpotifyAlbums), len(sqlxSpotifyAlbums))
	// }
	return gormSpotifyAlbums, nil
}

func upsertSpotifyAlbumsGorm(db *gorm.DB, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*SpotifyAlbum, error) {
	var albumSIds []spotify.ID
	for _, sAlbum := range sAlbums {
		albumSIds = append(albumSIds, sAlbum.ID)
	}
	var existingSpotifyAlbums []*SpotifyAlbum
	if err := db.
		Preload("Artists").
		Where("spotify_id IN ?", albumSIds).
		Find(&existingSpotifyAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddSpotifyAlbums(existingSpotifyAlbums)

	var newAlbums []*SpotifyAlbum
	for _, sAlbum := range sAlbums {
		if _, ok := bi.GetSpotifyAlbum(sAlbum); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sAlbum.Artists)
		if !ok {
			return nil, fmt.Errorf("bArtist not found")
		}
		bMetaAlbum, ok := bi.GetMetaAlbum(sAlbum)
		if !ok {
			return nil, fmt.Errorf("bMetaAlbum not found")
		}
		newAlbums = append(newAlbums, newSpotifyAlbum(sAlbum, bArtists, bMetaAlbum))
	}
	if len(newAlbums) == 0 {
		return existingSpotifyAlbums, nil
	}
	if err := db.Create(newAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddSpotifyAlbums(newAlbums)

	return append(existingSpotifyAlbums, newAlbums...), nil
}

func upsertSpotifyAlbumsSqlx(db *sqlx.DB, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*SpotifyAlbum, error) {
	return nil, fmt.Errorf("not implmemented")
}
