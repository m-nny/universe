package discsearch

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

func (a *App) ListingRelease(ctx context.Context, release *discogs.ListingRelease) (*ent.Album, error) {
	q := sanitizeQ(fmt.Sprintf("%s %s", release.Artist, release.Title))
	log.Printf("q: %v", q)
	albums, err := a.Spotify.SearchAlbum(ctx, q)
	if err != nil {
		return nil, err
	}
	albums = sliceutils.Unique(albums, func(item *ent.Album) string { return fmt.Sprintf("%d", item.ID) })
	result, err := mostSimilarAlbum(ctx, release, albums)
	if err != nil {
		return nil, err
	}
	log.Printf("master: %s - %s", release.Artist, release.Title)
	if result != nil {
		log.Printf("album: %s", result.SimplifiedName)
	} else {
		log.Printf("album: not found %v", release)
	}
	return result, nil
}

var sanitizeRgx = regexp.MustCompile(`[\(\)*\\\/\"\']`)

func sanitizeQ(q string) string {
	return sanitizeRgx.ReplaceAllString(q, "")
}

func mostSimilarAlbum(ctx context.Context, m *discogs.ListingRelease, arr []*ent.Album) (result *ent.Album, err error) {
	maxScore := 0
	for _, item := range arr {
		if len(item.Edges.Artists) == 0 {
			artists, err := item.QueryArtists().All(ctx)
			if err != nil {
				return nil, err
			}
			item.Edges.Artists = artists
		}
		score := albumSimilarity(m, item)
		if score > maxScore {
			maxScore = score
			result = item
		}
	}
	return
}

func albumSimilarity(m *discogs.ListingRelease, a *ent.Album) int {
	artistScores := sliceutils.Map(a.Edges.Artists, func(e *ent.Artist) int { return similaryScore(m.Artist, e.Name) })
	artistScore := sliceutils.Cnt(artistScores, sliceutils.Identity)
	titleScore := similaryScore(m.Title, a.Name)
	return artistScore + titleScore
}

func similaryScore(a, b string) int {
	score := 0
	if strings.Contains(a, b) {
		score += len(b)
	}
	if strings.Contains(b, a) {
		score += len(a)
	}
	return score
}
