package spotify_ent

import (
	"context"
)

func (s *Service) upsertUser(ctx context.Context, username string) error {
	return s.ent.User.
		Create().
		SetID(username).
		OnConflict().UpdateNewValues().
		Exec(ctx)
}
