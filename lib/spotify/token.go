package spotify

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/m-nny/universe/lib/jsoncache"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func getTokenCached(ctx context.Context, auth *spotifyauth.Authenticator) (*oauth2.Token, error) {
	return jsoncache.CachedExec("spotify_token", func() (*oauth2.Token, error) {
		return getToken(ctx, auth)
	})
}

func getToken(ctx context.Context, auth *spotifyauth.Authenticator) (*oauth2.Token, error) {
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
