package spotify_ent

import (
	"context"
	"slices"

	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/lib/utils/sliceutils"
	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

func (s *Service) ToAlbums(ctx context.Context, arr []spotify.SimpleAlbum) ([]*ent.Album, error) {
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
