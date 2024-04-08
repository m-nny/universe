package brain

import (
	"log"

	"github.com/jmoiron/sqlx"
)

type DiscogsSeller struct {
	Username        string
	SellingReleases []*DiscogsRelease
}

type discogsSellerReleases struct {
	DiscogsSeller string `db:"discogs_seller_username"`
	ReleaseId     uint   `db:"discogs_release_id"`
}

func upsertDiscogsSeller(db *sqlx.DB, username string) error {
	if _, err := db.Exec(`
		INSERT INTO discogs_sellers (username)
		VALUES (?)
		ON CONFLICT DO NOTHING
		`, username); err != nil {
		return err
	}
	return nil
}

func addDiscogsReleases(db *sqlx.DB, username string, releases []*DiscogsRelease) error {
	if len(releases) == 0 {
		return nil
	}
	if err := upsertDiscogsSeller(db, username); err != nil {
		return err
	}
	var sellingReleases []discogsSellerReleases
	for _, item := range releases {
		sellingReleases = append(sellingReleases, discogsSellerReleases{username, item.ID})
	}
	log.Printf("sellingReleases[0]=%v", sellingReleases[0])
	if _, err := db.NamedExec(`
		INSERT INTO discogs_seller_selling_releases (discogs_seller_username, discogs_release_id)
		VALUES (:discogs_seller_username, :discogs_release_id)
		ON CONFLICT DO NOTHING
	`, sellingReleases); err != nil {
		log.Printf("sellingReleases[0]=%v", sellingReleases[0])
		return err
	}
	log.Printf("sellingReleases[0]=%v", sellingReleases[0])
	return nil
}
