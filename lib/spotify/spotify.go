package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/jsoncache"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

type ID = spotify.ID

type SpotfyClient struct {
	ent    *ent.Client
	client *spotify.Client
}

type Config struct {
	ClientId     string
	ClientSecret string
	RedirectUrl  string
}

func LoadConfig() (*Config, error) {
	c := &Config{
		ClientId:     os.Getenv("spotify_clientId"),
		ClientSecret: os.Getenv("spotify_clientSecret"),
		RedirectUrl:  "http://localhost:3000/callback",
	}
	if c.ClientId == "" {
		return nil, fmt.Errorf("clientId ClientId is not set")
	}
	if c.ClientSecret == "" {
		return nil, fmt.Errorf("clientId ClientSecret is not set")
	}
	if c.RedirectUrl == "" {
		return nil, fmt.Errorf("clientId RedirectUrl is not set")
	}
	return c, nil
}

func New(ctx context.Context, config *Config, ent *ent.Client) (*SpotfyClient, error) {
	auth := spotifyauth.New(
		spotifyauth.WithClientID(config.ClientId),
		spotifyauth.WithClientSecret(config.ClientSecret),
		spotifyauth.WithRedirectURL(config.RedirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead),
	)
	token, err := getTokenCached(ctx, auth, ent)
	if err != nil {
		return nil, err
	}
	client := spotify.New(auth.Client(ctx, token))
	return &SpotfyClient{
		ent:    ent,
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
