package brain

import (
	"fmt"
	"log"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
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
		bArtists, ok := bi.GetArtists(sAlbum.Artists)
		if !ok {
			log.Printf("WTF sAlbum: %v", sAlbum)
			return nil, fmt.Errorf("could not get artists for %s, but it should be there", sAlbum.Name)
		}
		newAlbums = append(newAlbums, newMetaAlbum(sAlbum, bArtists))
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

func upsertSpotifyTracks(b *Brain, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*SpotifyTrack, error) {
	var existingTracks []*SpotifyTrack
	if err := b.gormDb.
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
	if len(newTracks) > 0 {
		if err := b.gormDb.Create(newTracks).Error; err != nil {
			return nil, err
		}
	}
	bi.AddSpotifyTracks(newTracks)
	return append(existingTracks, newTracks...), nil
}

// batchSaveAlbumTracks returns Brain representain of a spotify albums and tracks
//   - It will create new entries in DB if necessary for albums, tracks, artists
//   - It will deduplicate returned albums base on spotify.ID, this may result in len(result) < len(sAlbums)
//   - NOTE: it does not store all spotify.IDs of duplicated at the moment
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) batchSaveAlbumTracks(sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*SpotifyAlbum, []*SpotifyTrack, error) {
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

	if _, err := upsertMetaTracks(b, sTracks, bi); err != nil {
		return nil, nil, err
	}
	allTracks, err := upsertSpotifyTracks(b, sTracks, bi)
	if err != nil {
		return nil, nil, err
	}
	return allAlbums, allTracks, nil
}
