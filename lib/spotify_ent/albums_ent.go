package spotify_ent

import (
	"context"
	"slices"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/lib/jsoncache"
	"github.com/m-nny/universe/lib/utils/sliceutils"
	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

func (s *Service) SearchAlbum(ctx context.Context, q string) ([]*ent.Album, error) {
	results, err := jsoncache.CachedExec("discogs_search_"+q, func() (*spotify.SearchResult, error) {
		return s.spotify.Search(ctx, q, spotify.SearchTypeAlbum, spotify.Limit(50))
	})
	if err != nil {
		return nil, err
	}
	return s.toAlbums(ctx, results.Albums.Albums)
}

func (s *Service) toAlbumFull(ctx context.Context, a *spotify.FullAlbum) (*ent.Album, error) {
	album, err := s.toAlbum(ctx, a.SimpleAlbum)
	if err != nil {
		return nil, err
	}
	if _, err := s.toTracksWithAlbum(ctx, a.Tracks.Tracks, album); err != nil {
		return nil, err
	}
	return album, nil
}

func (s *Service) toAlbums(ctx context.Context, arr []spotify.SimpleAlbum) ([]*ent.Album, error) {
	return sliceutils.MapCtxErr(ctx, arr, s.toAlbum)
}

func (s *Service) toAlbum(ctx context.Context, a spotify.SimpleAlbum) (*ent.Album, error) {
	simplifiedName := utils.SimplifiedAlbumName(a)
	album, err := s.ent.Album.
		Query().
		Where(album.Similar(string(a.ID), simplifiedName)).
		Only(ctx)
	if err != nil {
		// Album does not exist yet
		return s._newAlbum(ctx, a, simplifiedName)
	}
	if slices.Contains(album.SpotifyIds, string(a.ID)) {
		return album, nil
	}
	// log.Printf("new album version: cur: %v new: %v", album, a.ID)
	return album.Update().AppendSpotifyIds([]string{string(a.ID)}).Save(ctx)
}

func (s *Service) _newAlbum(ctx context.Context, a spotify.SimpleAlbum, simplifiedName string) (*ent.Album, error) {
	artistIds, err := s.toArtists(ctx, a.Artists)
	if err != nil {
		return nil, err
	}
	return s.ent.Album.
		Create().
		AddArtistIDs(artistIds...).
		SetName(a.Name).
		SetSimplifiedName(simplifiedName).
		SetSpotifyIds([]string{string(a.ID)}).
		Save(ctx)
}
