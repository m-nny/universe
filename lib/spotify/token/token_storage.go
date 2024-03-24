package token

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

type TokenStorage interface {
	GetSpotifyToken(ctx context.Context, username string) (*oauth2.Token, error)
	StoreSpotifyToken(ctx context.Context, username string, spotifyToken *oauth2.Token) error
}

func GetToken(ctx context.Context, auth *spotifyauth.Authenticator, serverAddress string, tokenStorage TokenStorage, username string) (*oauth2.Token, error) {
	storedToken, err := tokenStorage.GetSpotifyToken(ctx, username)
	{
		var tokenExpiry time.Time
		if storedToken != nil {
			tokenExpiry = storedToken.Expiry
		}
		log.Printf("storedToken.Expiriry: %v storedToken.Valid(): %v err: %v", tokenExpiry, storedToken.Valid(), err)
	}
	if err == nil && storedToken.Valid() {
		return storedToken, nil
	}
	freshToken, err := GetFreshToken(ctx, auth, serverAddress)
	if err != nil {
		return nil, err
	}
	if err := tokenStorage.StoreSpotifyToken(ctx, username, freshToken); err != nil {
		return nil, err
	}
	return freshToken, nil
}

func GetFreshToken(ctx context.Context, auth *spotifyauth.Authenticator, serverAddress string) (*oauth2.Token, error) {
	state := "42"
	url := auth.AuthURL(state)
	log.Printf("Login using following url:\n%s", url)
	tokenCh := make(chan *oauth2.Token)
	errCh := make(chan error)
	callbackHandler := func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.Token(r.Context(), state, r)
		if err != nil {
			http.Error(w, "Couldn't get token", http.StatusNotFound)
			errCh <- err
			return
		}
		tokenCh <- token
	}
	serverMux := http.NewServeMux()
	serverMux.HandleFunc("/callback", callbackHandler)
	server := &http.Server{Addr: serverAddress, Handler: serverMux}
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Could not start server: %v", err)
			errCh <- err
		}
	}()
	select {
	case token := <-tokenCh:
		if err := server.Shutdown(ctx); err != nil {
			return nil, fmt.Errorf("could not stop server: %v", err)
		}
		return token, nil
	case err := <-errCh:
		return nil, err
	}
}
