package token

import (
	"context"

	"golang.org/x/oauth2"
)

type InMemoryTokenStorage struct {
	tokenMap map[string]*oauth2.Token
}

var _ TokenStorage = (*InMemoryTokenStorage)(nil)

func (e *InMemoryTokenStorage) GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error) {
	return e.tokenMap[username], nil
}

func (e *InMemoryTokenStorage) StoreSpotifyToken(ctx context.Context, username string, token *oauth2.Token) error {
	e.tokenMap[username] = token
	return nil
}

func NewInMemoryTokenStorage() *InMemoryTokenStorage {
	return &InMemoryTokenStorage{tokenMap: make(map[string]*oauth2.Token)}
}
