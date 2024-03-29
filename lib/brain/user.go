package brain

import (
	"context"
	"encoding/json"
	"log"

	"golang.org/x/oauth2"
	"gorm.io/gorm/clause"
)

type User struct {
	Username string `gorm:"primarykey"`
	// TODO add proper token
	SpotifyTokenStr []byte
}

func (u *User) SpotifyToken() (*oauth2.Token, error) {
	var token oauth2.Token
	if err := json.Unmarshal(u.SpotifyTokenStr, &token); err != nil {
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
		SpotifyTokenStr: spotifyTokenStr,
	}, nil
}

func (b *Brain) GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error) {
	var user User
	if err := b.gormDb.Where(User{Username: username}).First(&user).Error; err != nil {
		return nil, err
	}
	userToken, err := user.SpotifyToken()
	if err != nil {
		return nil, err
	}
	log.Printf("GetSpotifyToken(): userToken.tokenExpiry: %v", userToken.Expiry)
	return userToken, nil
}

func (b *Brain) StoreSpotifyToken(ctx context.Context, username string, spotifyToken *oauth2.Token) error {
	user, err := newUser(username, spotifyToken)
	if err != nil {
		return err
	}
	if err := b.gormDb.Clauses(clause.OnConflict{UpdateAll: true}).Create(&user).Error; err != nil {
		return nil
	}
	return nil
}
