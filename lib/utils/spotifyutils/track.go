package spotifyutils

import (
	"fmt"
	"strings"

	"github.com/zmb3/spotify/v2"
)

// simplifiedTrackName will return a string in form of
//
//	"<artist1>, <artist2> - <album name> [<album release year] - <track_num>. <track_name>"
func SimplifiedTrackName(t spotify.SimpleTrack, albumSimplifiedName string) string {
	msg := albumSimplifiedName
	msg += fmt.Sprintf(" - %02d.  %s", t.TrackNumber, t.Name)
	msg = strings.ToLower(msg)
	return msg
}
