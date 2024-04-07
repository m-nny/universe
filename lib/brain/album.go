package brain

import (
	"github.com/zmb3/spotify/v2"
)

// SaveAlbumsGorm returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveAlbumsGorm(fullAlbums []*spotify.FullAlbum) ([]*MetaAlbum, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sAlbum := range fullAlbums {
		// sAlbum.SimpleAlbum.Artists = sAlbum.Artists
		sAlbums = append(sAlbums, sAlbum.SimpleAlbum)
		sTracks = append(sTracks, sAlbum.Tracks.Tracks...)
	}
	albums, _, err := batchSaveAlbumTracksGorm(b.gormDb, sAlbums, sTracks)
	return albums, err
}

// SaveAlbumsSqlx returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveAlbumsSqlx(fullAlbums []*spotify.FullAlbum) ([]*MetaAlbum, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sAlbum := range fullAlbums {
		// sAlbum.SimpleAlbum.Artists = sAlbum.Artists
		sAlbums = append(sAlbums, sAlbum.SimpleAlbum)
		sTracks = append(sTracks, sAlbum.Tracks.Tracks...)
	}
	albums, _, err := batchSaveAlbumTracksSqlx(b.sqlxDb, sAlbums, sTracks)
	return albums, err
}

func (b *Brain) SaveSimpleAlbumsGorm(sAlbums []spotify.SimpleAlbum) ([]*MetaAlbum, error) {
	albums, _, err := batchSaveAlbumTracksGorm(b.gormDb, sAlbums, nil)
	return albums, err
}

func (b *Brain) SaveSimpleAlbumsSqlx(sAlbums []spotify.SimpleAlbum) ([]*MetaAlbum, error) {
	albums, _, err := batchSaveAlbumTracksSqlx(b.sqlxDb, sAlbums, nil)
	return albums, err
}

func (b *Brain) MetaAlbumCount() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaAlbum{}).Count(&count).Error
	return int(count), err
}
