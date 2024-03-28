package discogs

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/m-nny/universe/lib/jsoncache"
)

func (d *Service) SellerInventory(ctx context.Context, username string) ([]Listing, error) {
	return jsoncache.CachedExec("discogs/inventory/"+username, func() ([]Listing, error) {
		return d._SellerInventory(ctx, username)
	})
}

func (d *Service) _SellerInventory(ctx context.Context, username string) ([]Listing, error) {
	var results []Listing
	params := url.Values{}
	params.Set("status", "For Sale")
	params.Set("per_page", "100")
	discogsUrl := fmt.Sprintf("%s/users/%s/inventory?%s", d.config.BaseUrl, username, params.Encode())
	var inventory inventoryResult
	for {
		if err := get(ctx, discogsUrl, d.headers(), &inventory); err != nil {
			return nil, err
		}
		results = append(results, inventory.Listings...)
		if discogsUrl == inventory.nextPage() {
			break
		}
		discogsUrl = inventory.nextPage()
	}
	log.Printf("Inventory.totalItems=%d len(results)=%d", inventory.totalItems(), len(results))
	if inventory.totalItems() != len(results) {
		return results, fmt.Errorf("could not get all items got:%d total:%d", len(results), inventory.totalItems())
	}
	return results, nil
}

type inventoryResult struct {
	Pagination Pagination `json:"pagination"`
	Listings   []Listing  `json:"listings"`
}

func (i *inventoryResult) nextPage() string {
	return i.Pagination.Urls["next"]
}

func (i *inventoryResult) totalItems() int {
	return i.Pagination.Items
}

type Listing struct {
	AllowOffers     bool           `json:"allow_offers"`
	Audio           bool           `json:"audio"`
	Comments        string         `json:"comments"`
	Condition       string         `json:"condition"`
	ID              int            `json:"id"`
	Posted          string         `json:"posted"`
	Price           Price          `json:"price"`
	Release         ListingRelease `json:"release,omitempty"`
	ResourceURL     string         `json:"resource_url"`
	ShipsFrom       string         `json:"ships_from"`
	SleeveCondition string         `json:"sleeve_condition"`
	Status          string         `json:"status"`
	URI             string         `json:"uri"`
}
type Pagination struct {
	Items   int               `json:"items"`
	Page    int               `json:"page"`
	Pages   int               `json:"pages"`
	PerPage int               `json:"per_page"`
	Urls    map[string]string `json:"urls"`
}
type Price struct {
	Currency string  `json:"currency"`
	Value    float64 `json:"value"`
}

type ListingRelease struct {
	Artist        string `json:"artist"`
	CatalogNumber string `json:"catalog_number"`
	Description   string `json:"description"`
	Format        string `json:"format"`
	ID            int    `json:"id"`
	ResourceURL   string `json:"resource_url"`
	Thumbnail     string `json:"thumbnail"`
	Title         string `json:"title"`
	Year          int    `json:"year"`
}
