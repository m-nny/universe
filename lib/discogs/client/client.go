package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
	// Discogs API endpoint (optional).
	BaseUrl string
	// Currency to use (optional, default is USD).
	Currency string
	// UserAgent to to call discogs api with.
	UserAgent string

	// Token provided by discogs (optional).
	Token string

	// Secret provided by discogs (optional).
	Secret string
	// Key provided by discogs (optional).
	Key string
}

type Discogs struct {
	http   *http.Client
	config *Config
}

func New(config *Config) (*Discogs, error) {
	if config.UserAgent == "" {
		return nil, fmt.Errorf("config: UserAgent is empty")
	}
	if config.Secret == "" && config.Key != "" {
		return nil, fmt.Errorf("config: Secret is empty, but Key is not")
	}
	if config.Secret != "" && config.Key == "" {
		return nil, fmt.Errorf("config: Key is empty, but Secret is not")
	}
	return &Discogs{
		http:   &http.Client{},
		config: config,
	}, nil
}

func (d *Discogs) headers() http.Header {
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

func (d *Discogs) get(ctx context.Context, fullUrl string, result any) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fullUrl, nil)
	if err != nil {
		return err
	}
	req.Header = d.headers()
	resp, err := d.http.Do(req)
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
