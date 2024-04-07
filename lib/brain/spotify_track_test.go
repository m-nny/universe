package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

// Expected outputs
var (
	bSpotifyTrackOS = &SpotifyTrack{
		Artists:        []*Artist{bArtistLP},
		ID:             1,
		MetaTrackId:    1,
		Name:           "One Step Closer",
		SpotifyAlbumId: 1,
		SpotifyId:      "spotify:one_step",
	}
	bSpotifyTrackITE = &SpotifyTrack{
		Artists:        []*Artist{bArtistLP},
		ID:             2,
		MetaTrackId:    2,
		Name:           "In the end",
		SpotifyAlbumId: 1,
		SpotifyId:      "spotify:in_the_end",
	}
	bSpotifyTrackSC = &SpotifyTrack{
		Artists:        []*Artist{bArtistPR},
		ID:             3,
		MetaTrackId:    3,
		Name:           "Something Comforting",
		SpotifyAlbumId: 2,
		SpotifyId:      "spotify:something_comforting",
	}
)

func Test_upsertSpotifyTracksSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNSpotifyTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyTrack{bSpotifyTrackOS, bSpotifyTrackITE}
		got1, err := upsertSpotifyTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want1, got1); diff != "" {
			t.Errorf("upsertSpotifyTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertSpotifyTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want1, got2); diff != "" {
			t.Errorf("upsertSpotifyTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNSpotifyTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 2, checkNSpotifyTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNSpotifyTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE, sSimpleTrackSC}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyTrack{bSpotifyTrackOS}
		got1, err := upsertSpotifyTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want1, got1); diff != "" {
			t.Errorf("upsertSpotifyTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*SpotifyTrack{bSpotifyTrackITE}
		got2, err := upsertSpotifyTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want2, got2); diff != "" {
			t.Errorf("upsertSpotifyTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*SpotifyTrack{bSpotifyTrackSC}
		got3, err := upsertSpotifyTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackSC}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want3, got3); diff != "" {
			t.Errorf("upsertSpotifyTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 3, checkNSpotifyTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 3, checkNSpotifyTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}
	})
	t.Run("handles 0 args", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNSpotifyTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}

		want := []*SpotifyTrack{}
		got, err := upsertSpotifyTracksSqlx(sqlxDb, []spotify.SimpleTrack{}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want, got); diff != "" {
			t.Errorf("upsertSpotifyTracksSqlx() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 0, checkNSpotifyTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}
	})
}

var IGNORE_SPOTIFY_TRACK_FIELDS = cmpopts.IgnoreFields(SpotifyTrack{}, "Artists", "MetaTrack", "SpotifyAlbum")

func diffSpotifyTracks(want, got []*SpotifyTrack) string {
	return cmp.Diff(want, got, IGNORE_SPOTIFY_TRACK_FIELDS)
}

func checkNSpotifyTracksSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM spotify_tracks`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}

func checkNSpotifyTrackArtistsSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM spotify_track_artists`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}
