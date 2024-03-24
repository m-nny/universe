package spotify_ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/ent/user"
)

func getTokenCached(ctx context.Context, auth *spotifyauth.Authenticator, ent *ent.Client, username string) (*oauth2.Token, error) {
	storedToken, err := getStoredToken(ctx, ent, username)
	if err == nil && storedToken.Valid() {
		return storedToken, nil
	}
	freshToken, err := getFreshToken(ctx, auth)
	if err != nil {
		return nil, err
	}
	if err := storeToken(ctx, ent, freshToken, username); err != nil {
		return nil, err
	}
	return freshToken, nil
}

func getFreshToken(ctx context.Context, auth *spotifyauth.Authenticator) (*oauth2.Token, error) {
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
	server := &http.Server{Addr: ":3000", Handler: serverMux}
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

func getStoredToken(ctx context.Context, ent *ent.Client, username string) (*oauth2.Token, error) {
	u, err := ent.User.
		Query().
		Where(user.ID(username)).
		Only(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed querying root user: %w", err)
	}
	return u.SpotifyToken, nil
}

func storeToken(ctx context.Context, ent *ent.Client, token *oauth2.Token, username string) error {
	return ent.User.
		Create().
		SetID(username).
		SetSpotifyToken(token).
		OnConflict().UpdateNewValues().
		Exec(ctx)
}
