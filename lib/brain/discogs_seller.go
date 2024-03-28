package brain

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
	seller := DiscogsSeller{Username: username}
	if err := b.gormDb.Where("username = ?", username).FirstOrCreate(&seller).Error; err != nil {
		return err
	}
	if err := b.gormDb.Model(&seller).Association("SellingReleases").Replace(releases); err != nil {
		return err
	}
	return nil
}
