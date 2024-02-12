package discogs

import (
	"context"
	"fmt"
	"strings"
)

func (d *Service) Release(ctx context.Context, releaseId int) (*Release, error) {
	var release Release
	discogsUrl := fmt.Sprintf("%s/releases/%d", d.config.BaseUrl, releaseId)
	if err := get(ctx, discogsUrl, d.headers(), &release); err != nil {
		return nil, err
	}
	return &release, nil
}

type Release struct {
	Artists     []*Artist    `json:"artists"`
	ID          int          `json:"id"`
	Identifiers []Identifier `json:"identifiers"`
	MasterID    int          `json:"master_id"`
	MasterURL   string       `json:"master_url"`
	Title       string       `json:"title"`
	Tracks      []Track      `json:"tracklist"`
	URI         string       `json:"uri"`
	Year        int          `json:"year"`
}

func ArtistsString(artists []*Artist) string {
	var s []string
	for _, a := range artists {
		s = append(s, a.Name)
	}
	return strings.Join(s, " ")
}
