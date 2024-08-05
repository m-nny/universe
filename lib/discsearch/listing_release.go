package discsearch

import (
	"context"
	"fmt"
	"regexp"

	"github.com/m-nny/universe/lib/brain"
	"golang.org/x/exp/slog"
)

func (a *App) FindRelease(ctx context.Context, bRelease *brain.DiscogsRelease) (*brain.MetaAlbum, error) {
	if bRelease.SearchedMetaAlbum {
		return nil, nil
	}
	q := sanitizeQ(fmt.Sprintf("%s %s", bRelease.ArtistName, bRelease.Name))
	sAlbums, err := a.Spotify.SearchAlbum(ctx, q, 1)
	if err != nil {
		return nil, err
	}
	bAlbums, err := a.Brain.SaveSimpleAlbumsSqlx(sAlbums)
	if err != nil {
		return nil, err
	}
	bMetaAlbum, score, err := brain.MostSimilarAlbum(bRelease, bAlbums)
	if err != nil {
		return nil, err
	}
	if bMetaAlbum == nil {
		slog.Error("discsearch.FindRelease(): not found", "bRelease.DiscogsID", bRelease.DiscogsID, "bRelease", bRelease)
	} else {
		slog.Debug("discsearch.FindRelease(): found", "album", bMetaAlbum.SimplifiedName, "score", score)
	}
	if err := a.Brain.AssociateDiscogsRelease(bRelease, bMetaAlbum, score); err != nil {
		return nil, err
	}
	return bMetaAlbum, nil
}

var sanitizeRgx = regexp.MustCompile(`[\(\)*\\\/\"\'\=\~\!\#\&\?]`)

func sanitizeQ(q string) string {
	q = sanitizeRgx.ReplaceAllString(q, "")
	return q
}
