package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

// inputs
var (
	sSimpleTrackOS = spotify.SimpleTrack{
		ID:          "spotify:one_step",
		Name:        "One Step Closer",
		TrackNumber: 2,
		Artists:     []spotify.SimpleArtist{sArtistLP},
		Album:       sSimpleAlbumHT,
	}
	sSimpleTrackITE = spotify.SimpleTrack{
		ID:          "spotify:in_the_end",
		Name:        "In the end",
		TrackNumber: 8,
		Artists:     []spotify.SimpleArtist{sArtistLP},
		Album:       sSimpleAlbumHT,
	}
	sSimpleTrackSC = spotify.SimpleTrack{
		ID:          "spotify:something_comforting",
		Name:        "Something Comforting",
		TrackNumber: 11,
		Artists:     []spotify.SimpleArtist{sArtistPR},
		Album:       sSimpleAlbumN,
	}
	sSavedTrackOS = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: sSimpleTrackOS,
			Album:       sSimpleAlbumHT,
		},
	}
	sSavedTrackITE = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: sSimpleTrackITE,
			Album:       sSimpleAlbumHT,
		},
	}
	sSavedTrackSC = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: sSimpleTrackSC,
			Album:       sSimpleAlbumN,
		},
	}
)

// outputs
var (
	bMetaTrackOS = &MetaTrack{
		ID:             1,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: bMetaAlbumHT.SimplifiedName + " - 02. one step closer",
		MetaAlbumID:    bMetaAlbumHT.ID,
		MetaAlbum:      bMetaAlbumHT,
	}
	bMetaTrackITE = &MetaTrack{
		ID:             2,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: bMetaAlbumHT.SimplifiedName + " - 08. in the end",
		MetaAlbumID:    bMetaAlbumHT.ID,
		MetaAlbum:      bMetaAlbumHT,
	}
	bMetaTrackSC = &MetaTrack{
		ID:             3,
		Artists:        []*Artist{bArtistPR},
		SimplifiedName: bMetaAlbumN.SimplifiedName + " - 11. something comforting",
		MetaAlbumID:    bMetaAlbumN.ID,
		MetaAlbum:      bMetaAlbumN,
	}
)

func Test_upsertMetaTracksSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNMetaTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNMetaTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaTrack{bMetaTrackOS, bMetaTrackITE}
		got1, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("upsertMetaTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got2); diff != "" {
			t.Errorf("upsertMetaTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNMetaTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 2, checkNMetaTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNMetaTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNMetaTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaTrack{bMetaTrackOS}
		got1, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackOS}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("upsertMetaTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bMetaTrackITE}
		got2, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("upsertMetaTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaTrack{bMetaTrackSC}
		got3, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{sSimpleTrackSC}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want3, got3); diff != "" {
			t.Errorf("upsertMetaTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 3, checkNMetaTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 3, checkNMetaTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}
	})
	t.Run("handles 0 args", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNMetaTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNMetaTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}

		want := []*MetaTrack{}
		got, err := upsertMetaTracksSqlx(sqlxDb, []spotify.SimpleTrack{}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want, got); diff != "" {
			t.Errorf("upsertMetaTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 0, checkNMetaTracksSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
		if want, got := 0, checkNMetaTrackArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta track artists, but want %d rows", got, want)
		}
	})
}

var IGNORE_META_TRACK_FIELDS = cmpopts.IgnoreFields(MetaTrack{}, "MetaAlbum", "Artists")

func diffMetaTracks(want, got []*MetaTrack) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS, IGNORE_META_TRACK_FIELDS)
}

func checkNMetaTracksSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM meta_tracks`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}

func checkNMetaTrackArtistsSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM meta_track_artists`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}
