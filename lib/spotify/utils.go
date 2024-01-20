package spotify

import "github.com/zmb3/spotify/v2"

func getArtistNames(artists []spotify.SimpleArtist) (artistNames []string, artistIds []spotify.ID) {
	for _, artist := range artists {
		artistNames = append(artistNames, artist.Name)
		artistIds = append(artistIds, artist.ID)
	}
	return
}
