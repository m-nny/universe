package brain

import (
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Artist struct {
	ID        uint `gorm:"primarykey"`
	SpotifyId spotify.ID
	Name      string
	Albums    []*SpotifyAlbum `gorm:"many2many:spotify_album_artists;"`
}

func newArtist(sArtist spotify.SimpleArtist) *Artist {
	return &Artist{Name: sArtist.Name, SpotifyId: sArtist.ID}
}

// SaveArtists takes list of spotify Artists and returns Brain representain of them.
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned artists, this may result in len(result) < len(sArtists)
func (b *Brain) SaveArtists(sArtists []*spotify.FullArtist) ([]*Artist, error) {
	bi := newBrainIndex()
	sSimpleArtists := sliceutils.Map(sArtists, func(item *spotify.FullArtist) spotify.SimpleArtist { return item.SimpleArtist })
	return upsertArtists(b, sSimpleArtists, bi)
}

func (b *Brain) _saveArtists(sArtists []spotify.SimpleArtist) ([]*Artist, error) {
	return upsertArtists(b, sArtists, newBrainIndex())
}
