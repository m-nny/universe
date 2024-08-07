package brain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"

	"github.com/m-nny/universe/lib/utils/sliceutils"
	utils "github.com/m-nny/universe/lib/utils/spotifyutils"
)

type MetaTrack struct {
	SimplifiedName string `db:"simplified_name"`
	MetaAlbumID    string `db:"meta_album_id"`
	MetaAlbum      *MetaAlbum
	Artists        []*Artist
}

func (s *MetaTrack) String() string {
	if s == nil {
		return "<nil>"
	}
	return s.SimplifiedName
}

func newMetaTrack(sTrack spotify.SimpleTrack, bMetaAlbum *MetaAlbum, bArtists []*Artist) *MetaTrack {
	return &MetaTrack{
		SimplifiedName: utils.SimplifiedTrackName(sTrack, bMetaAlbum.SimplifiedName),
		MetaAlbumID:    bMetaAlbum.SimplifiedName,
		MetaAlbum:      bMetaAlbum,
		Artists:        bArtists,
	}
}

func upsertMetaTracksSqlx(db *sqlx.DB, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*MetaTrack, error) {
	if len(sTracks) == 0 {
		return []*MetaTrack{}, nil
	}
	sTracks = sliceutils.Unique(sTracks, bi.MustTrackSimplifiedName)
	var simpNames []string
	for _, sTrack := range sTracks {
		simpName, ok := bi.TrackSimplifiedName(sTrack)
		if !ok {
			return nil, fmt.Errorf("could not get simplified name for %s, but it should be there", sTrack.Name)
		}
		simpNames = append(simpNames, simpName)
	}

	var existingMetaTracks []*MetaTrack
	query, args, err := sqlx.In(`SELECT * FROM meta_tracks WHERE simplified_name IN (?)`, simpNames)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	if err := db.Select(&existingMetaTracks, query, args...); err != nil {
		return nil, err
	}
	bi.AddMetaTracks(existingMetaTracks)

	var newTracks []*MetaTrack
	for _, sTrack := range sTracks {
		if _, ok := bi.GetMetaTrack(sTrack); ok {
			continue
		}
		bMetaAlbum, ok := bi.GetMetaAlbum(sTrack.Album)
		if !ok {
			return nil, fmt.Errorf("could not find meta album for %s, but it should be there", sTrack.Name)
		}
		bArtists, ok := bi.GetArtists(sTrack.Artists)
		if !ok {
			return nil, fmt.Errorf("could not find meta album for %s, but it should be there", sTrack.Name)
		}
		newTracks = append(newTracks, newMetaTrack(sTrack, bMetaAlbum, bArtists))
	}
	if len(newTracks) == 0 {
		return existingMetaTracks, nil
	}
	if _, err := db.NamedExec(`
		INSERT INTO meta_tracks (simplified_name, meta_album_id)
		VALUES (:simplified_name, :meta_album_id)`, newTracks); err != nil {
		return nil, err
	}
	bi.AddMetaTracks(newTracks)
	type metaTrackArtistIds struct {
		MetaTrackId string     `db:"meta_track_id"`
		ArtistId    spotify.ID `db:"artist_id"`
	}
	var metaTrackArtsitIds []metaTrackArtistIds
	for _, bMetaTrack := range newTracks {
		for _, bArtist := range bMetaTrack.Artists {
			metaTrackArtsitIds = append(metaTrackArtsitIds, metaTrackArtistIds{
				MetaTrackId: bMetaTrack.SimplifiedName,
				ArtistId:    bArtist.SpotifyId,
			})
		}
	}
	_, err = db.NamedExec(`
		INSERT INTO meta_track_artists (meta_track_id, artist_id)
		VALUES (:meta_track_id, :artist_id)`, metaTrackArtsitIds)
	if err != nil {
		return nil, err
	}
	return append(existingMetaTracks, newTracks...), nil
}
