package spotify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	"golang.org/x/oauth2"
)

const tokenCacheLocation = ".cache/spotify.json"

func getToken(ctx context.Context) (*oauth2.Token, error) {
	if cachedToken, err := getCachedToken(); err == nil && cachedToken != nil {
		return cachedToken, nil
	}
	token, err := getTokenUsingApi(ctx)
	if err != nil {
		return nil, err
	}
	if err := setCachedToken(token); err != nil {
		return nil, err
	}
	return token, nil
}

func getTokenUsingApi(ctx context.Context) (*oauth2.Token, error) {
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

func setCachedToken(token *oauth2.Token) error {
	dir := path.Dir(tokenCacheLocation)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	tokenJson, err := json.MarshalIndent(token, "  ", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(tokenCacheLocation, tokenJson, os.ModePerm)
}

func getCachedToken() (*oauth2.Token, error) {
	data, err := os.ReadFile(tokenCacheLocation)
	if err != nil {
		return nil, err
	}
	var token *oauth2.Token

	if err := json.Unmarshal(data, &token); err != nil {
		return nil, err
	}
	return token, nil
}
