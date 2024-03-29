package brain

import (
	"fmt"
	"strings"

	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type DiscogsRelease struct {
	ID          uint `gorm:"primarykey"`
	DiscogsID   int
	Name        string
	ArtistName  string
	MetaAlbumId *uint
	MetaAlbum   *MetaAlbum
}

func newDiscogsRelease(release discogs.ListingRelease) *DiscogsRelease {
	r := &DiscogsRelease{
		DiscogsID:  release.ID,
		Name:       release.Title,
		ArtistName: release.Artist,
	}
	return r
}

func (dr *DiscogsRelease) addMetaAlbum(bMetaAlbum *MetaAlbum) {
	if bMetaAlbum == nil {
		return
	}
	dr.MetaAlbum = bMetaAlbum
	dr.MetaAlbumId = &bMetaAlbum.ID
}

func (b *Brain) SaveDiscorgsReleases(dReleases []discogs.ListingRelease) ([]*DiscogsRelease, error) {
	dReleases = sliceutils.Unique(dReleases, func(item discogs.ListingRelease) int { return item.ID })
	var discogsIds []int
	for _, dRelease := range dReleases {
		discogsIds = append(discogsIds, dRelease.ID)
	}

	var existingReleases []*DiscogsRelease
	if err := b.gormDb.
		Where("discogs_id IN ?", discogsIds).
		Find(&existingReleases).Error; err != nil {
		return nil, err
	}
	releaseMap := sliceutils.ToMap(existingReleases, func(item *DiscogsRelease) int { return item.DiscogsID })

	var newReleases []*DiscogsRelease
	for _, dRelease := range dReleases {
		if _, ok := releaseMap[dRelease.ID]; ok {
			continue
		}
		newReleases = append(newReleases, newDiscogsRelease(dRelease))
	}
	if len(newReleases) == 0 {
		return existingReleases, nil
	}
	if err := b.gormDb.Create(newReleases).Error; err != nil {
		return nil, err
	}
	return append(existingReleases, newReleases...), nil
}

func (b *Brain) AssociateDiscogsRelease(bRelease *DiscogsRelease, bMetaAlbum *MetaAlbum) error {
	if bMetaAlbum == nil {
		return fmt.Errorf("cannot set MetaAlbum to nil")
	}
	if err := b.gormDb.Model(bRelease).UpdateColumn("meta_album_id", bMetaAlbum.ID).Error; err != nil {
		return err
	}
	bRelease.addMetaAlbum(bMetaAlbum)
	return nil
}

func MostSimilarAlbum(dRelease *DiscogsRelease, bAlbums []*MetaAlbum) (*MetaAlbum, int, error) {
	var result *MetaAlbum
	maxScore := 0
	for _, bAlbum := range bAlbums {
		if len(bAlbum.Artists) == 0 {
			return nil, 0, fmt.Errorf("Albums should have artists populated")
		}
		score := albumSimilarity(dRelease, bAlbum)
		if score > maxScore {
			maxScore = score
			result = bAlbum
		}
	}
	return result, maxScore, nil
}

func albumSimilarity(dRelease *DiscogsRelease, eAlbum *MetaAlbum) int {
	artistScore := sliceutils.Sum(eAlbum.Artists, func(e *Artist) int { return similaryScore(dRelease.ArtistName, e.Name) })
	if artistScore == 0 {
		return 0
	}
	titleScore := similaryScore(dRelease.Name, eAlbum.AnyName)
	if titleScore == 0 {
		return 0
	}
	return artistScore + titleScore
}

func similaryScore(a, b string) int {
	a = strings.ToLower(a)
	b = strings.ToLower(b)
	score := 0
	if strings.Contains(a, b) {
		score += len(b)
	}
	if strings.Contains(b, a) {
		score += len(a)
	}
	return score
}
