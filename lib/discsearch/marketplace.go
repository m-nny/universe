package discsearch

import (
	"context"

	"github.com/m-nny/universe/ent"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/slices"
)

func (a *App) Inventory(ctx context.Context, sellerId string) ([]*ent.Album, error) {
	inventory, err := a.Discogs.SellerInventory(ctx, sellerId)
	if err != nil {
		return nil, err
	}
	albums, err := slices.MapCtxErr(ctx, inventory,
		func(ctx context.Context, release *discogs.Listing) (*ent.Album, error) {
			return a.ListingRelease(ctx, release.Release)
		})
	if err != nil {
		return nil, err
	}
	return albums, nil
}
