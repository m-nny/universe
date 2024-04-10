package brain

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type DiscogsRelease struct {
	ArtistName string `db:"artist_name"`
	DiscogsID  int    `db:"discogs_id"`
	Format     string
	Name       string

	MetaAlbumId              *string `db:"meta_album_id"`
	MetaAlbumSimilariryScore int     `db:"meta_album_score"`
	SearchedMetaAlbum        bool    `db:"searched_meta_album"`

	MetaAlbum *MetaAlbum
}

func newDiscogsRelease(release discogs.ListingRelease) *DiscogsRelease {
	r := &DiscogsRelease{
		DiscogsID:  release.ID,
		Name:       release.Title,
		ArtistName: release.Artist,
		Format:     release.Format,

		SearchedMetaAlbum: false,
	}
	return r
}

func (b *Brain) SaveDiscorgsReleases(dReleases []discogs.ListingRelease, username string) ([]*DiscogsRelease, error) {
	bReleases, err := upsertDiscogsReleases(b.sqlxDb, dReleases)
	if err != nil {
		return nil, err
	}
	if err := addDiscogsReleases(b.sqlxDb, username, bReleases); err != nil {
		return nil, err
	}
	return bReleases, nil
}

func upsertDiscogsReleases(db *sqlx.DB, dReleases []discogs.ListingRelease) ([]*DiscogsRelease, error) {
	if len(dReleases) == 0 {
		return []*DiscogsRelease{}, nil
	}
	dReleases = sliceutils.Unique(dReleases, func(item discogs.ListingRelease) int { return item.ID })
	var discogsIds []int
	for _, dRelease := range dReleases {
		discogsIds = append(discogsIds, dRelease.ID)
	}

	query, args, err := sqlx.In(`SELECT * FROM discogs_releases WHERE discogs_id IN (?)`, discogsIds)
	if err != nil {
		return nil, err
	}
	var existingReleases []*DiscogsRelease
	if err := db.Select(&existingReleases, query, args...); err != nil {
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
	if _, err := db.NamedExec(`
		INSERT INTO discogs_releases (discogs_id, name, artist_name, format)
		VALUES (:discogs_id, :name, :artist_name, :format)`, newReleases); err != nil {
		return nil, err
	}
	return append(existingReleases, newReleases...), nil
}

func (b *Brain) AssociateDiscogsRelease(bRelease *DiscogsRelease, bMetaAlbum *MetaAlbum, score int) error {
	bRelease.MetaAlbum = bMetaAlbum
	if bMetaAlbum != nil {
		bRelease.MetaAlbumId = &bMetaAlbum.SimplifiedName
	}
	bRelease.MetaAlbumSimilariryScore = score
	bRelease.SearchedMetaAlbum = true
	_, err := b.sqlxDb.NamedExec(`
		UPDATE discogs_releases 
		SET  meta_album_id=:meta_album_id, meta_album_score=:meta_album_score, searched_meta_album=:searched_meta_album
		WHERE id=:id
	`, bRelease)
	if err != nil {
		return err
	}
	return nil
}

func MostSimilarAlbum(dRelease *DiscogsRelease, bAlbums []*MetaAlbum) (*MetaAlbum, int, error) {
	var result *MetaAlbum
	maxScore := 0
	for _, bAlbum := range bAlbums {
		if len(bAlbum.Artists) == 0 {
			return nil, 0, fmt.Errorf("albums should have artists populated")
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
	titleScore := similaryScore(dRelease.Name, eAlbum.AnyName)
	if titleScore == 0 {
		return 0
	}
	artistScore := sliceutils.Sum(eAlbum.Artists,
		func(e *Artist) int { return similaryScore(dRelease.ArtistName, e.Name) })
	if artistScore == 0 {
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
