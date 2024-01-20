package spotify

import "github.com/zmb3/spotify/v2"

var _albums = make(map[spotify.ID]*Album)

type Album struct {
	Name        string
	Id          ID
	ArtistNames []string
	ArtistIds   []ID
}

func getAlbum(simpleAlbum spotify.SimpleAlbum) *Album {
	id := simpleAlbum.ID
	if val := _albums[id]; val != nil {
		return val
	}
	artistNames, artistIds := getArtistNames(simpleAlbum.Artists)
	album := &Album{
		Name:        simpleAlbum.Name,
		Id:          id,
		ArtistNames: artistNames,
		ArtistIds:   artistIds,
	}
	_albums[id] = album
	return album
}
