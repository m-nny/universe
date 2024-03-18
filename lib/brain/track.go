package brain

import (
	"log"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	lspotify "github.com/m-nny/universe/lib/spotify"
)

type Track struct {
	gorm.Model
	SpotifyId spotify.ID
	Name      string
	AlbumId   uint
	Album     *Album
	Artists   []*Artist `gorm:"many2many:track_artists;"`
}

func newTrack(sTrack spotify.SimpleTrack, bAlbum *Album, bArtists []*Artist) *Track {
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
		log.Printf("sTrack: %s", sTrack.Name)

		log.Printf("  sTrack.Album: %+v", sTrack.Album.Name)
		log.Printf("  sTrack.SimpleTrack.Album: %+v", sTrack.SimpleTrack.Album.Name)

		log.Printf("  sTrack.Artists: %+v", lspotify.SArtistsString(sTrack.Artists))
		log.Printf("  sTrack.SimpleTrack.Artists: %+v", lspotify.SArtistsString(sTrack.SimpleTrack.Artists))
		sAlbums = append(sAlbums, sTrack.Album)

		// we are using sTrack.Album to associate it with bAlbum later
		sTrack.SimpleTrack.Album = sTrack.Album
		sTracks = append(sTracks, sTrack.SimpleTrack)
	}
	_, tracks, err := b.batchSave(sAlbums, sTracks)
	return tracks, err
}
