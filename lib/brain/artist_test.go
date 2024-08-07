package brain

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
)

var (
	sArtistLP = spotify.SimpleArtist{
		ID:   "spotify:linkin_park",
		Name: "Linkin Park",
	}
	sArtistJZ = spotify.SimpleArtist{
		ID:   "spotify:jay_z",
		Name: "Jay-Z",
	}
	sArtistPR = spotify.SimpleArtist{
		ID:   "spotify:porter_robinson",
		Name: "Porter Robinson",
	}
	sArtistDT = spotify.SimpleArtist{
		ID:   "spotify:dream_theater",
		Name: "Dream Theater",
	}

	bArtistLP = &Artist{
		Name:      "Linkin Park",
		SpotifyId: "spotify:linkin_park",
	}
	bArtistJZ = &Artist{
		Name:      "Jay-Z",
		SpotifyId: "spotify:jay_z",
	}
	bArtistPR = &Artist{
		Name:      "Porter Robinson",
		SpotifyId: "spotify:porter_robinson",
	}
	bArtistDT = &Artist{
		Name:      "Dream Theater",
		SpotifyId: "spotify:dream_theater",
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
	t.Run("handles 0 args", func(t *testing.T) {
		brain := getInmemoryBrain(t)

		want1 := []*Artist{}
		got1, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{}, newBrainIndex())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("upsertArtistsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if wantN, nArtists := 0, checkNArtistsSqlx(t, brain.sqlxDb); nArtists != wantN {
			t.Fatalf("sqlx have %d rows, but want %d rows", nArtists, wantN)
		}
	})
}

var SORT_ARTISTS = cmpopts.SortSlices(func(a, b *Artist) bool {
	return strings.Compare(string(a.SpotifyId), string(b.SpotifyId)) < 0
})

func diffArtists(want, got []*Artist) string {
	return cmp.Diff(want, got)
}

func checkNArtistsSqlx(tb testing.TB, db *sqlx.DB) int {
	sqlxArtists, err := getAllArtistsSqlx(db)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	// tb.Logf("There are %d artists in sqlx db:\n", len(sqlxArtists))
	// for idx, item := range sqlxArtists {
	// 	tb.Logf("[%d/%d] artist: %+v", idx+1, len(sqlxArtists), item)
	// }
	// tb.Logf("---------")
	return len(sqlxArtists)
}
