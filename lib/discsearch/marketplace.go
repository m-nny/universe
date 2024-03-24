package discsearch

import (
	"context"
	"log"

	"github.com/m-nny/universe/lib/brain"
	"github.com/m-nny/universe/lib/discogs"
)

func (a *App) Inventory(ctx context.Context, sellerId string) ([]*brain.MetaAlbum, error) {
	inventory, err := a.Discogs.SellerInventory(ctx, sellerId)
	if err != nil {
		return nil, err
	}
	// inventory = inventory[:10]
	dReleases := make([]discogs.ListingRelease, len(inventory))
	for idx, dListing := range inventory {
		dReleases[idx] = dListing.Release
	}
	bReleases, err := a.Brain.SaveDiscorgsReleases(dReleases)
	if err != nil {
		return nil, err
	}
	var bAlbums []*brain.MetaAlbum
	for idx, bRelease := range bReleases {
		log.Printf("[%d/%d] release: %+v", idx+1, len(bReleases), bRelease.Name)
		log.Printf("==================================")
		log.Printf("release: %s - %s", bRelease.ArtistName, bRelease.Name)
		if bRelease.MetaAlbum != nil {
			bAlbums = append(bAlbums, bRelease.MetaAlbum)
			continue
		}
		bMetaAlbum, score, err := a.FindRelease(ctx, bRelease)
		if err != nil {
			return nil, err
		}
		if bMetaAlbum != nil {
			log.Printf("album: %s score: %d", bMetaAlbum.SimplifiedName, score)
		} else {
			log.Printf("album: not found release_id: %d %+v", bRelease.DiscogsID, bRelease)
		}
		if bMetaAlbum == nil {
			continue
		}
		if err := a.Brain.AssociateDiscogsRelease(bRelease, bMetaAlbum); err != nil {
			return nil, err
		}
		bAlbums = append(bAlbums, bMetaAlbum)
	}
	log.Printf("Found %d bAlbums for %d bReleases for %d dReleases", len(bAlbums), len(bReleases), len(dReleases))
	log.Printf("%+v", bAlbums)
	return bAlbums, nil
}
