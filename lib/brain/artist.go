package brain

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Artist struct {
	ID        uint       `gorm:"primarykey"`
	SpotifyId spotify.ID `db:"spotify_id"`
	Name      string
	Albums    []*SpotifyAlbum `gorm:"many2many:spotify_album_artists;"`
}

func newArtist(sArtist spotify.SimpleArtist) *Artist {
	return &Artist{Name: sArtist.Name, SpotifyId: sArtist.ID}
}

// SaveArtists takes list of spotify Artists and returns Brain representain of them.
//   - It will create new entries in DB if necessary
//   - It will deduplicate returned artists, this may result in len(result) < len(sArtists)
func (b *Brain) SaveArtists(sArtists []*spotify.FullArtist) ([]*Artist, error) {
	sSimpleArtists := sliceutils.Map(sArtists, func(item *spotify.FullArtist) spotify.SimpleArtist { return item.SimpleArtist })
	return upsertArtists(b, sSimpleArtists, newBrainIndex())
}

func upsertArtists(b *Brain, sArtists []spotify.SimpleArtist, bi *brainIndex) ([]*Artist, error) {
	sqlxBi := bi.Clone()
	sqlxArtists, err := _sqlxUpsertArtists(b.sqlxDb, sArtists, sqlxBi)
	if err != nil {
		return nil, err
	}
	gormArtists, err := _gormUpsertArtists(b.gormDb, sArtists, bi)
	if err != nil {
		return nil, err
	}
	// TODO(m-nny): check sqlxArtists == gormArtists
	if len(gormArtists) != len(sqlxArtists) {
		return nil, fmt.Errorf("len(gormArtists) != len(sqlxArtists): %d != %d", len(gormArtists), len(sqlxArtists))
	}
	return gormArtists, nil
}
func _gormUpsertArtists(db *gorm.DB, sArtists []spotify.SimpleArtist, bi *brainIndex) ([]*Artist, error) {
	sArtists = sliceutils.Unique(sArtists, func(item spotify.SimpleArtist) spotify.ID { return item.ID })
	var spotifyIds []spotify.ID
	for _, sArtist := range sArtists {
		spotifyIds = append(spotifyIds, sArtist.ID)
	}

	var existingArtists []*Artist
	if err := db.
		Where("spotify_id IN ?", spotifyIds).
		Find(&existingArtists).Error; err != nil {
		return nil, err
	}
	bi.AddArtists(existingArtists)

	var newArtists []*Artist
	for _, sArtist := range sArtists {
		if _, ok := bi.GetArtist(sArtist); ok {
			continue
		}
		newArtists = append(newArtists, newArtist(sArtist))
	}
	// All artists are already created, can exit
	if len(newArtists) == 0 {
		return existingArtists, nil
	}
	if err := db.Create(newArtists).Error; err != nil {
		return nil, err
	}
	bi.AddArtists(newArtists)
	return append(existingArtists, newArtists...), nil
}
func _sqlxUpsertArtists(db *sqlx.DB, sArtists []spotify.SimpleArtist, bi *brainIndex) ([]*Artist, error) {
	sArtists = sliceutils.Unique(sArtists, func(item spotify.SimpleArtist) spotify.ID { return item.ID })
	var spotifyIds []spotify.ID
	for _, sArtist := range sArtists {
		spotifyIds = append(spotifyIds, sArtist.ID)
	}

	var existingArtists []*Artist
	query, args, err := sqlx.In(`SELECT * FROM artists WHERE spotify_id IN (?)`, spotifyIds)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	if err := db.Select(&existingArtists, query, args...); err != nil {
		return nil, err
	}
	bi.AddArtists(existingArtists)

	var newArtists []*Artist
	for _, sArtist := range sArtists {
		if _, ok := bi.GetArtist(sArtist); ok {
			continue
		}
		newArtists = append(newArtists, newArtist(sArtist))
	}
	// All artists are already created, can exit
	if len(newArtists) == 0 {
		return existingArtists, nil
	}
	rows, err := db.NamedQuery(`INSERT INTO artists (spotify_id, name) VALUES (:spotify_id, :name) RETURNING id`, newArtists)
	if err != nil {
		return nil, err
	}
	for idx := 0; rows.Next(); idx++ {
		if err := rows.Scan(&newArtists[idx].ID); err != nil {
			return nil, err
		}
		log.Printf("row: %+v", newArtists[idx])
	}
	bi.AddArtists(newArtists)
	return append(existingArtists, newArtists...), nil
}

func getAllSqlxArtists(db *sqlx.DB) ([]*Artist, error) {
	var existingArtists []*Artist
	query, args, err := sqlx.In(`SELECT * FROM artists`)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	if err := db.Select(&existingArtists, query, args...); err != nil {
		return nil, err
	}
	return existingArtists, nil
}
