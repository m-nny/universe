package spotify

import (
	"context"
	"errors"
	"log"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/user"
	"github.com/zmb3/spotify/v2"
)

func (s *SpotfyClient) GetAllPlaylists(ctx context.Context) ([]*ent.Playlist, error) {
	plists, err := s.getUserPlaylists(ctx)
	if err == nil && len(plists) > 0 {
		return plists, nil
	}
	rawPlists, err := s._GetAllPlaylists(ctx)
	if err != nil {
		return nil, err
	}
	if err := s.saveUserPlaylists(ctx, rawPlists); err != nil {
		return nil, err
	}
	return s.getUserPlaylists(ctx)
}
func (s *SpotfyClient) _GetAllPlaylists(ctx context.Context) (playlists []spotify.SimplePlaylist, err error) {
	for resp, err := s.client.CurrentUsersPlaylists(ctx); err == nil; err = s.client.NextPage(ctx, resp) {
		log.Printf("len(resp.Tracks)=%d", len(resp.Playlists))
		playlists = append(playlists, resp.Playlists...)
	}
	if !errors.Is(err, spotify.ErrNoMorePages) && err != nil {
		return nil, err
	}
	return playlists, nil
}

func (s *SpotfyClient) getUserPlaylists(ctx context.Context) ([]*ent.Playlist, error) {
	plists, err := s.ent.User.
		Query().
		Where(user.ID(rootUserName)).
		QueryPlaylists().
		All(ctx)
	return plists, err
}

func (s *SpotfyClient) saveUserPlaylists(ctx context.Context, rawPlists []spotify.SimplePlaylist) error {
	err := s.ent.Playlist.
		MapCreateBulk(rawPlists, func(pc *ent.PlaylistCreate, i int) {
			p := rawPlists[i]
			pc.
				SetID(string(p.ID)).
				SetName(p.Name).
				SetSnaphotID(p.SnapshotID).
				SetOwnerID(rootUserName)
		}).
		OnConflict().
		UpdateNewValues().
		Exec(ctx)
	return err
}
