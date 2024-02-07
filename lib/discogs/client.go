package discogs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

type Config struct {
	// Discogs API endpoint (optional).
	BaseUrl string
	// UserAgent to to call discogs api with.
	UserAgent string

	// Token provided by discogs (optional).
	Token string

	// Secret provided by discogs (optional).
	Secret string
	// Key provided by discogs (optional).
	Key string
}

type Service struct {
	http   *http.Client
	config *Config
}

func LoadConfig() (*Config, error) {
	c := &Config{
		BaseUrl:   "https://api.discogs.com",
		UserAgent: os.Getenv("discogs_userAgent"),
		Token:     os.Getenv("discogs_token"),
	}
	if c.Token == "" {
		return nil, fmt.Errorf("discogs Token is not set")
	}
	if c.UserAgent == "" {
		return nil, fmt.Errorf("config: UserAgent is empty")
	}
	if c.Secret == "" && c.Key != "" {
		return nil, fmt.Errorf("config: Secret is empty, but Key is not")
	}
	if c.Secret != "" && c.Key == "" {
		return nil, fmt.Errorf("config: Key is empty, but Secret is not")
	}
	return c, nil
}

func New() (*Service, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}
	return &Service{
		http:   &http.Client{},
		config: config,
	}, nil
}

func (d *Service) headers() http.Header {
	h := http.Header{}
	h.Add("User-Agent", d.config.UserAgent)
	if token := d.config.Token; token != "" {
		h.Add("Authorization", "Discogs token="+token)
	}
	if d.config.Secret != "" {
		h.Add("Authorization", fmt.Sprintf("Discogs secret=%s key=%s", d.config.Secret, d.config.Key))
	}
	return h
}

func get(ctx context.Context, fullUrl string, h http.Header, result any) error {
	log.Printf("[discogs] GET: %s", fullUrl)
	req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
	if err != nil {
		return err
	}
	req.Header = h
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("http status: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return err
	}
	return nil
}
