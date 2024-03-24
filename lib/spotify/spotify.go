package spotify

import (
	"context"
	"fmt"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/brain"
)

type ID = spotify.ID

type Service struct {
	brain   *brain.Brain
	spotify *spotify.Client
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
		return nil, fmt.Errorf("spotify ClientId is not set")
	}
	if c.ClientSecret == "" {
		return nil, fmt.Errorf("spotify ClientSecret is not set")
	}
	if c.RedirectUrl == "" {
		return nil, fmt.Errorf("spotify RedirectUrl is not set")
	}
	return c, nil
}

func New(ctx context.Context, ent *ent.Client, brain *brain.Brain, username string) (*Service, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	auth := spotifyauth.New(
		spotifyauth.WithClientID(config.ClientId),
		spotifyauth.WithClientSecret(config.ClientSecret),
		spotifyauth.WithRedirectURL(config.RedirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead),
	)
	token, err := getTokenCached(ctx, auth, ent, username)
	if err != nil {
		return nil, err
	}
	client := spotify.New(auth.Client(ctx, token))
	return &Service{
		brain:   brain,
		spotify: client,
	}, nil
}
