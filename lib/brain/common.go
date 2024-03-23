package brain

import (
	"fmt"
	"log"
	"slices"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/spotify/utils"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

func upsertArtists(b *Brain, sArtists []spotify.SimpleArtist, bi *brainIndex) ([]*Artist, error) {
	sArtists = sliceutils.Unique(sArtists, func(item spotify.SimpleArtist) spotify.ID { return item.ID })
	var spotifyIds []spotify.ID
	for _, sAlbum := range sArtists {
		spotifyIds = append(spotifyIds, sAlbum.ID)
	}

	var existingArtists []*Artist
	if err := b.gormDb.
		Where("spotify_id IN ?", spotifyIds).
		Find(&existingArtists).Error; err != nil {
		return nil, err
	}
	bi.AddArtists(existingArtists)

	var newArtists []*Artist
	for _, sArtist := range sArtists {
		if _, ok := bi.GetArtist(sArtist); ok {
			continue
		}
		newArtists = append(newArtists, newArtist(sArtist))
	}
	// All artists are already created, can exit
	if len(newArtists) == 0 {
		return existingArtists, nil
	}
	if err := b.gormDb.Create(newArtists).Error; err != nil {
		return nil, err
	}
	bi.AddArtists(newArtists)
	return append(existingArtists, newArtists...), nil
}

func upsertMetaAlbums(b *Brain, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*MetaAlbum, error) {
	sAlbums = sliceutils.Unique(sAlbums, utils.SimplifiedAlbumName)
	var albumSimps []string
	for _, sAlbum := range sAlbums {
		albumSimps = append(albumSimps, utils.SimplifiedAlbumName(sAlbum))
	}

	var existingMetaAlbums []*MetaAlbum
	if err := b.gormDb.
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
		newAlbums = append(newAlbums, newMetaAlbum(sAlbum))
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

func upsertTracks(b *Brain, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*Track, error) {
	var existingTracks []*Track
	if err := b.gormDb.
		Preload("Artists").
		Where("spotify_id IN ?", sliceutils.Map(sTracks, func(item spotify.SimpleTrack) spotify.ID { return item.ID })).
		Find(&existingTracks).Error; err != nil {
		return nil, err
	}

	var newTracks []*Track
	for _, sTrack := range sTracks {
		if slices.ContainsFunc(existingTracks, func(item *Track) bool { return item.SpotifyId == sTrack.ID }) {
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
			return nil, fmt.Errorf("could not find album for %s, but it should be there", sTrack.Name)
		}
		newTracks = append(newTracks, newTrack(sTrack, bSpotifyAlbum, bArtists))
	}
	if len(newTracks) > 0 {
		if err := b.gormDb.Create(newTracks).Error; err != nil {
			return nil, err
		}
	}
	allTracks := append(existingTracks, newTracks...)
	return allTracks, nil
}

// batchSaveAlbumTracks returns Brain representain of a spotify albums and tracks
//   - It will create new entries in DB if necessary for albums, tracks, artists
//   - It will deduplicate returned albums base on spotify.ID, this may result in len(result) < len(sAlbums)
//   - NOTE: it does not store all spotify.IDs of duplicated at the moment
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) batchSaveAlbumTracks(sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*SpotifyAlbum, []*Track, error) {
	sAlbums = sliceutils.Unique(sAlbums, func(item spotify.SimpleAlbum) spotify.ID { return item.ID })
	sTracks = sliceutils.Unique(sTracks, func(item spotify.SimpleTrack) spotify.ID { return item.ID })

	var sArtists []spotify.SimpleArtist
	for _, sAlbum := range sAlbums {
		sArtists = append(sArtists, sAlbum.Artists...)
	}
	for _, sTrack := range sTracks {
		sArtists = append(sArtists, sTrack.Artists...)
	}

	bi := newBrainIndex()
	_, err := upsertArtists(b, sArtists, bi)
	if err != nil {
		return nil, nil, err
	}

	if _, err := upsertMetaAlbums(b, sAlbums, bi); err != nil {
		return nil, nil, err
	}
	allAlbums, err := upsertSpotifyAlbums(b, sAlbums, bi)
	if err != nil {
		return nil, nil, err
	}

	allTracks, err := upsertTracks(b, sTracks, bi)
	if err != nil {
		return nil, nil, err
	}
	return allAlbums, allTracks, nil
}
