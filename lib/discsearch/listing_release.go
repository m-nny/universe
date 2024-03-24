package discsearch

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/m-nny/universe/lib/brain"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

func (a *App) ListingRelease(ctx context.Context, release discogs.ListingRelease) (*brain.MetaAlbum, error) {
	q := sanitizeQ(fmt.Sprintf("%s %s", release.Artist, release.Title))
	log.Printf("q: %v", q)
	sAlbums, err := a.Spotify.SearchAlbum(ctx, q)
	if err != nil {
		return nil, err
	}
	bAlbums, err := a.Brain.SaveSimpleAlbums(sAlbums)
	if err != nil {
		return nil, err
	}
	result, err := mostSimilarAlbum(release, bAlbums)
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

func mostSimilarAlbum(dRelease discogs.ListingRelease, bAlbums []*brain.MetaAlbum) (*brain.MetaAlbum, error) {
	var result *brain.MetaAlbum
	maxScore := 0
	for _, bAlbum := range bAlbums {
		if len(bAlbum.Artists) == 0 {
			return nil, fmt.Errorf("Albums should have artists populated")
		}
		score := albumSimilarity(dRelease, bAlbum)
		if score > maxScore {
			maxScore = score
			result = bAlbum
		}
	}
	return result, nil
}

func albumSimilarity(dRelease discogs.ListingRelease, eAlbum *brain.MetaAlbum) int {
	artistScore := sliceutils.Sum(eAlbum.Artists, func(e *brain.Artist) int { return similaryScore(dRelease.Artist, e.Name) })
	titleScore := similaryScore(dRelease.Title, eAlbum.AnyName)
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
