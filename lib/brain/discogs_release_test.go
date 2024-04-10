package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/m-nny/universe/lib/discogs"
	"github.com/zmb3/spotify/v2"
)

// inputs
var (
	dListingReleaseDT = discogs.ListingRelease{
		Artist:        "Dream Theater",
		CatalogNumber: "7559-62742-2",
		Description:   "Dream Theater - Six Degrees Of Inner Turbulence (2xCD, Album)",
		Format:        "2xCD, Album",
		ID:            2519936,
		ResourceURL:   "https://api.discogs.com/releases/2519936",
		Thumbnail:     "https://i.discogs.com/_G1SttzK5tjvItnMongUP2WBfFEe11JHvBY1yPxxP88/rs:fit/g:sm/q:40/h:150/w:150/czM6Ly9kaXNjb2dz/LWRhdGFiYXNlLWlt/YWdlcy9SLTI1MTk5/MzYtMTI4ODUwOTU5/NC5qcGVn.jpeg",
		Title:         "Six Degrees Of Inner Turbulence",
		Year:          2002,
	}
	dListingReleaseIR = discogs.ListingRelease{
		Artist:        "I Romans",
		CatalogNumber: "YEP 00680",
		Description:   "I Romans - Coniglietto (7\")",
		Format:        "7\"",
		ID:            1485800,
		ResourceURL:   "https://api.discogs.com/releases/1485800",
		Thumbnail:     "https://i.discogs.com/CNJkbQgvZD_Pw4aUNxJeMgNzf7qsTgKoJ5SCEb2ikok/rs:fit/g:sm/q:40/h:150/w:150/czM6Ly9kaXNjb2dz/LWRhdGFiYXNlLWlt/YWdlcy9SLTE0ODU4/MDAtMTUwODg3ODYy/MS0xNTg2LmpwZWc.jpeg",
		Title:         "Coniglietto",
		Year:          1976,
	}
)

// outputs
var (
	bDiscogsReleaseDT = &DiscogsRelease{
		ArtistName: "Dream Theater",
		DiscogsID:  2519936,
		Format:     "2xCD, Album",
		Name:       "Six Degrees Of Inner Turbulence",
	}
	bDiscogsReleaseIR = &DiscogsRelease{
		ArtistName: "I Romans",
		DiscogsID:  1485800,
		Format:     `7"`,
		Name:       "Coniglietto",
	}
)

func Test_upsertDiscogsRelease(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}

		want1 := []*DiscogsRelease{bDiscogsReleaseDT, bDiscogsReleaseIR}
		got1, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{dListingReleaseDT, dListingReleaseIR})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffDiscogsReleases(want1, got1); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{dListingReleaseDT, dListingReleaseIR})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffDiscogsReleases(want1, got2); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}

		want1 := []*DiscogsRelease{bDiscogsReleaseDT}
		got1, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{dListingReleaseDT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffDiscogsReleases(want1, got1); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*DiscogsRelease{bDiscogsReleaseIR}
		got2, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{dListingReleaseIR})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffDiscogsReleases(want2, got2); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}

	})
	t.Run("handles 0 args", func(t *testing.T) {
		sqlxDb := getInmemoryBrain(t).sqlxDb
		// username := "test_username"
		if want, got := 0, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}

		want := []*DiscogsRelease{}
		got, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffDiscogsReleases(want, got); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 0, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}
	})
}

func Test_AssociateDiscogsRelease(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		sqlxDb := brain.sqlxDb
		// username := "test_username"
		if want, got := 0, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}

		want1 := []*DiscogsRelease{bDiscogsReleaseDT}
		got1, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{dListingReleaseDT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffDiscogsReleases(want1, got1); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(sqlxDb, []spotify.SimpleArtist{sArtistDT}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsSqlx(sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumDT}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		if err := brain.AssociateDiscogsRelease(got1[0], bMetaAlbumDT, 100); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		got2, err := upsertDiscogsReleases(sqlxDb, []discogs.ListingRelease{dListingReleaseDT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		wantDT := *bDiscogsReleaseDT
		wantDT.MetaAlbumId = &bMetaAlbumDT.SimplifiedName
		wantDT.MetaAlbumSimilariryScore = 100
		wantDT.SearchedMetaAlbum = true
		want2 := []*DiscogsRelease{&wantDT}
		if diff := diffDiscogsReleases(want2, got2); diff != "" {
			t.Errorf("upsertDiscogsReleases() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 1, checkNDiscogsReleasesSqlx(t, sqlxDb); got != want {
			t.Fatalf("sqlx has %d discogs releases, but want %d rows", got, want)
		}
	})
}

func diffDiscogsReleases(want, got []*DiscogsRelease) string {
	return cmp.Diff(want, got, cmpopts.SortSlices(func(a, b *DiscogsRelease) bool { return a.DiscogsID < b.DiscogsID }))
}

func checkNDiscogsReleasesSqlx(tb testing.TB, db *sqlx.DB) int {
	var cnt int
	if err := db.Get(&cnt, `SELECT COUNT(*) FROM discogs_releases`); err != nil {
		tb.Fatalf("err: %v", err)
	}
	return cnt
}
