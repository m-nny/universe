package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
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

func Test_upsertSpotifyAlbumsGorm(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNSpotifyAlbumsGorm(t, brain.gormDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		// Setup Artists & MetaAlbums
		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyAlbum{bSpotifyAlbumHT, bSpotifyAlbumHT20, bSpotifyAlbumN}
		got1, err := upsertSpotifyAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want1, got1); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsGorm() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertSpotifyAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want1, got2); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsGorm() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNSpotifyAlbumsGorm(t, brain.gormDb); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		// Setup Artists & MetaAlbums
		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(brain.gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyAlbum{bSpotifyAlbumHT}
		got1, err := upsertSpotifyAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want1, got1); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsGorm() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*SpotifyAlbum{bSpotifyAlbumHT20}
		got2, err := upsertSpotifyAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT20}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want2, got2); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsGorm() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*SpotifyAlbum{bSpotifyAlbumN}
		got3, err := upsertSpotifyAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyAlbums(want3, got3); diff != "" {
			t.Fatalf("upsertSpotifyAlbumsGorm() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_SPOTIFY_ALBUM_FIELDS = cmpopts.IgnoreFields(SpotifyAlbum{}, "Artists", "MetaAlbum")

func diffSpotifyAlbums(want, got []*SpotifyAlbum) string {
	return cmp.Diff(want, got, IGNORE_SPOTIFY_ALBUM_FIELDS)
}

func checkNSpotifyAlbumsGorm(tb testing.TB, db *gorm.DB) int {
	var gormSpotifyAlbums []SpotifyAlbum
	if err := db.Find(&gormSpotifyAlbums).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	return len(gormSpotifyAlbums)
}
