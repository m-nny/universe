package spotify

import "github.com/zmb3/spotify/v2"

func getArtistNames(artists []spotify.SimpleArtist) (artistIds []string) {
	for _, artist := range artists {
		artistIds = append(artistIds, artist.ID.String())
	}
	return
}
