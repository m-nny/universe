package brain

import "gorm.io/gorm"

type DiscogsSeller struct {
	Username        string            `gorm:"primarykey"`
	SellingReleases []*DiscogsRelease `gorm:"many2many:discogs_seller_selling_releases"`
}

func newDiscogsSeller(username string) *User {
	return &User{
		Username: username,
	}
}

func (b *Brain) upsertDiscogsUser(username string, releases []*DiscogsRelease) error {
	seller := newDiscogsSeller(username)
	if err := b.gormDb.Where("username = ?", username).FirstOrCreate(&seller).Error; err != nil {
		return err
	}
	return b.gormDb.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&seller).Association("SellingReleases").Clear(); err != nil {
			return err
		}
		if err := tx.Model(&seller).Association("SellingReleases").Append(releases); err != nil {
			return err
		}
		return nil
	})
}
