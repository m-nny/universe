package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
)

var (
	sAlbum1 = &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:      ("spotify:hybrid_theory"),
			Name:    "Hybrid Theory",
			Artists: []spotify.SimpleArtist{sArtist1},
		},
	}
	sAlbum2 = &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:      ("spotify:hybryd_theory_20"),
			Name:    "Hybrid Theory (20th Anniversary Edition)",
			Artists: []spotify.SimpleArtist{sArtist1},
		},
	}
	sAlbum3 = &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:      ("spotify:nurture"),
			Name:    "Nurture",
			Artists: []spotify.SimpleArtist{sArtist2},
		},
	}
	bAlbum1 = &SpotifyAlbum{
		ID:          1,
		Name:        "Hybrid Theory",
		SpotifyId:   "spotify:hybrid_theory",
		Artists:     []*Artist{bArtist1},
		MetaAlbumId: 1,
	}
	bAlbum2 = &SpotifyAlbum{
		ID:          2,
		Name:        "Hybrid Theory (20th Anniversary Edition)",
		SpotifyId:   "spotify:hybryd_theory_20",
		Artists:     []*Artist{bArtist1},
		MetaAlbumId: 1,
	}
	bAlbum3 = &SpotifyAlbum{
		ID:          3,
		Name:        "Nurture",
		SpotifyId:   "spotify:nurture",
		Artists:     []*Artist{bArtist2},
		MetaAlbumId: 2,
	}
)

func TestToAlbums(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := logAllAlbums(t, brain); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*SpotifyAlbum{bAlbum1, bAlbum2, bAlbum3}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbum1, sAlbum2, sAlbum3})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffAlbums(want1, got1); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbum1, sAlbum2, sAlbum3})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffAlbums(want1, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := logAllAlbums(t, brain); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*SpotifyAlbum{bAlbum1}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbum1})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffAlbums(want1, got1); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*SpotifyAlbum{bAlbum2}
		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbum2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffAlbums(want2, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*SpotifyAlbum{bAlbum3}
		got3, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbum3})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffAlbums(want3, got3); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_ALBUM_FIELDS = cmpopts.IgnoreFields(SpotifyAlbum{}, "MetaAlbum")

func diffAlbum(want, got *SpotifyAlbum) string {
	return cmp.Diff(want, got, IGNORE_ALBUM_FIELDS)
}

func diffAlbums(want, got []*SpotifyAlbum) string {
	return cmp.Diff(want, got, IGNORE_ALBUM_FIELDS)
}

func logAllAlbums(tb testing.TB, brain *Brain) int {
	var allAlbums []SpotifyAlbum
	if err := brain.gormDb.Find(&allAlbums).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	tb.Logf("There are %d albums in db:\n", len(allAlbums))
	for idx, item := range allAlbums {
		tb.Logf("[%d/%d] album: %+v", idx+1, len(allAlbums), item)
	}
	return len(allAlbums)
}
