package brain

import (
	"slices"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Album struct {
	gorm.Model
	SpotifyId string
	Name      string
	Artists   []*Artist `gorm:"many2many:album_artists;"`
}

func newAlbum(sAlbum *spotify.SimpleAlbum, bArtists []*Artist) *Album {
	return &Album{Name: sAlbum.Name, SpotifyId: sAlbum.ID.String(), Artists: bArtists}
}

// toAlbum returns Brain representain of a spotify album
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) toAlbum(sAlbum *spotify.SimpleAlbum, bArtists []*Artist) (*Album, error) {
	var album Album
	if err := b.gormDb.
		Preload("Artists").
		Where(&Album{SpotifyId: sAlbum.ID.String()}).
		Attrs(newAlbum(sAlbum, bArtists)).
		FirstOrCreate(&album).Error; err != nil {
		return nil, err
	}
	return &album, nil
}

// ToAlbums returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) ToAlbums(sAlbums []*spotify.SimpleAlbum) ([]*Album, error) {
	sAlbums = sliceutils.Uniqe(sAlbums, func(item *spotify.SimpleAlbum) string { return item.ID.String() })
	sArtists := sliceutils.FlatMap(sAlbums, func(item *spotify.SimpleAlbum) []*spotify.SimpleArtist { return sliceutils.MapP(item.Artists) })
	allArtists, err := b.ToArtists(sArtists)
	if err != nil {
		return nil, err
	}
	// map of spotifyId to Artist
	bArtistMap := sliceutils.ToMap(allArtists, func(item *Artist) string { return item.SpotifyId })

	spotifyIds := sliceutils.Map(sAlbums, func(item *spotify.SimpleAlbum) string { return item.ID.String() })
	var existingAlbums []*Album
	if err := b.gormDb.
		Where("spotify_id IN ?", spotifyIds).
		Find(&existingAlbums).Error; err != nil {
		return nil, err
	}

	var newAlbums []*Album
	for _, sAlbum := range sAlbums {
		if slices.ContainsFunc(existingAlbums, func(item *Album) bool { return item.SpotifyId == sAlbum.ID.String() }) {
			continue
		}
		bArtists := sliceutils.Map(sAlbum.Artists, func(item spotify.SimpleArtist) *Artist { return bArtistMap[item.ID.String()] })
		newAlbums = append(newAlbums, newAlbum(sAlbum, bArtists))
	}
	// All artists are already created, can exit
	if len(newAlbums) == 0 {
		return existingAlbums, nil
	}
	if err := b.gormDb.Create(newAlbums).Error; err != nil {
		return nil, err
	}
	return append(existingAlbums, newAlbums...), nil
}
