package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

var (
	sArtistLP = spotify.SimpleArtist{
		ID:   "spotify:linkin_park",
		Name: "Linkin Park",
	}
	sArtistPR = spotify.SimpleArtist{
		ID:   "spotify:porter_robinson",
		Name: "Porter Robinson",
	}
	bArtistLP = &Artist{
		ID:        1,
		Name:      "Linkin Park",
		SpotifyId: "spotify:linkin_park",
	}
	bArtistPR = &Artist{
		ID:        2,
		Name:      "Porter Robinson",
		SpotifyId: "spotify:porter_robinson",
	}
)

func Test_upsertArtistsSqlx(t *testing.T) {
	t.Run("returns same ID, when called multiple times with same SpotifyId", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if wantN, nArtists := 0, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}

		want1 := []*Artist{bArtistLP, bArtistPR}
		got1, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("upsertArtistsSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got2); diff != "" {
			t.Errorf("upsertArtistsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 2, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if wantN, nArtists := 0, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}

		want1 := []*Artist{bArtistLP}
		got1, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{sArtistLP}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("upsertArtistsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 1, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}

		want2 := []*Artist{bArtistPR}
		got2, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{sArtistPR}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want2, got2); diff != "" {
			t.Errorf("upsertArtistsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 2, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}
	})
}

func Test_upsertArtistsGorm(t *testing.T) {
	t.Run("returns same ID, when called multiple times with same SpotifyId", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if wantN, nArtists := 0, checkNArtistsGorm(t, brain.gormDb); nArtists != wantN {
			t.Fatalf("gorm have %d rows, but want %d rows", nArtists, wantN)
		}

		want1 := []*Artist{bArtistLP, bArtistPR}
		got1, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got2); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 2, checkNArtistsGorm(t, brain.gormDb); nArtists != wantN {
			t.Fatalf("gorm have %d rows, but want %d rows", nArtists, wantN)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if wantN, nArtists := 0, checkNArtistsGorm(t, brain.gormDb); nArtists != wantN {
			t.Fatalf("gorm have %d rows, but want %d rows", nArtists, wantN)
		}

		want1 := []*Artist{bArtistLP}
		got1, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 1, checkNArtistsGorm(t, brain.gormDb); nArtists != wantN {
			t.Fatalf("gorm have %d rows, but want %d rows", nArtists, wantN)
		}

		want2 := []*Artist{bArtistPR}
		got2, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistPR}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want2, got2); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 2, checkNArtistsGorm(t, brain.gormDb); nArtists != wantN {
			t.Fatalf("gorm have %d rows, but want %d rows", nArtists, wantN)
		}
	})
}
func diffArtists(want, got []*Artist) string {
	return cmp.Diff(want, got)
}

func checkNArtistsSqlx(tb testing.TB, db *sqlx.DB) int {
	sqlxArtists, err := getAllArtistsSqlx(db)
	tb.Logf("There are %d artists in sqlx db:\n", len(sqlxArtists))
	for idx, item := range sqlxArtists {
		tb.Logf("[%d/%d] artist: %+v", idx+1, len(sqlxArtists), item)
	}
	tb.Logf("---------")
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	return len(sqlxArtists)
}

func checkNArtistsGorm(tb testing.TB, db *gorm.DB) int {
	gormArtists, err := getAllArtistsGorm(db)
	tb.Logf("There are %d artists in gorm db:\n", len(gormArtists))
	for idx, item := range gormArtists {
		tb.Logf("[%d/%d] artist: %+v", idx+1, len(gormArtists), item)
	}
	tb.Logf("---------")
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	return len(gormArtists)
}
