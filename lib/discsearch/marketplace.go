package discsearch

import (
	"context"

	"github.com/m-nny/universe/lib/brain"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

func (a *App) Inventory(ctx context.Context, sellerId string) ([]*brain.MetaAlbum, error) {
	inventory, err := a.Discogs.SellerInventory(ctx, sellerId)
	if err != nil {
		return nil, err
	}
	inventory = inventory[:10]
	albums, err := sliceutils.MapCtxErr(ctx, inventory,
		func(ctx context.Context, release discogs.Listing) (*brain.MetaAlbum, error) {
			return a.ListingRelease(ctx, release.Release)
		})
	if err != nil {
		return nil, err
	}
	return albums, nil
}
