package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

// Inputs
var (
	sSimpleAlbumHT = spotify.SimpleAlbum{
		ID:      "spotify:hybrid_theory",
		Name:    "Hybrid Theory",
		Artists: []spotify.SimpleArtist{sArtistLP},
	}
	sSimpleAlbumHT20 = spotify.SimpleAlbum{
		ID:      "spotify:hybryd_theory_20",
		Name:    "Hybrid Theory (20th Anniversary Edition)",
		Artists: []spotify.SimpleArtist{sArtistLP},
	}
	sSimpleAlbumN = spotify.SimpleAlbum{
		ID:      "spotify:nurture",
		Name:    "Nurture",
		Artists: []spotify.SimpleArtist{sArtistPR},
	}
	sFullAlbumHT   = &spotify.FullAlbum{SimpleAlbum: sSimpleAlbumHT}
	sFullAlbumHT20 = &spotify.FullAlbum{SimpleAlbum: sSimpleAlbumHT20}
	sFullAlbumN    = &spotify.FullAlbum{SimpleAlbum: sSimpleAlbumN}
)

// Expected outputs
var (
	bMetaAlbumHT = &MetaAlbum{
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: "linkin park - hybrid theory",
	}
	bMetaAlbumN = &MetaAlbum{
		Artists:        []*Artist{bArtistPR},
		SimplifiedName: "porter robinson - nurture",
	}
)

func Test_upsertMetaAlbumsSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		if want, got := 0, checkNMetaAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}
		if want, got := 0, checkNMetaAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta album artists, but want %d rows", got, want)
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bMetaAlbumN}
		got1, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Fatalf("upsertMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Fatalf("upsertMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 2, checkNMetaAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}
		if want, got := 2, checkNMetaAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta album artists, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		if want, got := 0, checkNMetaAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d rows, but want %d rows", got, want)
		}
		if want, got := 0, checkNMetaAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta album artists, but want %d rows", got, want)
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Fatalf("checkNMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT20}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want2, got2); diff != "" {
			t.Fatalf("checkNMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bMetaAlbumN}
		got3, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want3, got3); diff != "" {
			t.Fatalf("checkNMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 2, checkNMetaAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}
		if want, got := 2, checkNMetaAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta album artists, but want %d rows", got, want)
		}
	})
	t.Run("handles 0 args", func(t *testing.T) {
		brain := getInmemoryBrain(t)

		want1 := []*MetaAlbum{}
		got1, err := upsertMetaAlbumsSqlx(brain.sqlxDb, []spotify.SimpleAlbum{}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("upsertMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 0, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}
	})
}

var IGNORE_META_ALBUM_FIELDS = cmpopts.IgnoreFields(MetaAlbum{}, "AnyName", "Artists")

func diffMetaAlbums(want, got []*MetaAlbum) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS)
}

func checkNMetaAlbumsSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM meta_albums`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}

func checkNMetaAlbumArtistsSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM meta_album_artists`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}
