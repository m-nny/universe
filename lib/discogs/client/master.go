package client

import (
	"context"
	"fmt"
)

func (d *Discogs) Master(ctx context.Context, releaseId int) (*Master, error) {
	var release Master
	discogsUrl := fmt.Sprintf("%s/masters/%d", d.config.BaseUrl, releaseId)
	if err := d.get(ctx, discogsUrl, &release); err != nil {
		return nil, err
	}
	return &release, nil
}

type Master struct {
	Artists        []Artist `json:"artists"`
	ID             int      `json:"id"`
	MainRelease    int      `json:"main_release"`
	MainReleaseURL string   `json:"main_release_url"`
	ResourceURL    string   `json:"resource_url"`
	Title          string   `json:"title"`
	Tracks         []Track  `json:"tracklist"`
	URI            string   `json:"uri"`
	Year           int      `json:"year"`
}
