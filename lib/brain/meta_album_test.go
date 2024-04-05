package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
		if nAlbums := checkNMetaAlbumsGorm(t, brain); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bMetaAlbumN}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT, sFullAlbumHT20, sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT, sFullAlbumHT20, sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsGorm(t, brain); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumHT20})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want2, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bMetaAlbumN}
		got3, err := brain.SaveAlbums([]*spotify.FullAlbum{sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want3, got3); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_upsertMetaAlbumsGorm(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := checkNMetaAlbumsGorm(t, brain); nAlbums != 0 {
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
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertMetaAlbumsGorm(brain.gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumHT20, sSimpleAlbumN}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_META_ALBUM_FIELDS = cmpopts.IgnoreFields(MetaAlbum{}, "AnyName")

func diffMetaAlbums(want, got []*MetaAlbum) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS)
}

func checkNMetaAlbumsGorm(tb testing.TB, brain *Brain) int {
	var allAlbums []MetaAlbum
	if err := brain.gormDb.Find(&allAlbums).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	// tb.Logf("There are %d albums in db:\n", len(allAlbums))
	// for idx, item := range allAlbums {
	// 	tb.Logf("[%d/%d] album: %+v", idx+1, len(allAlbums), item)
	// }
	return len(allAlbums)
}
