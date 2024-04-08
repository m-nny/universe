package brain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

type SpotifyAlbum struct {
	ID          uint
	SpotifyId   spotify.ID `db:"spotify_id"`
	Name        string
	Artists     []*Artist
	MetaAlbumId uint `db:"meta_album_id"`
	MetaAlbum   *MetaAlbum
}

func (s *SpotifyAlbum) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.Name
}

func newSpotifyAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist, bMetaAlbum *MetaAlbum) *SpotifyAlbum {
	return &SpotifyAlbum{
		Name:        sAlbum.Name,
		SpotifyId:   sAlbum.ID,
		Artists:     bArtists,
		MetaAlbumId: bMetaAlbum.ID,
		MetaAlbum:   bMetaAlbum,
	}
}

type SpotifyAlbumArtist struct {
	SpotifyAlbumId uint       `db:"spotify_album_id"`
	ArtistId       spotify.ID `db:"artist_id"`
}

func upsertSpotifyAlbumsSqlx(db *sqlx.DB, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*SpotifyAlbum, error) {
	if len(sAlbums) == 0 {
		return []*SpotifyAlbum{}, nil
	}
	var albumSIds []spotify.ID
	for _, sAlbum := range sAlbums {
		albumSIds = append(albumSIds, sAlbum.ID)
	}
	query, args, err := sqlx.In(`SELECT * FROM spotify_albums WHERE spotify_id IN (?)`, albumSIds)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	var existingSpotifyAlbums []*SpotifyAlbum
	if err := db.Select(&existingSpotifyAlbums, query, args...); err != nil {
		return nil, err
	}
	bi.AddSpotifyAlbums(existingSpotifyAlbums)

	var newAlbums []*SpotifyAlbum
	for _, sAlbum := range sAlbums {
		if _, ok := bi.GetSpotifyAlbum(sAlbum); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sAlbum.Artists)
		if !ok {
			return nil, fmt.Errorf("bArtist not found")
		}
		bMetaAlbum, ok := bi.GetMetaAlbum(sAlbum)
		if !ok {
			return nil, fmt.Errorf("bMetaAlbum not found")
		}
		newAlbums = append(newAlbums, newSpotifyAlbum(sAlbum, bArtists, bMetaAlbum))
	}
	if len(newAlbums) == 0 {
		return existingSpotifyAlbums, nil
	}
	rows, err := db.NamedQuery(`INSERT INTO spotify_albums (spotify_id, name, meta_album_id) VALUES (:spotify_id, :name, :meta_album_id) RETURNING id`, newAlbums)
	if err != nil {
		return nil, err
	}
	for idx := 0; rows.Next(); idx++ {
		if err := rows.Scan(&newAlbums[idx].ID); err != nil {
			return nil, err
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	bi.AddSpotifyAlbums(newAlbums)

	var spotifyAlbumArtists []SpotifyAlbumArtist
	for _, bAlbum := range newAlbums {
		for _, bArtist := range bAlbum.Artists {
			spotifyAlbumArtists = append(spotifyAlbumArtists, SpotifyAlbumArtist{bAlbum.ID, bArtist.SpotifyId})
		}
	}
	_, err = db.NamedExec(`INSERT INTO spotify_album_artists (spotify_album_id, artist_id) VALUES (:spotify_album_id, :artist_id)`, spotifyAlbumArtists)
	if err != nil {
		return nil, err
	}
	return append(existingSpotifyAlbums, newAlbums...), nil
}
