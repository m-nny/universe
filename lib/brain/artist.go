package brain

import (
	"slices"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Artist struct {
	gorm.Model
	SpotifyId spotify.ID
	Name      string
	Albums    []*Album `gorm:"many2many:album_artists;"`
}

func newArtist(sArtist *spotify.SimpleArtist) *Artist {
	return &Artist{Name: sArtist.Name, SpotifyId: sArtist.ID}
}

func (b *Brain) ToArtist(sArtist *spotify.SimpleArtist) (*Artist, error) {
	var artist Artist
	if err := b.gormDb.
		Where(&Artist{SpotifyId: sArtist.ID}).
		Attrs(newArtist(sArtist)).
		FirstOrCreate(&artist).Error; err != nil {
		return nil, err
	}
	return &artist, nil
}

// SaveArtists takes list of spotify Artists and returns Brain representain of them.
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned artists, this may result in len(result) < len(sArtists)
func (b *Brain) SaveArtists(sArtists []*spotify.SimpleArtist) ([]*Artist, error) {
	sArtists = sliceutils.Unique(sArtists, func(item *spotify.SimpleArtist) spotify.ID { return item.ID })
	spotifyIds := sliceutils.Map(sArtists, func(item *spotify.SimpleArtist) spotify.ID { return item.ID })

	var existingArtists []*Artist
	if err := b.gormDb.
		Where("spotify_id IN ?", spotifyIds).
		Find(&existingArtists).Error; err != nil {
		return nil, err
	}

	var newArtists []*Artist
	for _, sArtist := range sArtists {
		if slices.ContainsFunc(existingArtists, func(item *Artist) bool { return item.SpotifyId == sArtist.ID }) {
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
	return append(existingArtists, newArtists...), nil
}

func (b *Brain) toArtistsMap(sArtists []*spotify.SimpleArtist) (map[spotify.ID]*Artist, error) {
	bArtists, err := b.SaveArtists(sArtists)
	if err != nil {
		return nil, err
	}
	return sliceutils.ToMap(bArtists, func(item *Artist) spotify.ID { return item.SpotifyId }), nil
}
