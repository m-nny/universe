package brain

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

type User struct {
	Username string `gorm:"primarykey"`
	// TODO add proper token
	SpotifyTokenStr string
}

func (u *User) SpotifyToken() (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal([]byte(u.SpotifyTokenStr), &token); err != nil {
		return nil, err
	}
	return &token, nil
}

func newUser(username string, spotifyToken *oauth2.Token) (*User, error) {
	spotifyTokenStr, err := json.Marshal(spotifyToken)
	if err != nil {
		return nil, err
	}
	return &User{
		Username:        username,
		SpotifyTokenStr: string(spotifyTokenStr),
	}, nil
}

func (b *Brain) GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error) {
	var user User
	if err := b.gormDb.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return user.SpotifyToken()
}

func (b *Brain) StoreSpotifyToken(ctx context.Context, username string, spotifyToken *oauth2.Token) error {
	user, err := newUser(username, spotifyToken)
	if err != nil {
		return err
	}
	return b.gormDb.Create(&user).Error
}
