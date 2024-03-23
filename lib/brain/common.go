package brain

import (
	"fmt"
	"log"
	"slices"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/spotify/utils"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

// batchSaveAlbumTracks returns Brain representain of a spotify albums and tracks
//   - It will create new entries in DB if necessary for albums, tracks, artists
//   - It will deduplicate returned albums base on spotify.ID, this may result in len(result) < len(sAlbums)
//   - NOTE: it does not store all spotify.IDs of duplicated at the moment
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) batchSaveAlbumTracks(sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*SpotifyAlbum, []*Track, error) {
	sAlbums = sliceutils.Unique(sAlbums, utils.SimplifiedAlbumName)
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

	var albumSIds []spotify.ID
	var albumSimpNames []string
	for _, sAlbum := range sAlbums {
		albumSIds = append(albumSIds, sAlbum.ID)
		albumSimpNames = append(albumSimpNames, utils.SimplifiedAlbumName(sAlbum))
	}
	var existingAlbums []*SpotifyAlbum
	if err := b.gormDb.
		Preload("Artists").
		Where("spotify_id IN ?", albumSIds).
		Or("simplified_name in ?", albumSimpNames).
		Find(&existingAlbums).Error; err != nil {
		return nil, nil, err
	}
	ai := newAlbumIndex(existingAlbums)

	var newAlbums []*SpotifyAlbum
	for _, sAlbum := range sAlbums {
		if _, ok := ai.Get(sAlbum); ok {
			continue
		}
		bArtists := sliceutils.Map(sAlbum.Artists, func(item spotify.SimpleArtist) *Artist { return bArtistMap[item.ID] })
		newAlbums = append(newAlbums, newSpotifyAlbum(sAlbum, bArtists))
	}
	if len(newAlbums) > 0 {
		if err := b.gormDb.Create(newAlbums).Error; err != nil {
			return nil, nil, err
		}
		ai.Add(newAlbums)
	}
	allAlbums := append(existingAlbums, newAlbums...)

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
		bAlbum, ok := ai.Get(sTrack.Album)
		if !ok {
			log.Printf("WTF sTrack: %v", sTrack)
			return nil, nil, fmt.Errorf("could not find album for %s, but it should be there", sTrack.Name)
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
