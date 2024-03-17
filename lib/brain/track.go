package brain

import (
	"slices"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Track struct {
	gorm.Model
	SpotifyId string
	Name      string
	AlbumId   uint
	Album     *Album
	Artists   []*Artist `gorm:"many2many:track_artists;"`
}

func newTrack(sTrack *spotify.SimpleTrack, bAlbum *Album, bArtists []*Artist) *Track {
	return &Track{Name: sTrack.Name, SpotifyId: sTrack.ID.String(), AlbumId: bAlbum.ID, Album: bAlbum, Artists: bArtists}
}

// ToTracks returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) ToTracks(sTracks []*spotify.SimpleTrack) ([]*Track, error) {
	sTracks = sliceutils.Uniqe(sTracks, func(item *spotify.SimpleTrack) string { return item.ID.String() })

	sArtists := sliceutils.FlatMap(sTracks, func(item *spotify.SimpleTrack) []*spotify.SimpleArtist { return sliceutils.MapP(item.Artists) })
	allArtists, err := b.ToArtists(sArtists)
	if err != nil {
		return nil, err
	}
	// map of spotifyId to Artist
	bArtistMap := sliceutils.ToMap(allArtists, func(item *Artist) string { return item.SpotifyId })

	sAlbums := sliceutils.Map(sTracks, func(item *spotify.SimpleTrack) *spotify.SimpleAlbum { return &item.Album })
	allAlbums, err := b.ToAlbums(sAlbums)
	if err != nil {
		return nil, err
	}
	// map of spotifyId to Album
	bAlbumMap := sliceutils.ToMap(allAlbums, func(item *Album) string { return item.SpotifyId })

	spotifyIds := sliceutils.Map(sTracks, func(item *spotify.SimpleTrack) string { return item.ID.String() })
	var existingTracks []*Track
	if err := b.gormDb.
		Preload("Artists").
		Where("spotify_id IN ?", spotifyIds).
		Find(&existingTracks).Error; err != nil {
		return nil, err
	}

	var newTracks []*Track
	for _, sTrack := range sTracks {
		if slices.ContainsFunc(existingTracks, func(item *Track) bool { return item.SpotifyId == sTrack.ID.String() }) {
			continue
		}
		bArtists := sliceutils.Map(sTrack.Artists, func(item spotify.SimpleArtist) *Artist { return bArtistMap[item.ID.String()] })
		bAlbum := bAlbumMap[sTrack.Album.ID.String()]
		newTracks = append(newTracks, newTrack(sTrack, bAlbum, bArtists))
	}
	// All artists are already created, can exit
	if len(newTracks) == 0 {
		return existingTracks, nil
	}
	if err := b.gormDb.Create(newTracks).Error; err != nil {
		return nil, err
	}
	return append(existingTracks, newTracks...), nil
}
