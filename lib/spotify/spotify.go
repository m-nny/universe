package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/m-nny/universe/lib/jsoncache"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type SpotfyClient struct {
	client *spotify.Client
}

func New(ctx context.Context) (*SpotfyClient, error) {
	clientId := os.Getenv("spotify_clientId")
	clientSecret := os.Getenv("spotify_clientSecret")
	redirectUrl := "http://localhost:3000/callback"
	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("clientId or clientSecret is not set")
	}
	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientId),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead),
	)
	token, err := getTokenCached(ctx, auth)
	if err != nil {
		return nil, err
	}
	client := spotify.New(auth.Client(ctx, token))
	return &SpotfyClient{
		client: client,
	}, nil
}

func (s *SpotfyClient) GetAllPlaylists(ctx context.Context) (playlists []spotify.SimplePlaylist, err error) {
	return jsoncache.CachedExec("spotify_savedPlaylists", func() ([]spotify.SimplePlaylist, error) {
		return s._GetAllPlaylists(ctx)
	})
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
