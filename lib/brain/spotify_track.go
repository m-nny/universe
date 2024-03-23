package brain

import (
	"github.com/zmb3/spotify/v2"
)

type SpotifyTrack struct {
	ID             uint `gorm:"primarykey"`
	SpotifyId      spotify.ID
	Name           string
	SpotifyAlbumId uint
	SpotifyAlbum   *SpotifyAlbum
	Artists        []*Artist `gorm:"many2many:track_artists;"`
	MetaTrackId    uint
	MetaTrack      *MetaTrack
}

func newSpotifyTrack(sTrack spotify.SimpleTrack, bSpotifyAlbum *SpotifyAlbum, bArtists []*Artist, bMetaTrack *MetaTrack) *SpotifyTrack {
	return &SpotifyTrack{
		Name:           sTrack.Name,
		SpotifyId:      sTrack.ID,
		SpotifyAlbumId: bSpotifyAlbum.ID,
		SpotifyAlbum:   bSpotifyAlbum,
		Artists:        bArtists,
		MetaTrackId:    bMetaTrack.ID,
		MetaTrack:      bMetaTrack,
	}
}

// SaveTracks returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveTracks(savedTracks []spotify.SavedTrack) ([]*SpotifyTrack, error) {
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
