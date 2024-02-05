package spotify

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/album"
	"github.com/m-nny/universe/lib/utils"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) GetAlbumsById(ctx context.Context, ids []spotify.ID) ([]*ent.Album, error) {
	rawAlbums, err := s.spotify.GetAlbums(ctx, ids)
	if err != nil {
		return nil, err
	}
	return utils.SliceMapCtxErr(ctx, rawAlbums, s.toAlbumFull)
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
func (s *Service) toAlbum(ctx context.Context, a spotify.SimpleAlbum) (*ent.Album, error) {
	simplifiedName := simplifiedAlbumName(a)
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
	log.Printf("new album version: cur: %v new: %v", album, a.ID)
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

// simplifiedAlbumName will return a string in form of "<artist1>, <artist2> - <album name> [<album release year]"
func simplifiedAlbumName(a spotify.SimpleAlbum) string {
	artistNames := strings.Join(
		utils.SliceMap(a.Artists, func(a spotify.SimpleArtist) string { return a.Name }),
		", ",
	)
	releaseYear := a.ReleaseDateTime().Year()
	msg := fmt.Sprintf("%s - %s [%d]", artistNames, a.Name, releaseYear)
	msg = strings.ToLower(msg)
	return msg
}
