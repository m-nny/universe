package brain

import (
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
)

type Artist struct {
	ID        uint
	SpotifyId spotify.ID `db:"spotify_id"`
	Name      string
	Albums    []*SpotifyAlbum
}

func newArtist(sArtist spotify.SimpleArtist) *Artist {
	return &Artist{Name: sArtist.Name, SpotifyId: sArtist.ID}
}

func upsertArtistsSqlx(db *sqlx.DB, sArtists []spotify.SimpleArtist, bi *brainIndex) ([]*Artist, error) {
	if len(sArtists) == 0 {
		return []*Artist{}, nil
	}
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
	}
	bi.AddArtists(newArtists)
	return append(existingArtists, newArtists...), nil
}

func getAllArtistsSqlx(db *sqlx.DB) ([]Artist, error) {
	var existingArtists []Artist
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
