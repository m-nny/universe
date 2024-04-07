package discsearch

import (
	"context"
	"fmt"
	"log"
	"regexp"

	"github.com/m-nny/universe/lib/brain"
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
		log.Printf("album: not found release_id: %d %+v", bRelease.DiscogsID, bRelease)
	} else {
		log.Printf("album: %s score: %d", bMetaAlbum.SimplifiedName, score)
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
