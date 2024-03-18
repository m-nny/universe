package brain

import (
	"log"
	"slices"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

func (b *Brain) batchSaveAlbumTracks(sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*Album, []*Track, error) {
	sAlbums = sliceutils.Unique(sAlbums, func(item spotify.SimpleAlbum) spotify.ID { return item.ID })
	sTracks = sliceutils.Unique(sTracks, func(item spotify.SimpleTrack) spotify.ID { return item.ID })

	var sArtists []spotify.SimpleArtist
	for _, sAlbum := range sAlbums {
		sArtists = append(sArtists, sAlbum.Artists...)
	}
	for _, sTrack := range sTracks {
		sArtists = append(sArtists, sTrack.Artists...)
	}

	bArtistMap, err := b.toArtistsMap(sArtists)
	if err != nil {
		return nil, nil, err
	}

	var existingAlbums []*Album
	if err := b.gormDb.
		Preload("Artists").
		Where("spotify_id IN ?", sliceutils.Map(sAlbums, func(item spotify.SimpleAlbum) spotify.ID { return item.ID })).
		Find(&existingAlbums).Error; err != nil {
		return nil, nil, err
	}

	var newAlbums []*Album
	for _, sAlbum := range sAlbums {
		if slices.ContainsFunc(existingAlbums, func(item *Album) bool { return item.SpotifyId == sAlbum.ID }) {
			continue
		}
		bArtists := sliceutils.Map(sAlbum.Artists, func(item spotify.SimpleArtist) *Artist { return bArtistMap[item.ID] })
		newAlbums = append(newAlbums, newAlbum(sAlbum, bArtists))
	}
	if len(newAlbums) > 0 {
		if err := b.gormDb.Create(newAlbums).Error; err != nil {
			return nil, nil, err
		}
	}
	allAlbums := append(existingAlbums, newAlbums...)
	bAlbumMap := sliceutils.ToMap(allAlbums, func(item *Album) spotify.ID { return item.SpotifyId })

	var existingTracks []*Track
	if err := b.gormDb.
		Preload("Artists").
		Where("spotify_id IN ?", sliceutils.Map(sTracks, func(item spotify.SimpleTrack) spotify.ID { return item.ID })).
		Find(&existingTracks).Error; err != nil {
		return nil, nil, err
	}

	var newTracks []*Track
	for _, sTrack := range sTracks {
		if slices.ContainsFunc(existingTracks, func(item *Track) bool { return item.SpotifyId == sTrack.ID }) {
			continue
		}
		bArtists := sliceutils.Map(sTrack.Artists, func(item spotify.SimpleArtist) *Artist { return bArtistMap[item.ID] })
		bAlbum := bAlbumMap[sTrack.Album.ID]
		if bAlbum == nil {
			log.Printf("WTF sTrack: %v", sTrack)
		}
		newTracks = append(newTracks, newTrack(sTrack, bAlbum, bArtists))
	}
	if len(newTracks) > 0 {
		if err := b.gormDb.Create(newTracks).Error; err != nil {
			return nil, nil, err
		}
	}
	allTracks := append(existingTracks, newTracks...)
	return allAlbums, allTracks, nil
}
