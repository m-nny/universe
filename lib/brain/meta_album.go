package brain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

type MetaAlbum struct {
	ID             uint      `gorm:"primarykey"`
	SimplifiedName string    `db:"simplified_name"`
	AnyName        string    `db:"any_name"`
	Artists        []*Artist `gorm:"many2many:meta_album_artists;"`
}

func newMetaAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist) *MetaAlbum {
	return &MetaAlbum{
		SimplifiedName: utils.SimplifiedAlbumName(sAlbum),
		AnyName:        sAlbum.Name,
		Artists:        bArtists,
	}
}

type MetaAlbumArtist struct {
	MetaAlbumId uint `db:"meta_album_id"`
	ArtistId    uint `db:"artist_id"`
}

func upsertMetaAlbums(b *Brain, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*MetaAlbum, error) {
	sqlxBi := bi.Clone()
	sqlxMetaAlbums, err := upsertMetaAlbumsSqlx(b.sqlxDb, sAlbums, sqlxBi)
	if err != nil {
		return nil, err
	}
	gormMetaAlbums, err := upsertMetaAlbumsGorm(b.gormDb, sAlbums, bi)
	if err != nil {
		return nil, err
	}
	// TODO(m-nny): check sqlxMetaAlbums == gormMetaAlbums
	if len(gormMetaAlbums) != len(sqlxMetaAlbums) {
		return nil, fmt.Errorf("len(gormMetaAlbums) != len(sqlxMetaAlbums): %d != %d", len(gormMetaAlbums), len(sqlxMetaAlbums))
	}
	return gormMetaAlbums, nil

}

func upsertMetaAlbumsGorm(db *gorm.DB, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*MetaAlbum, error) {
	sAlbums = sliceutils.Unique(sAlbums, utils.SimplifiedAlbumName)
	var simpNames []string
	for _, sAlbum := range sAlbums {
		simpNames = append(simpNames, utils.SimplifiedAlbumName(sAlbum))
	}

	var existingMetaAlbums []*MetaAlbum
	if err := db.
		Preload("Artists").
		Where("simplified_name IN ?", simpNames).
		Find(&existingMetaAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddMetaAlbums(existingMetaAlbums)

	// var newMetaAlbums []*MetaAlbum
	var newAlbums []*MetaAlbum
	for _, sAlbum := range sAlbums {
		if _, ok := bi.GetMetaAlbum(sAlbum); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sAlbum.Artists)
		if !ok {
			// log.Printf("could not get artists for %s, but they should be there", sAlbum.Name)
			return nil, fmt.Errorf("could not get artists for %s, but they should be there", sAlbum.Name)
		}
		newAlbums = append(newAlbums, newMetaAlbum(sAlbum, bArtists))
	}
	if len(newAlbums) == 0 {
		return existingMetaAlbums, nil
	}
	if err := db.Create(newAlbums).Error; err != nil {
		return nil, err
	}
	bi.AddMetaAlbums(newAlbums)
	return append(existingMetaAlbums, newAlbums...), nil
}

func upsertMetaAlbumsSqlx(db *sqlx.DB, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*MetaAlbum, error) {
	sAlbums = sliceutils.Unique(sAlbums, utils.SimplifiedAlbumName)
	var simpNames []string
	for _, sAlbum := range sAlbums {
		simpNames = append(simpNames, utils.SimplifiedAlbumName(sAlbum))
	}

	var existingMetaAlbums []*MetaAlbum
	query, args, err := sqlx.In(`SELECT * FROM meta_albums WHERE simplified_name IN (?)`, simpNames)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	if err := db.Select(&existingMetaAlbums, query, args...); err != nil {
		return nil, err
	}
	bi.AddMetaAlbums(existingMetaAlbums)

	var newAlbums []*MetaAlbum
	for _, sAlbum := range sAlbums {
		if _, ok := bi.GetMetaAlbum(sAlbum); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sAlbum.Artists)
		if !ok {
			return nil, fmt.Errorf("could not get artists for %s, but they should be there", sAlbum.Name)
		}
		newAlbums = append(newAlbums, newMetaAlbum(sAlbum, bArtists))
	}
	if len(newAlbums) == 0 {
		return existingMetaAlbums, nil
	}
	rows, err := db.NamedQuery(`INSERT INTO meta_albums (simplified_name, any_name) VALUES (:simplified_name, :any_name) RETURNING id`, newAlbums)
	if err != nil {
		return nil, err
	}
	for idx := 0; rows.Next(); idx++ {
		if err := rows.Scan(&newAlbums[idx].ID); err != nil {
			return nil, err
		}
	}
	bi.AddMetaAlbums(newAlbums)
	var metaAlbumArtists []MetaAlbumArtist
	for _, bAlbum := range newAlbums {
		for _, bArtist := range bAlbum.Artists {
			metaAlbumArtists = append(metaAlbumArtists, MetaAlbumArtist{bAlbum.ID, bArtist.ID})
		}
	}
	_, err = db.NamedExec(`INSERT INTO meta_album_artists (meta_album_id, artist_id) VALUES (:meta_album_id, :artist_id)`, metaAlbumArtists)
	if err != nil {
		return nil, err
	}
	return append(existingMetaAlbums, newAlbums...), nil
}

// SaveAlbums returns Brain representain of a spotify album
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned albums, this may result in len(result) < len(sAlbums)
//   - NOTE: Does not debupe based on simplified name
func (b *Brain) SaveAlbums(fullAlbums []*spotify.FullAlbum) ([]*MetaAlbum, error) {
	var sAlbums []spotify.SimpleAlbum
	var sTracks []spotify.SimpleTrack
	for _, sAlbum := range fullAlbums {
		// sAlbum.SimpleAlbum.Artists = sAlbum.Artists
		sAlbums = append(sAlbums, sAlbum.SimpleAlbum)
		sTracks = append(sTracks, sAlbum.Tracks.Tracks...)
	}
	albums, _, err := b.batchSaveAlbumTracks(sAlbums, sTracks)
	return albums, err
}

func (b *Brain) SaveSimpleAlbums(sAlbums []spotify.SimpleAlbum) ([]*MetaAlbum, error) {
	albums, _, err := b.batchSaveAlbumTracks(sAlbums, nil)
	return albums, err
}

func (b *Brain) MetaAlbumCount() (int, error) {
	var count int64
	err := b.gormDb.Model(&MetaAlbum{}).Count(&count).Error
	return int(count), err
}

func getAllMetaAlbumsSqlx(db *sqlx.DB) ([]MetaAlbum, error) {
	var bMetaAlbums []MetaAlbum
	if err := db.Select(&bMetaAlbums, `SELECT * FROM meta_albums`); err != nil {
		return nil, err
	}
	return bMetaAlbums, nil
}

func cntMetaAlbumArtistsSqlx(db *sqlx.DB) (int, error) {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM meta_album_artists`); err != nil {
		return 0, err
	}
	return cnt, nil
}

func getAllMetaAlbumsGorm(db *gorm.DB) ([]MetaAlbum, error) {
	var gormMetaAlbums []MetaAlbum
	if err := db.Find(&gormMetaAlbums).Error; err != nil {
		return nil, err
	}
	return gormMetaAlbums, nil
}
