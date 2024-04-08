package brain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type User struct {
	Username string
	// TODO add proper token
	SpotifyTokenStr []byte `db:"spotify_token_str"`
	SavedTracks     []*MetaTrack
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
	if err := b.sqlxDb.GetContext(ctx, &user, "SELECT * FROM users WHERE username = ?", username); err != nil {
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
	if _, err := b.sqlxDb.NamedExecContext(ctx, `
		INSERT INTO users (username, spotify_token_str) VALUES (:username, :spotify_token_str)
		ON CONFLICT DO UPDATE SET spotify_token_str = excluded.spotify_token_str
	`, user); err != nil {
		return nil
	}
	return nil
}

type UserSavedTrack struct {
	Username    string `db:"user_username"`
	MetaTrackId uint   `db:"meta_track_id"`
}

func upsertUser(db *sqlx.DB, username string) error {
	if _, err := db.Exec(`
		INSERT INTO users (username)
		VALUES (?)
		ON CONFLICT DO NOTHING`, username); err != nil {
		return fmt.Errorf("could not upsert user: %w", err)
	}
	return nil
}

func addSavedTracksSqlx(db *sqlx.DB, username string, tracks []*MetaTrack) error {
	if err := upsertUser(db, username); err != nil {
		return err
	}
	var userSavedTracks []UserSavedTrack
	for _, bMetaTrack := range tracks {
		userSavedTracks = append(userSavedTracks, UserSavedTrack{username, bMetaTrack.ID})
	}
	if _, err := db.NamedExec(`
		INSERT INTO user_saved_tracks (user_username, meta_track_id)
		VALUES (:user_username, :meta_track_id)
		ON CONFLICT DO NOTHING`, userSavedTracks); err != nil {
		return err
	}
	return nil
}
