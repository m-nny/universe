package brain

import (
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

// batchSaveAlbumTracks returns Brain representain of a spotify albums and tracks
//   - It will create new entries in DB if necessary for albums, tracks, artists
//   - It will deduplicate returned albums base on spotify.ID, this may result in len(result) < len(sAlbums)
//   - NOTE: it does not store all spotify.IDs of duplicated at the moment
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) batchSaveAlbumTracks(sAlbums []spotify.SimpleAlbum, sTracks []spotify.SimpleTrack) ([]*MetaAlbum, []*MetaTrack, error) {
	sAlbums = sliceutils.Unique(sAlbums, func(item spotify.SimpleAlbum) spotify.ID { return item.ID })
	sTracks = sliceutils.Unique(sTracks, func(item spotify.SimpleTrack) spotify.ID { return item.ID })

	var sArtists []spotify.SimpleArtist
	for _, sAlbum := range sAlbums {
		sArtists = append(sArtists, sAlbum.Artists...)
	}
	for _, sTrack := range sTracks {
		sArtists = append(sArtists, sTrack.Artists...)
	}

	bi := newBrainIndex()
	_, err := upsertArtists(b, sArtists, bi)
	if err != nil {
		return nil, nil, err
	}

	metaAlbums, err := upsertMetaAlbums(b, sAlbums, bi)
	if err != nil {
		return nil, nil, err
	}
	if _, err := upsertSpotifyAlbums(b, sAlbums, bi); err != nil {
		return nil, nil, err
	}

	metaTracks, err := upsertMetaTracks(b, sTracks, bi)
	if err != nil {
		return nil, nil, err
	}
	if _, err := upsertSpotifyTracks(b, sTracks, bi); err != nil {
		return nil, nil, err
	}
	return metaAlbums, metaTracks, nil
}
