package brain

import (
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type Track struct {
	gorm.Model
	SpotifyId spotify.ID
	Name      string
	AlbumId   uint
	Album     *SpotifyAlbum
	Artists   []*Artist `gorm:"many2many:track_artists;"`
}

func newTrack(sTrack spotify.SimpleTrack, bAlbum *SpotifyAlbum, bArtists []*Artist) *Track {
	return &Track{Name: sTrack.Name, SpotifyId: sTrack.ID, AlbumId: bAlbum.ID, Album: bAlbum, Artists: bArtists}
}

// SaveTracks returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveTracks(savedTracks []spotify.SavedTrack) ([]*Track, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sTrack := range savedTracks {
		sAlbums = append(sAlbums, sTrack.Album)

		// we are using sTrack.Album to associate it with bAlbum later
		sTrack.SimpleTrack.Album = sTrack.Album
		sTracks = append(sTracks, sTrack.SimpleTrack)
	}
	_, tracks, err := b.batchSaveAlbumTracks(sAlbums, sTracks)
	return tracks, err
}
