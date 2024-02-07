package discogs

import (
	"context"
	"fmt"
)

func (d *Service) Release(ctx context.Context, releaseId int) (*Release, error) {
	var release Release
	discogsUrl := fmt.Sprintf("%s/releases/%d", d.config.BaseUrl, releaseId)
	if err := d.get(ctx, discogsUrl, &release); err != nil {
		return nil, err
	}
	return &release, nil
}

type Release struct {
	Artists     []Artist     `json:"artists"`
	ID          int          `json:"id"`
	Identifiers []Identifier `json:"identifiers"`
	MasterID    int          `json:"master_id"`
	MasterURL   string       `json:"master_url"`
	Title       string       `json:"title"`
	Tracks      []Track      `json:"tracklist"`
	URI         string       `json:"uri"`
	Year        int          `json:"year"`
}
