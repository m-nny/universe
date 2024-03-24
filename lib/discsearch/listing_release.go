package discsearch

import (
	"context"
	"fmt"
	"regexp"

	"github.com/m-nny/universe/lib/brain"
)

func (a *App) FindRelease(ctx context.Context, bRelease *brain.DiscogsRelease) (*brain.MetaAlbum, int, error) {
	q := sanitizeQ(fmt.Sprintf("%s %s", bRelease.ArtistName, bRelease.Name))
	if len(q) > 50 {
		return nil, 0, nil
	}
	sAlbums, err := a.Spotify.SearchAlbum(ctx, q, 1)
	if err != nil {
		return nil, 0, err
	}
	bAlbums, err := a.Brain.SaveSimpleAlbums(sAlbums)
	if err != nil {
		return nil, 0, err
	}
	result, score, err := brain.MostSimilarAlbum(bRelease, bAlbums)
	return result, score, err
}

var sanitizeRgx = regexp.MustCompile(`[\(\)*\\\/\"\'\=\~\!\#\&\?]`)

func sanitizeQ(q string) string {
	q = sanitizeRgx.ReplaceAllString(q, "")
	return q
}
