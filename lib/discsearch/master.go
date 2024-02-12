package discsearch

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils"
)

func (a *App) Master(ctx context.Context, masterId int) (*ent.Album, error) {
	master, err := a.Discogs.Master(ctx, masterId)
	if err != nil {
		return nil, err
	}
	masterArtists := discogs.ArtistsString(master.Artists)
	q := fmt.Sprintf("%s %s", masterArtists, master.Title)
	log.Printf("q: %v", q)
	albums, err := a.Spotify.SearchAlbum(ctx, q)
	if err != nil {
		return nil, err
	}
	albums = utils.SliceUniqe(albums, func(item *ent.Album) string { return fmt.Sprintf("%d", item.ID) })
	result, err := mostSimilarAlbum(ctx, master, albums)
	if err != nil {
		return nil, err
	}
	if result == nil {
		return nil, fmt.Errorf("spotify: album not found")
	}
	log.Printf("master: %s - %s", masterArtists, master.Title)
	log.Printf("album: %s", result.SimplifiedName)
	return result, nil
}

func mostSimilarAlbum(ctx context.Context, m *discogs.Master, arr []*ent.Album) (result *ent.Album, err error) {
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

func albumSimilarity(m *discogs.Master, a *ent.Album) int {
	artistScores := utils.SliceStar(m.Artists, a.Edges.Artists,
		func(m *discogs.Artist, e *ent.Artist) int { return similaryScore(m.Name, e.Name) })
	artistScore := utils.SliceCnt(artistScores, utils.Identity)
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
