package brain

import (
	"fmt"

	"github.com/zmb3/spotify/v2"
)

type SpotifyAlbum struct {
	ID          uint `gorm:"primarykey"`
	SpotifyId   spotify.ID
	Name        string
	Artists     []*Artist `gorm:"many2many:spotify_album_artists;"`
	MetaAlbumId uint
	MetaAlbum   *MetaAlbum
}

func (s *SpotifyAlbum) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", s.Name)
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
	var albumSIds []spotify.ID
	for _, sAlbum := range sAlbums {
		albumSIds = append(albumSIds, sAlbum.ID)
	}
	var existingSpotifyAlbums []*SpotifyAlbum
	if err := b.gormDb.
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
	if err := b.gormDb.Create(newAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddSpotifyAlbums(newAlbums)

	return append(existingSpotifyAlbums, newAlbums...), nil
}
