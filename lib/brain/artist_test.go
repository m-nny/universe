package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

var (
	sArtist1 = &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{
			ID:   spotify.ID("spotify:linkin_park"),
			Name: "Linkin Park",
		},
	}
	sArtist2 = &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{
			ID:   spotify.ID("spotify:porter_robinson"),
			Name: "Porter Robinson",
		},
	}
	bArtist1 = &Artist{
		Model:     gorm.Model{ID: 1},
		Name:      "Linkin Park",
		SpotifyId: "spotify:linkin_park",
	}
	bArtist2 = &Artist{
		Model:     gorm.Model{ID: 2},
		Name:      "Porter Robinson",
		SpotifyId: "spotify:porter_robinson",
	}
)

func TestToArtist(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		got1, err := brain.ToArtist(sArtist1)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtist(bArtist1, got1); diff != "" {
			t.Errorf("ToArtist() mismatch (-want +got):\n%s", diff)
		}
		got2, err := brain.ToArtist(sArtist1)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtist(bArtist1, got2); diff != "" {
			t.Errorf("ToArtist() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		got1, err := brain.ToArtist(sArtist1)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtist(bArtist1, got1); diff != "" {
			t.Errorf("ToArtist() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.ToArtist(sArtist2)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtist(bArtist2, got2); diff != "" {
			t.Errorf("ToArtist() mismatch (-want +got):\n%s", diff)
		}
	})
}

func TestToArtists(t *testing.T) {
	t.Run("returns same ID, when called multiple times with same SpotifyId", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		logAllArists(t, brain)

		want1 := []*Artist{bArtist1, bArtist2}
		got1, err := brain.ToArtists([]*spotify.FullArtist{sArtist1, sArtist2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("ToArtists() mismatch (-want +got):\n%s", diff)
		}
		logAllArists(t, brain)

		got2, err := brain.ToArtists([]*spotify.FullArtist{sArtist1, sArtist2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got2); diff != "" {
			t.Errorf("ToArtists() mismatch (-want +got):\n%s", diff)
		}
		logAllArists(t, brain)
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		logAllArists(t, brain)

		want1 := []*Artist{bArtist1}
		got1, err := brain.ToArtists([]*spotify.FullArtist{sArtist1})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("ToArtists() mismatch (-want +got):\n%s", diff)
		}
		logAllArists(t, brain)

		want2 := []*Artist{bArtist1, bArtist2}
		got2, err := brain.ToArtists([]*spotify.FullArtist{sArtist1, sArtist1, sArtist2, sArtist2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want2, got2); diff != "" {
			t.Errorf("ToArtists() mismatch (-want +got):\n%s", diff)
		}
		logAllArists(t, brain)
	})
}

func diffArtist(want, got *Artist) string {
	return cmp.Diff(want, got, cmpopts.IgnoreFields(Artist{}, "Model.CreatedAt", "Model.UpdatedAt", "Model.DeletedAt"))
}

func diffArtists(want, got []*Artist) string {
	return cmp.Diff(want, got, cmpopts.IgnoreFields(Artist{}, "Model.CreatedAt", "Model.UpdatedAt", "Model.DeletedAt"))
}

func logAllArists(tb testing.TB, brain *Brain) int {
	var allArists []Artist
	if err := brain.gormDb.Find(&allArists).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	tb.Logf("There are %d artists in db:\n", len(allArists))
	for idx, item := range allArists {
		tb.Logf("[%d/%d] artist: %+v", idx+1, len(allArists), item)
	}
	return len(allArists)
}
