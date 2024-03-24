package spotifyutils

import (
	"strings"

	"github.com/zmb3/spotify/v2"
)

func SArtistsString(artists []spotify.SimpleArtist) string {
	var s []string
	for _, a := range artists {
		s = append(s, a.Name)
	}
	return strings.Join(s, " ")
}
