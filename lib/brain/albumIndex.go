package brain

import (
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/spotify/utils"
)

type brainIndex struct {
	artistsMap      map[spotify.ID]*Artist
	metaAlbumMap    map[string]*MetaAlbum
	spotifyAlbumMap map[spotify.ID]*SpotifyAlbum
}

func newBrainIndex() *brainIndex {
	return &brainIndex{
		artistsMap:      make(map[spotify.ID]*Artist),
		metaAlbumMap:    make(map[string]*MetaAlbum),
		spotifyAlbumMap: make(map[spotify.ID]*SpotifyAlbum),
	}
}

func (ai *brainIndex) AddMetaAlbums(bMetaAlbums []*MetaAlbum) *brainIndex {
	for _, bAlbum := range bMetaAlbums {
		ai.metaAlbumMap[bAlbum.SimplifiedName] = bAlbum
	}
	return ai
}

func (ai *brainIndex) GetMetaAlbum(sAlbum spotify.SimpleAlbum) (*MetaAlbum, bool) {
	simpName := utils.SimplifiedAlbumName(sAlbum)
	bMetaAlbum, ok := ai.metaAlbumMap[simpName]
	return bMetaAlbum, ok
}

func (ai *brainIndex) AddSpotifyAlbums(bSpotifyAlbums []*SpotifyAlbum) *brainIndex {
	for _, bAlbum := range bSpotifyAlbums {
		ai.spotifyAlbumMap[bAlbum.SpotifyId] = bAlbum
	}
	return ai
}

func (ai *brainIndex) GetSpotifyAlbum(sAlbum spotify.SimpleAlbum) (*SpotifyAlbum, bool) {
	bSpotifyAlbum, ok := ai.spotifyAlbumMap[sAlbum.ID]
	return bSpotifyAlbum, ok
}

func (ai *brainIndex) AddArtists(bArtists []*Artist) *brainIndex {
	for _, bArtist := range bArtists {
		ai.artistsMap[bArtist.SpotifyId] = bArtist
	}
	return ai
}

func (ai *brainIndex) GetArtist(sArtist spotify.SimpleArtist) (*Artist, bool) {
	bArtist, ok := ai.artistsMap[sArtist.ID]
	return bArtist, ok
}

func (ai *brainIndex) GetArtists(sArtists []spotify.SimpleArtist) ([]*Artist, bool) {
	var bArtists []*Artist
	for _, sArtist := range sArtists {
		bArtist, ok := ai.artistsMap[sArtist.ID]
		if !ok {
			return nil, false
		}
		bArtists = append(bArtists, bArtist)
	}
	return bArtists, true
}
