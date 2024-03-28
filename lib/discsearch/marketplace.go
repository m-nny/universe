package discsearch

import (
	"context"
	"log"
	"slices"
	"strings"

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
	var missedBRelease []brain.DiscogsRelease
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
			missedBRelease = append(missedBRelease, *bRelease)
			continue
		}
		if err := a.Brain.AssociateDiscogsRelease(bRelease, bMetaAlbum); err != nil {
			return nil, err
		}
		bAlbums = append(bAlbums, bMetaAlbum)
	}
	log.Printf("Found %d bAlbums for %d bReleases for %d dReleases", len(bAlbums), len(bReleases), len(dReleases))
	slices.SortFunc(missedBRelease, func(a, b brain.DiscogsRelease) int {
		if val := strings.Compare(a.ArtistName, b.ArtistName); val != 0 {
			return val
		}
		return strings.Compare(a.Name, b.Name)
	})
	log.Printf("Missed %d bReleases", len(missedBRelease))
	for idx, bRelease := range missedBRelease {
		log.Printf("%3d. %08d %s - %s", idx+1, bRelease.DiscogsID, bRelease.ArtistName, bRelease.Name)
		log.Printf("              https://www.discogs.com/release/%d", bRelease.DiscogsID)
	}
	return bAlbums, nil
}
