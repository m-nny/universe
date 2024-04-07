package spotifyutils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zmb3/spotify/v2"

	iterutils "github.com/m-nny/universe/lib/utils/iterutils"
)

var simplifyAlbumNameRegex = func() *regexp.Regexp {
	blocklistItems := []string{
		` \((\d{4} )?Remastered( Version)?\)`,
		` \(Bonus Edition\)`,
		` \(Collector's Edition\)`,
		` \(Deluxe Edition\)`,
		` \(Deluxe Version\)`,
		` \(Deluxe\)`,
		` \(Expanded Edition\)`,
		` \(Explicit Version\)`,
		` \(Extended Edition\)`,
		` \(Special Edition\)`,
		` \(Standard Version\)`,
		` \(The Complete Edition\)`,
		` \(Tour Edition\)`,
		` \(Wembley Edition\)`,
		` \(20th Anniversary Edition\)`,
		` Deluxe`,
	}

	regex := strings.Join(iterutils.Map(blocklistItems, func(s string) string {
		return fmt.Sprintf("(%s)", s)
	}), "|")
	return regexp.MustCompile(regex)
}()

// SimplifiedAlbumName will return a string in form of "<artist1>, <artist2> - <album name>"
func SimplifiedAlbumName(a spotify.SimpleAlbum) string {
	artistNames := strings.Join(
		iterutils.Map(a.Artists, func(a spotify.SimpleArtist) string { return a.Name }),
		", ",
	)
	// releaseYear := a.ReleaseDateTime().Year()
	// albumName := strings.ReplaceAll(a.Name, "(Deluxe Edition)", "")
	albumName := simplifyAlbumNameRegex.ReplaceAllString(a.Name, "")
	msg := fmt.Sprintf("%s - %s", artistNames, albumName)
	msg = strings.ToLower(msg)
	return msg
}
