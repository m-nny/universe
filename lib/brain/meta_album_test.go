package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
)

var (
	sAlbumHT = &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:      "spotify:hybrid_theory",
			Name:    "Hybrid Theory",
			Artists: []spotify.SimpleArtist{sArtistLP},
		},
	}
	sAlbumHT20 = &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:      "spotify:hybryd_theory_20",
			Name:    "Hybrid Theory (20th Anniversary Edition)",
			Artists: []spotify.SimpleArtist{sArtistLP},
		},
	}
	sAlbumN = &spotify.FullAlbum{
		SimpleAlbum: spotify.SimpleAlbum{
			ID:      "spotify:nurture",
			Name:    "Nurture",
			Artists: []spotify.SimpleArtist{sArtistPR},
		},
	}
	bMetaAlbumHT = &MetaAlbum{
		ID:      1,
		Artists: []*Artist{bArtistLP},

		SimplifiedName: "linkin park - hybrid theory",
	}
	bmetaAlbumN = &MetaAlbum{
		ID:             2,
		Artists:        []*Artist{bArtistPR},
		SimplifiedName: "porter robinson - nurture",
	}
)

func TestToAlbums(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := logAllAlbums(t, brain); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bmetaAlbumN}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbumHT, sAlbumHT20, sAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMeteaAlbums(want1, got1); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbumHT, sAlbumHT20, sAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMeteaAlbums(want1, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nAlbums := logAllAlbums(t, brain); nAlbums != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbumHT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMeteaAlbums(want1, got1); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbumHT20})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMeteaAlbums(want2, got2); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bmetaAlbumN}
		got3, err := brain.SaveAlbums([]*spotify.FullAlbum{sAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMeteaAlbums(want3, got3); diff != "" {
			t.Errorf("ToAlbums() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_META_ALBUM_FIELDS = cmpopts.IgnoreFields(MetaAlbum{}, "AnyName")

func diffMetaAlbum(want, got *MetaAlbum) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS)
}

func diffMeteaAlbums(want, got []*MetaAlbum) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS)
}

func logAllAlbums(tb testing.TB, brain *Brain) int {
	var allAlbums []MetaAlbum
	if err := brain.gormDb.Find(&allAlbums).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	tb.Logf("There are %d albums in db:\n", len(allAlbums))
	for idx, item := range allAlbums {
		tb.Logf("[%d/%d] album: %+v", idx+1, len(allAlbums), item)
	}
	return len(allAlbums)
}
