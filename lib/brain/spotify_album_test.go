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
	bSpotifyAlbumHT = &SpotifyAlbum{
		ID:          1,
		Artists:     []*Artist{bArtistLP},
		SpotifyId:   "spotify:hybrid_theory",
		Name:        "Hybrid Theory",
		MetaAlbumId: 1,
	}
	bSpotifyAlbumHT20 = &SpotifyAlbum{
		ID:          2,
		Artists:     []*Artist{bArtistLP},
		SpotifyId:   "spotify:hybryd_theory_20",
		Name:        "Hybrid Theory (20th Anniversary Edition)",
		MetaAlbumId: 1,
	}
	bSpotifyAlbumN = &SpotifyAlbum{
		ID:          3,
		Artists:     []*Artist{bArtistPR},
		SpotifyId:   "spotify:nurture",
		Name:        "Nurture",
		MetaAlbumId: 2,
	}
)

func Test_upsertSpotifyAlbumsSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		if want, got := 0, checkNSpotifyAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify albums, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify album artists, but want %d rows", got, want)
		}

		// Setup Artists & MetaAlbums
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyAlbum{bSpotifyAlbumHT, bSpotifyAlbumHT20, bSpotifyAlbumN}
		got1, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want1, got1); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want1, got2); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 3, checkNSpotifyAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify albums, but want %d rows", got, want)
		}
		if want, got := 3, checkNSpotifyAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify album artists, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		if want, got := 0, checkNSpotifyAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify albums, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify album artists, but want %d rows", got, want)
		}

		// Setup Artists & MetaAlbums
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyAlbum{bSpotifyAlbumHT}
		got1, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want1, got1); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*SpotifyAlbum{bSpotifyAlbumHT20}
		got2, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT20}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want2, got2); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*SpotifyAlbum{bSpotifyAlbumN}
		got3, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want3, got3); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 3, checkNSpotifyAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify albums, but want %d rows", got, want)
		}
		if want, got := 3, checkNSpotifyAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify album artists, but want %d rows", got, want)
		}
	})
	t.Run("handles 0 args", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		if want, got := 0, checkNSpotifyAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify albums, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify album artists, but want %d rows", got, want)
		}

		want := []*SpotifyAlbum{}
		got, err := upsertSpotifyAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want, got); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 0, checkNSpotifyAlbumsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify albums, but want %d rows", got, want)
		}
		if want, got := 0, checkNSpotifyAlbumArtistsSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d spotify album artists, but want %d rows", got, want)
		}
	})
}

var IGNORE_SPOTIFY_ALBUM_FIELDS = cmpopts.IgnoreFields(SpotifyAlbum{}, "Artists", "MetaAlbum")

func diffSpotifyAlbums(want, got []*SpotifyAlbum) string {
	return cmp.Diff(want, got, IGNORE_SPOTIFY_ALBUM_FIELDS)
}

func checkNSpotifyAlbumsSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM spotify_albums`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}

func checkNSpotifyAlbumArtistsSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM spotify_album_artists`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}
