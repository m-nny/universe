package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

var (
	sTrack1 = &spotify.SimpleTrack{
		ID:      ("spotify:one_step"),
		Name:    "One Step Closer",
		Artists: []spotify.SimpleArtist{*sArtist1},
		Album:   *sAlbum1,
	}
	sTrack2 = &spotify.SimpleTrack{
		ID:      ("spotify:in_the_end"),
		Name:    "In the end",
		Artists: []spotify.SimpleArtist{*sArtist1},
		Album:   *sAlbum1,
	}
	sTrack3 = &spotify.SimpleTrack{
		ID:      ("spotify:something_comforting"),
		Name:    "Something Comforting",
		Artists: []spotify.SimpleArtist{*sArtist2},
		Album:   *sAlbum2,
	}
	bTrack1 = &Track{
		Model:     gorm.Model{ID: 1},
		Name:      "One Step Closer",
		SpotifyId: "spotify:one_step",
		Artists:   []*Artist{bArtist1},
		Album:     bAlbum1,
		AlbumId:   bAlbum1.ID,
	}
	bTrack2 = &Track{
		Model:     gorm.Model{ID: 2},
		Name:      "In the end",
		SpotifyId: "spotify:in_the_end",
		Artists:   []*Artist{bArtist1},
		Album:     bAlbum1,
		AlbumId:   bAlbum1.ID,
	}
	bTrack3 = &Track{
		Model:     gorm.Model{ID: 3},
		Name:      "Something Comforting",
		SpotifyId: "spotify:something_comforting",
		Artists:   []*Artist{bArtist2},
		Album:     bAlbum2,
		AlbumId:   bAlbum2.ID,
	}
)

func TestToTracks(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nTracks := logAllTracks(t, brain); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*Track{bTrack1, bTrack2}
		got1, err := brain.ToTracks([]*spotify.SimpleTrack{sTrack1, sTrack2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffTracks(want1, got1); diff != "" {
			t.Errorf("ToTracks() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.ToTracks([]*spotify.SimpleTrack{sTrack1, sTrack2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffTracks(want1, got2); diff != "" {
			t.Errorf("ToTracks() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nTracks := logAllTracks(t, brain); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*Track{bTrack1}
		got1, err := brain.ToTracks([]*spotify.SimpleTrack{sTrack1})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffTracks(want1, got1); diff != "" {
			t.Errorf("ToTracks() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*Track{bTrack2}
		got2, err := brain.ToTracks([]*spotify.SimpleTrack{sTrack2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffTracks(want2, got2); diff != "" {
			t.Errorf("ToTracks() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_TRACK_FIELDS = cmpopts.IgnoreFields(Track{}, "Model.CreatedAt", "Model.UpdatedAt", "Model.DeletedAt", "Album")

func diffTrack(want, got *Track) string {
	return cmp.Diff(want, got, IGNORE_ALBUM_FIELDS, IGNORE_ARTIST_FIELDS, IGNORE_TRACK_FIELDS)
}

func diffTracks(want, got []*Track) string {
	return cmp.Diff(want, got, IGNORE_ALBUM_FIELDS, IGNORE_ARTIST_FIELDS, IGNORE_TRACK_FIELDS)
}

func logAllTracks(tb testing.TB, brain *Brain) int {
	var allTracks []Track
	if err := brain.gormDb.Find(&allTracks).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	tb.Logf("There are %d tracks in db:\n", len(allTracks))
	for idx, item := range allTracks {
		tb.Logf("[%d/%d] track: %+v", idx+1, len(allTracks), item)
	}
	return len(allTracks)
}
