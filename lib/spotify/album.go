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
	return utils.SliceMapCtxErr(
		ctx,
		rawAlbums,
		func(ctx context.Context, a *spotify.FullAlbum) (*ent.Album, error) {
			return s.toAlbum(ctx, a.SimpleAlbum)
		})
}

func (s *Service) getAlbumFromSpotify(ctx context.Context, a spotify.SimpleAlbum) (*ent.Album, error) {
	simplifiedName := simplifiedAlbumName(a)
	return s.ent.Album.
		Query().
		Where(
			album.Or(
				album.SpotifyIdContains(string(a.ID)),
				album.SimplifiedName(simplifiedName),
			),
		).
		Only(ctx)
}

func (s *Service) toAlbum(ctx context.Context, rawAlbum spotify.SimpleAlbum) (*ent.Album, error) {
	album, err := s.getAlbumFromSpotify(ctx, rawAlbum)
	if err != nil {
		// Album does not exist yet
		return s._newAlbum(ctx, rawAlbum)
	}
	if slices.Contains(album.SpotifyIds, string(rawAlbum.ID)) {
		return album, nil
	}
	log.Printf("new album version: cur: %v new: %v", album, rawAlbum.ID)
	return album.Update().AppendSpotifyIds([]string{string(rawAlbum.ID)}).Save(ctx)
}

func (s *Service) _newAlbum(ctx context.Context, a spotify.SimpleAlbum) (*ent.Album, error) {
	artistIds, err := s.toArtists(ctx, a.Artists)
	if err != nil {
		return nil, err
	}
	spotifyIds := []string{string(a.ID)}
	simplfiedName := simplifiedAlbumName(a)
	return s.ent.Album.
		Create().
		SetSpotifyIds(spotifyIds).
		SetName(string(a.Name)).
		SetSimplifiedName(simplfiedName).
		AddArtistIDs(artistIds...).
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
