package brain

import (
	"github.com/zmb3/spotify/v2"
)

// SaveTracksGorm returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveTracksGorm(savedTracks []spotify.SavedTrack, username string) ([]*MetaTrack, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sFullTrack := range savedTracks {
		sAlbum := sFullTrack.Album
		if sAlbum.ID == "" {
			sAlbum = sFullTrack.SimpleTrack.Album
		}

		// we are using sTrack.Album to associate it with bAlbum later
		sTrack := sFullTrack.SimpleTrack
		sTrack.Album = sAlbum

		sAlbums = append(sAlbums, sAlbum)
		sTracks = append(sTracks, sTrack)
	}
	_, tracks, err := batchSaveAlbumTracksGorm(b.gormDb, sAlbums, sTracks)
	if err := addSavedTracksGorm(b.gormDb, username, tracks); err != nil {
		return nil, err
	}

	return tracks, err
}

// SaveTracksSqlx returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveTracksSqlx(savedTracks []spotify.SavedTrack, username string) ([]*MetaTrack, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sFullTrack := range savedTracks {
		sAlbum := sFullTrack.Album
		if sAlbum.ID == "" {
			sAlbum = sFullTrack.SimpleTrack.Album
		}

		// we are using sTrack.Album to associate it with bAlbum later
		sTrack := sFullTrack.SimpleTrack
		sTrack.Album = sAlbum

		sAlbums = append(sAlbums, sAlbum)
		sTracks = append(sTracks, sTrack)
	}
	_, tracks, err := batchSaveAlbumTracksSqlx(b.sqlxDb, sAlbums, sTracks)
	if err := addSavedTracksSqlx(b.sqlxDb, username, tracks); err != nil {
		return nil, err
	}

	return tracks, err
}

func (b *Brain) MetaTrackCountGorm() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaTrack{}).Count(&count).Error
	return int(count), err
}

func (b *Brain) MetaTrackCountSqlx() (int, error) {
	var cnt int
	if err := b.sqlxDb.Get(&cnt, `SELECT COUNT(*) FROM meta_tracks`); err != nil {
		return 0, err
	}
	return cnt, nil
}
