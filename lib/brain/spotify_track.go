package brain

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

type SpotifyTrack struct {
	ID             uint
	SpotifyId      spotify.ID `db:"spotify_id"`
	Name           string
	SpotifyAlbumId uint `db:"spotify_album_id"`
	SpotifyAlbum   *SpotifyAlbum
	Artists        []*Artist
	MetaTrackId    uint `db:"meta_track_id"`
	MetaTrack      *MetaTrack
}

func newSpotifyTrack(sTrack spotify.SimpleTrack, bSpotifyAlbum *SpotifyAlbum, bArtists []*Artist, bMetaTrack *MetaTrack) *SpotifyTrack {
	return &SpotifyTrack{
		Name:           sTrack.Name,
		SpotifyId:      sTrack.ID,
		SpotifyAlbumId: bSpotifyAlbum.ID,
		SpotifyAlbum:   bSpotifyAlbum,
		Artists:        bArtists,
		MetaTrackId:    bMetaTrack.ID,
		MetaTrack:      bMetaTrack,
	}
}

func upsertSpotifyTracksSqlx(db *sqlx.DB, sTracks []spotify.SimpleTrack, bi *brainIndex) ([]*SpotifyTrack, error) {
	if len(sTracks) == 0 {
		return []*SpotifyTrack{}, nil
	}
	var spotifyIds []spotify.ID
	for _, sTrack := range sTracks {
		spotifyIds = append(spotifyIds, sTrack.ID)
	}
	query, args, err := sqlx.In(`SELECT * FROM spotify_tracks WHERE spotify_id IN (?)`, spotifyIds)
	if err != nil {
		return nil, err
	}
	query = db.Rebind(query)
	var existingTracks []*SpotifyTrack
	if err := db.Select(&existingTracks, query, args...); err != nil {
		return nil, err
	}
	bi.AddSpotifyTracks(existingTracks)

	var newTracks []*SpotifyTrack
	for _, sTrack := range sTracks {
		if _, ok := bi.GetSpotifyTrack(sTrack); ok {
			continue
		}
		bArtists, ok := bi.GetArtists(sTrack.Artists)
		if !ok {
			return nil, fmt.Errorf("could not find artist for %s, but it should be there", sTrack.Name)
		}
		bSpotifyAlbum, ok := bi.GetSpotifyAlbum(sTrack.Album)
		if !ok {
			return nil, fmt.Errorf("could not find spotify album for %s, but it should be there", sTrack.Name)
		}
		bMetaTrack, ok := bi.GetMetaTrack(sTrack)
		if !ok {
			return nil, fmt.Errorf("could not find meta track for %s, but it should be there", sTrack.Name)
		}
		newTracks = append(newTracks, newSpotifyTrack(sTrack, bSpotifyAlbum, bArtists, bMetaTrack))
	}
	if len(newTracks) == 0 {
		return existingTracks, nil
	}
	rows, err := db.NamedQuery(`INSERT INTO spotify_tracks (spotify_id, name, spotify_album_id, meta_track_id) VALUES (:spotify_id, :name, :spotify_album_id, :meta_track_id) RETURNING id`, newTracks)
	if err != nil {
		return nil, err
	}
	for idx := 0; rows.Next(); idx++ {
		if err := rows.Scan(&newTracks[idx].ID); err != nil {
			return nil, err
		}
	}
	bi.AddSpotifyTracks(newTracks)
	var spotifyTrackArtsits []map[string]any
	for _, bMetaTrack := range newTracks {
		for _, bArtist := range bMetaTrack.Artists {
			spotifyTrackArtsits = append(spotifyTrackArtsits, map[string]any{
				"spotify_track_id": bMetaTrack.ID,
				"artist_id":        bArtist.ID,
			})
		}
	}
	_, err = db.NamedExec(`INSERT INTO spotify_track_artists (spotify_track_id, artist_id) VALUES (:spotify_track_id, :artist_id)`, spotifyTrackArtsits)
	if err != nil {
		return nil, err
	}
	return append(existingTracks, newTracks...), nil
}
