package spotify

import (
	"context"
	"fmt"

	"github.com/m-nny/universe/ent"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) saveArtistsDeep(ctx context.Context, tracks []spotify.SavedTrack, username string) ([]*ent.Track, error) {
	var rawArtists []spotify.SimpleArtist
	for _, track := range tracks {
		for _, artist := range track.Artists {
			rawArtists = append(rawArtists, artist)
		}
	}
	// artistMap, err := s.batchToArtists(ctx, rawArtists)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, fmt.Errorf("saveArtistsDeep: not implemnted")
}
