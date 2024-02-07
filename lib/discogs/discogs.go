package discogs

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/discogs/client"
	"github.com/m-nny/universe/lib/spotify"
)

type Config = client.Config

func LoadConfig() (*Config, error) {
	c := &Config{
		BaseUrl:   "https://api.discogs.com",
		UserAgent: os.Getenv("discogs_userAgent"),
		Currency:  "GBP",
		Token:     os.Getenv("discogs_token"),
	}
	if c.Token == "" {
		return nil, fmt.Errorf("discogs Token is not set")
	}
	return c, nil
}

type Service struct {
	spotify *spotify.Service
	ent     *ent.Client
	discogs *client.Discogs
}

func New(spotify *spotify.Service, ent *ent.Client) (*Service, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	discogs, err := client.New(config)
	if err != nil {
		return nil, err
	}
	return &Service{
		spotify: spotify,
		ent:     ent,
		discogs: discogs,
	}, nil
}

func (s *Service) GetRelease(ctx context.Context, id int) error {
	release, err := s.discogs.Release(ctx, id)
	if err != nil {
		return err
	}
	log.Printf("release: %+v", release)
	return nil
}
