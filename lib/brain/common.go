package brain

import (
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

// batchSaveAlbumTracksGorm returns Brain representain of a spotify albums and tracks
//   - It will create new entries in DB if necessary for albums, tracks, artists
//   - It will deduplicate returned albums base on spotify.ID, this may result in len(bMetaAlbums) < len(sAlbums)
//   - It will deduplicate returned tracks base on spotify.ID, this may result in len(bMetaTracks) < len(sTracks)
func batchSaveAlbumTracksGorm(db *gorm.DB, sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*MetaAlbum, []*MetaTrack, error) {
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

	_, err := upsertArtistsGorm(db, sArtists, bi)
	if err != nil {
		return nil, nil, err
	}

	metaAlbums, err := upsertMetaAlbumsGorm(db, sAlbums, bi)
	if err != nil {
		return nil, nil, err
	}
	if _, err := upsertSpotifyAlbumsGorm(db, sAlbums, bi); err != nil {
		return nil, nil, err
	}

	metaTracks, err := upsertMetaTracksGorm(db, sTracks, bi)
	if err != nil {
		return nil, nil, err
	}
	if _, err := upsertSpotifyTracksGorm(db, sTracks, bi); err != nil {
		return nil, nil, err
	}
	return metaAlbums, metaTracks, nil
}

// batchSaveAlbumTracksSqlx returns Brain representain of a spotify albums and tracks
//   - It will create new entries in DB if necessary for albums, tracks, artists
//   - It will deduplicate returned albums base on spotify.ID, this may result in len(bMetaAlbums) < len(sAlbums)
//   - It will deduplicate returned tracks base on spotify.ID, this may result in len(bMetaTracks) < len(sTracks)
func batchSaveAlbumTracksSqlx(db *sqlx.DB, sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*MetaAlbum, []*MetaTrack, error) {
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

	_, err := upsertArtistsSqlx(db, sArtists, bi)
	if err != nil {
		return nil, nil, err
	}

	metaAlbums, err := upsertMetaAlbumsSqlx(db, sAlbums, bi)
	if err != nil {
		return nil, nil, err
	}
	if _, err := upsertSpotifyAlbumsSqlx(db, sAlbums, bi); err != nil {
		return nil, nil, err
	}

	metaTracks, err := upsertMetaTracksSqlx(db, sTracks, bi)
	if err != nil {
		return nil, nil, err
	}
	if _, err := upsertSpotifyTracksSqlx(db, sTracks, bi); err != nil {
		return nil, nil, err
	}
	return metaAlbums, metaTracks, nil
}
