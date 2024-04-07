package brain

import "fmt"

type DiscogsSeller struct {
	Username        string
	SellingReleases []*DiscogsRelease
}

func newDiscogsSeller(username string) *DiscogsSeller {
	return &DiscogsSeller{
		Username: username,
	}
}

func (b *Brain) addDiscogsReleases(username string, releases []*DiscogsRelease) error {
	return fmt.Errorf("not implemented")
	// seller := newDiscogsSeller(username)
	// if err := b.gormDb.Where("username = ?", username).FirstOrCreate(&seller).Error; err != nil {
	// 	return err
	// }
	// return b.gormDb.Transaction(func(tx *gorm.DB) error {
	// 	if err := tx.Model(&seller).Association("SellingReleases").Clear(); err != nil {
	// 		return err
	// 	}
	// 	if err := tx.Model(&seller).Association("SellingReleases").Append(releases); err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })
}
