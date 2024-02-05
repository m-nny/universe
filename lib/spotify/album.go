package spotify

import (
	"context"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/album"
	"github.com/zmb3/spotify/v2"
)

func (s *Service) getAlbumFromSpotify(ctx context.Context, t spotify.SimpleAlbum) (*ent.Album, error) {
	return s.ent.Album.
		Query().
		Where(album.SpotifyIdContains(string(t.ID))).
		Only(ctx)
}

func (s *Service) toAlbum(ctx context.Context, a spotify.SimpleAlbum) (int, error) {
	// Check if already have it
	if album, err := s.getAlbumFromSpotify(ctx, a); err == nil {
		return album.ID, nil
	}

	artistIds, err := s.toArtists(ctx, a.Artists)
	if err != nil {
		return 0, err
	}
	spotifyIds := []string{string(a.ID)}
	album, err := s.ent.Album.Create().
		SetSpotifyIds(spotifyIds).
		SetName(string(a.Name)).
		AddArtistIDs(artistIds...).
		Save(ctx)
	if err != nil {
		return 0, err
	}
	return album.ID, nil
}
