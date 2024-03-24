package token

import (
	"context"
	"fmt"

	"golang.org/x/oauth2"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/user"
)

type EntTokenStorage struct {
	ent *ent.Client
}

var _ TokenStorage = (*EntTokenStorage)(nil)

func (e *EntTokenStorage) GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error) {
	u, err := e.ent.User.
		Query().
		Where(user.ID(username)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying root user: %w", err)
	}
	return u.SpotifyToken, nil
}

func (e *EntTokenStorage) StoreSpotifyToken(ctx context.Context, username string, token *oauth2.Token) error {
	return e.ent.User.
		Create().
		SetID(username).
		SetSpotifyToken(token).
		OnConflict().UpdateNewValues().
		Exec(ctx)
}

func NewEntTokenStorage(ent *ent.Client) *EntTokenStorage {
	return &EntTokenStorage{ent: ent}
}
