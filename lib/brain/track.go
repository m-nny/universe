package brain

import (
	"slices"

	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Track struct {
	gorm.Model
	SpotifyId spotify.ID
	Name      string
	AlbumId   uint
	Album     *Album
	Artists   []*Artist `gorm:"many2many:track_artists;"`
}

func newTrack(sTrack spotify.SavedTrack, bAlbum *Album, bArtists []*Artist) *Track {
	return &Track{Name: sTrack.Name, SpotifyId: sTrack.ID, AlbumId: bAlbum.ID, Album: bAlbum, Artists: bArtists}
}

// SaveTracks returns Brain representain of a spotify tracks
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sTracks)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveTracks(sTracks []spotify.SavedTrack) ([]*Track, error) {
	sTracks = sliceutils.Unique(sTracks, func(item spotify.SavedTrack) spotify.ID { return item.ID })
	sAlbums := sliceutils.Map(sTracks, func(item spotify.SavedTrack) *spotify.SimpleAlbum { return &item.Album })

	sArtists := sliceutils.FlatMap(sTracks, func(item spotify.SavedTrack) []*spotify.SimpleArtist { return sliceutils.MapP(item.Artists) })
	sArtists = append(sArtists, sliceutils.FlatMap(sAlbums, func(item *spotify.SimpleAlbum) []*spotify.SimpleArtist { return sliceutils.MapP(item.Artists) })...)

	bArtistMap, err := b.toArtistsMap(sArtists)
	if err != nil {
		return nil, err
	}

	bAlbumMap, err := b.toAlbumsMap(sAlbums, bArtistMap)
	if err != nil {
		return nil, err
	}

	spotifyIds := sliceutils.Map(sTracks, func(item spotify.SavedTrack) spotify.ID { return item.ID })
	var existingTracks []*Track
	if err := b.gormDb.
		Preload("Artists").
		Where("spotify_id IN ?", spotifyIds).
		Find(&existingTracks).Error; err != nil {
		return nil, err
	}

	var newTracks []*Track
	for _, sTrack := range sTracks {
		if slices.ContainsFunc(existingTracks, func(item *Track) bool { return item.SpotifyId == sTrack.ID }) {
			continue
		}
		bArtists := sliceutils.Map(sTrack.Artists, func(item spotify.SimpleArtist) *Artist { return bArtistMap[item.ID] })
		bAlbum := bAlbumMap[sTrack.Album.ID]
		newTracks = append(newTracks, newTrack(sTrack, bAlbum, bArtists))
	}
	// All artists are already created, can exit
	if len(newTracks) == 0 {
		return existingTracks, nil
	}
	if err := b.gormDb.Create(newTracks).Error; err != nil {
		return nil, err
	}
	return append(existingTracks, newTracks...), nil
}
