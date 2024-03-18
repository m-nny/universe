package brain

import (
	"fmt"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

type AlbumId uint

type Album struct {
	gorm.Model
	SpotifyId spotify.ID
	Name      string
	Artists   []*Artist `gorm:"many2many:album_artists;"`
}

func (s *Album) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", s.Name)
}

func newAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist) *Album {
	return &Album{Name: sAlbum.Name, SpotifyId: sAlbum.ID, Artists: bArtists}
}

// SaveAlbums returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveAlbums(fullAlbums []*spotify.FullAlbum) ([]*Album, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sAlbum := range fullAlbums {
		// sAlbum.SimpleAlbum.Artists = sAlbum.Artists
		sAlbums = append(sAlbums, sAlbum.SimpleAlbum)
		sTracks = append(sTracks, sAlbum.Tracks.Tracks...)
	}
	albums, _, err := b.batchSave(sAlbums, sTracks)
	return albums, err
}
