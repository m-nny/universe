package brain

import (
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/spotify/utils"
)

type albumIndex struct {
	idMap       map[spotify.ID]*Album
	simpNameMap map[string]*Album
}

func newAlbumIndex(bAlbums []*Album) *albumIndex {
	ai := &albumIndex{
		idMap:       make(map[spotify.ID]*Album),
		simpNameMap: make(map[string]*Album),
	}
	return ai.Add(bAlbums)
}

func (ai *albumIndex) Add(bAlbums []*Album) *albumIndex {
	for _, bAlbum := range bAlbums {
		ai.idMap[bAlbum.SpotifyId] = bAlbum
		ai.simpNameMap[bAlbum.SimplifiedName] = bAlbum
	}
	return ai
}

func (ai *albumIndex) Get(sAlbum spotify.SimpleAlbum) (*Album, bool) {
	if val, ok := ai.idMap[sAlbum.ID]; ok {
		return val, true
	}
	simpName := utils.SimplifiedAlbumName(sAlbum)
	if val, ok := ai.simpNameMap[simpName]; ok {
		return val, true
	}
	return nil, false
}
