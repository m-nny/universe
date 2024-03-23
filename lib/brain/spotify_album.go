package brain

import (
	"fmt"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/spotify/utils"
)

type SpotifyAlbum struct {
	gorm.Model
	SpotifyId      spotify.ID
	Name           string
	Artists        []*Artist `gorm:"many2many:spotify_album_artists;"`
	SimplifiedName string
	MetaAlbumId    uint
	MetaAlbum      *MetaAlbum
}

func (s *SpotifyAlbum) String() string {
	if s == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s", s.Name)
}

func newSpotifyAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist, bMetaAlbum *MetaAlbum) *SpotifyAlbum {
	return &SpotifyAlbum{
		Name:           sAlbum.Name,
		SpotifyId:      sAlbum.ID,
		Artists:        bArtists,
		SimplifiedName: utils.SimplifiedAlbumName(sAlbum),
		MetaAlbumId:    bMetaAlbum.ID,
		MetaAlbum:      bMetaAlbum,
	}
}

// SaveAlbums returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveAlbums(fullAlbums []*spotify.FullAlbum) ([]*SpotifyAlbum, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sAlbum := range fullAlbums {
		// sAlbum.SimpleAlbum.Artists = sAlbum.Artists
		sAlbums = append(sAlbums, sAlbum.SimpleAlbum)
		sTracks = append(sTracks, sAlbum.Tracks.Tracks...)
	}
	albums, _, err := b.batchSaveAlbumTracks(sAlbums, sTracks)
	return albums, err
}
