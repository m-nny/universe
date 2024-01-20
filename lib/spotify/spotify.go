package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

func GetAllTracks(ctx context.Context) error {
	client, err := GetClient(ctx)
	if err != nil {
		return err
	}
	tracks, err := client.CurrentUsersTracks(ctx)
	if err != nil {
		return err
	}
	log.Printf("tracks: %+v", tracks)
	return err
}

func GetClient(ctx context.Context) (*spotify.Client, error) {
	clientId := os.Getenv("spotify_clientId")
	clientSecret := os.Getenv("spotify_clientSecret")
	redirectUrl := "http://localhost:3000/callback"
	state := "42"
	if clientId == "" || clientSecret == "" {
		return nil, fmt.Errorf("clientId or clientSecret is not set")
	}
	auth := spotifyauth.New(
		spotifyauth.WithClientID(clientId),
		spotifyauth.WithClientSecret(clientSecret),
		spotifyauth.WithRedirectURL(redirectUrl),
		spotifyauth.WithScopes(spotifyauth.ScopeUserLibraryRead),
	)
	url := auth.AuthURL(state)
	log.Printf("Login using following url:\n%s", url)
	clientCh := make(chan *spotify.Client)
	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			return
		}
		client := spotify.New(auth.Client(r.Context(), token))
		clientCh <- client
	}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/callback", callbackHandler)
	server := &http.Server{Addr: ":3000", Handler: serverMux}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Could not start server: %v", err)
		}
	}()
	client := <-clientCh
	if err := server.Shutdown(ctx); err != nil {
		return nil, fmt.Errorf("could not stop server: %v", err)
	}
	return client, nil
}
