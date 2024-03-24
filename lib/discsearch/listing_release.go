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

func (a *App) ListingRelease(ctx context.Context, release discogs.ListingRelease) (*ent.Album, error) {
	q := sanitizeQ(fmt.Sprintf("%s %s", release.Artist, release.Title))
	log.Printf("q: %v", q)
	salbums, err := a.Spotify.SearchAlbum(ctx, q)
	if err != nil {
		return nil, err
	}
	albums, err := a.SpotifyEnt.ToAlbums(ctx, salbums)
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

func mostSimilarAlbum(ctx context.Context, dRelease discogs.ListingRelease, eAlbums []*ent.Album) (*ent.Album, error) {
	var result *ent.Album
	maxScore := 0
	for _, eAlbum := range eAlbums {
		if len(eAlbum.Edges.Artists) == 0 {
			artists, err := eAlbum.QueryArtists().All(ctx)
			if err != nil {
				return nil, err
			}
			eAlbum.Edges.Artists = artists
		}
		score := albumSimilarity(dRelease, eAlbum)
		if score > maxScore {
			maxScore = score
			result = eAlbum
		}
	}
	return result, nil
}

func albumSimilarity(dRelease discogs.ListingRelease, eAlbum *ent.Album) int {
	artistScore := sliceutils.Sum(eAlbum.Edges.Artists, func(e *ent.Artist) int { return similaryScore(dRelease.Artist, e.Name) })
	titleScore := similaryScore(dRelease.Title, eAlbum.Name)
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
