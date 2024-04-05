package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/jmoiron/sqlx"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
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
		ID:             1,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: "linkin park - hybrid theory",
	}
	bMetaAlbumN = &MetaAlbum{
		ID:             2,
		Artists:        []*Artist{bArtistPR},
		SimplifiedName: "porter robinson - nurture",
	}
)

func Test_SaveAlbums(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsGorm(t, brain.gormDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bMetaAlbumN}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT, sFullAlbumHT20, sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("SaveAlbums() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT, sFullAlbumHT20, sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Errorf("SaveAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsGorm(t, brain.gormDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("SaveAlbums() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT20})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want2, got2); diff != "" {
			t.Errorf("SaveAlbums() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bMetaAlbumN}
		got3, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want3, got3); diff != "" {
			t.Errorf("SaveAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_upsertMetaAlbumsGorm(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsGorm(t, brain.gormDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bMetaAlbumN}
		got1, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsGorm(t, brain.gormDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT20}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want2, got2); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bMetaAlbumN}
		got3, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want3, got3); diff != "" {
			t.Errorf("upsertArtistsGorm() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_upsertMetaAlbumsSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsSqlx(t, brain.sqlxDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bMetaAlbumN}
		got1, err := upsertMetaAlbumsSqlx(brain.sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("upsertMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertMetaAlbumsSqlx(brain.sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Errorf("upsertMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsSqlx(t, brain.sqlxDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		// Setup Artists
		bi := newBrainIndex()
		if _, err := upsertArtistsSqlx(brain.sqlxDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := upsertMetaAlbumsSqlx(brain.sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("checkNMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := upsertMetaAlbumsSqlx(brain.sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumHT20}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want2, got2); diff != "" {
			t.Errorf("checkNMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bMetaAlbumN}
		got3, err := upsertMetaAlbumsSqlx(brain.sqlxDb, []spotify.SimpleAlbum{sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want3, got3); diff != "" {
			t.Errorf("checkNMetaAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_META_ALBUM_FIELDS = cmpopts.IgnoreFields(MetaAlbum{}, "AnyName")

func diffMetaAlbums(want, got []*MetaAlbum) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS)
}

func checkNMetaAlbumsGorm(tb testing.TB, db *gorm.DB) int {
	gormMetaAlbums, err := getAllMetaAlbumsGorm(db)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	// tb.Logf("There are %d artists in gorm db:\n", len(gormMetaAlbums))
	// for idx, item := range gormMetaAlbums {
	// 	tb.Logf("[%d/%d] artist: %+v", idx+1, len(gormMetaAlbums), item)
	// }
	// tb.Logf("---------")
	return len(gormMetaAlbums)
}

func checkNMetaAlbumsSqlx(tb testing.TB, db *sqlx.DB) int {
	sqlxMetaAlbums, err := getAllMetaAlbumsSqlx(db)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	// tb.Logf("There are %d artists in sqlx db:\n", len(sqlxMetaAlbums))
	// for idx, item := range sqlxMetaAlbums {
	// 	tb.Logf("[%d/%d] artist: %+v", idx+1, len(sqlxMetaAlbums), item)
	// }
	// tb.Logf("---------")
	return len(sqlxMetaAlbums)
}
