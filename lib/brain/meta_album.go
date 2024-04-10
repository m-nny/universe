package brain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
	"github.com/m-nny/universe/lib/utils/spotifyutils"
)

type MetaAlbum struct {
	SimplifiedName string `db:"simplified_name"`
	AnyName        string `db:"any_name"`
	Artists        []*Artist
}

func newMetaAlbum(sAlbum spotify.SimpleAlbum, bArtists []*Artist) *MetaAlbum {
	return &MetaAlbum{
		SimplifiedName: spotifyutils.SimplifiedAlbumName(sAlbum),
		AnyName:        sAlbum.Name,
		Artists:        bArtists,
	}
}

func upsertMetaAlbumsSqlx(db *sqlx.DB, sAlbums []spotify.SimpleAlbum, bi *brainIndex) ([]*MetaAlbum, error) {
	if len(sAlbums) == 0 {
		return []*MetaAlbum{}, nil
	}
	sAlbums = sliceutils.Unique(sAlbums, spotifyutils.SimplifiedAlbumName)
	var simpNames []string
	for _, sAlbum := range sAlbums {
		simpNames = append(simpNames, spotifyutils.SimplifiedAlbumName(sAlbum))
	}

	query, args, err := sqlx.In(`SELECT * FROM meta_albums WHERE simplified_name IN (?)`, simpNames)
	if err != nil {
		return nil, fmt.Errorf("could not prepare query for existing meta albums: %w", err)
	}
	query = db.Rebind(query)
	var existingMetaAlbums []*MetaAlbum
	if err := db.Select(&existingMetaAlbums, query, args...); err != nil {
		return nil, fmt.Errorf("could not get existing meta albums: %w", err)
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
	if _, err = db.NamedExec(`
		INSERT INTO meta_albums (simplified_name, any_name)
		VALUES (:simplified_name, :any_name)`, newAlbums); err != nil {
		return nil, fmt.Errorf("could not insert meta albums: %w", err)
	}
	bi.AddMetaAlbums(newAlbums)
	type metaAlbumArtistIds struct {
		MetaAlbumId string     `db:"meta_album_id"`
		ArtistId    spotify.ID `db:"artist_id"`
	}
	var metaAlbumArtistsIds []metaAlbumArtistIds
	for _, bAlbum := range newAlbums {
		for _, bArtist := range bAlbum.Artists {
			metaAlbumArtistsIds = append(metaAlbumArtistsIds, metaAlbumArtistIds{bAlbum.SimplifiedName, bArtist.SpotifyId})
		}
	}
	_, err = db.NamedExec(`
		INSERT INTO meta_album_artists (meta_album_id, artist_id)
		VALUES (:meta_album_id, :artist_id)`, metaAlbumArtistsIds)
	if err != nil {
		return nil, fmt.Errorf("could not insert meta album artists: %w", err)
	}
	return append(existingMetaAlbums, newAlbums...), nil
}
