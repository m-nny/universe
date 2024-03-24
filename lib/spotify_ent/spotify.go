package spotify_ent

import (
	"context"

	"github.com/m-nny/universe/ent"
)

type Service struct {
	ent     *ent.Client
}

func New(ctx context.Context, ent *ent.Client, username string) *Service {
	return &Service{
		ent:     ent,
	}
}
