package spotifyutils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/xiter"
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

	regex := xiter.SliceJoinFn(blocklistItems, "|", func(s string) string { return fmt.Sprintf("(%s)", s) })
	return regexp.MustCompile(regex)
}()

// SimplifiedAlbumName will return a string in form of "<artist1>, <artist2> - <album name>"
func SimplifiedAlbumName(a spotify.SimpleAlbum) string {
	artistNames := xiter.SliceJoinFn(a.Artists, ", ", func(a spotify.SimpleArtist) string { return a.Name })
	// releaseYear := a.ReleaseDateTime().Year()
	// albumName := strings.ReplaceAll(a.Name, "(Deluxe Edition)", "")
	albumName := simplifyAlbumNameRegex.ReplaceAllString(a.Name, "")
	msg := fmt.Sprintf("%s - %s", artistNames, albumName)
	msg = strings.ToLower(msg)
	return msg
}
