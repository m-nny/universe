package spotify

import "github.com/zmb3/spotify/v2"

func getArtistNames(artists []spotify.SimpleArtist) (artistNames []string, artistIds []string) {
	for _, artist := range artists {
		artistNames = append(artistNames, artist.Name)
		artistIds = append(artistIds, artist.ID.String())
	}
	return
}
