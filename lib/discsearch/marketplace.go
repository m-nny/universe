package discsearch

import (
	"context"
	"slices"
	"strings"

	"github.com/schollz/progressbar/v3"

	"github.com/m-nny/universe/lib/brain"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/logutils"
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
	bReleases, err := a.Brain.SaveDiscorgsReleases(dReleases, sellerId)
	if err != nil {
		return nil, err
	}
	var bAlbums []*brain.MetaAlbum
	var missedBRelease []brain.DiscogsRelease
	bar := progressbar.Default(int64(len(bReleases)))
	for _, bRelease := range bReleases {
		bar.Add(1)
		if bRelease.MetaAlbum != nil {
			bAlbums = append(bAlbums, bRelease.MetaAlbum)
			continue
		}
		// log.Printf("[%d/%d] release: %+v", idx+1, len(bReleases), bRelease.Name)
		// log.Printf("==================================")
		// log.Printf("release: %s - %s", bRelease.ArtistName, bRelease.Name)
		bMetaAlbum, err := a.FindRelease(ctx, bRelease)
		if err != nil {
			return nil, err
		}
		if bMetaAlbum == nil {
			missedBRelease = append(missedBRelease, *bRelease)
			continue
		}
		bAlbums = append(bAlbums, bMetaAlbum)
	}
	logutils.Debugf("Found %d bAlbums for %d bReleases for %d dReleases", len(bAlbums), len(bReleases), len(dReleases))
	slices.SortFunc(missedBRelease, func(a, b brain.DiscogsRelease) int {
		if val := strings.Compare(a.ArtistName, b.ArtistName); val != 0 {
			return val
		}
		return strings.Compare(a.Name, b.Name)
	})
	logutils.Debugf("Missed %d bReleases", len(missedBRelease))
	// for idx, bRelease := range missedBRelease {
	// 	log.Printf("%3d. %08d %s - %s", idx+1, bRelease.DiscogsID, bRelease.ArtistName, bRelease.Name)
	// 	log.Printf("              https://www.discogs.com/release/%d", bRelease.DiscogsID)
	// }
	return bAlbums, nil
}
