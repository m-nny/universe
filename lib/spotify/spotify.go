package spotify

import (
	"context"
	"log"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

var auth *spotifyauth.Authenticator

func getAuth() *spotifyauth.Authenticator {
	clientId := os.Getenv("spotify_clientId")
	clientSecret := os.Getenv("spotify_clientSecret")
	redirectUrl := "http://localhost:3000/callback"
	if clientId == "" || clientSecret == "" {
		log.Fatalf("clientId or clientSecret is not set")
	}
	return spotifyauth.New(
		spotifyauth.WithClientID(clientId),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead),
	)
}

func GetAllTracks(ctx context.Context) error {
	auth = getAuth()
	token, err := getToken(ctx)
	if err != nil {
		return err
	}
	client := spotify.New(auth.Client(ctx, token))
	tracks, err := client.CurrentUsersTracks(ctx)
	if err != nil {
		return err
	}
	log.Printf("tracks: %+v", tracks)
	return err
}
