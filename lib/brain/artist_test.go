package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

func getInmemoryBrain(tb testing.TB) *Brain {
	dbName := "file::memory:?cache=shared"
	brain, err := New(dbName)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	return brain
}

func TestToArtist(t *testing.T) {
	brain := getInmemoryBrain(t)
	sArtist1 := &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{
			ID:   spotify.ID("spotify:linkin_park"),
			Name: "Linkin Park",
		},
	}
	sArtist2 := &spotify.FullArtist{
		SimpleArtist: spotify.SimpleArtist{
			ID:   spotify.ID("spotify:porter_robinson"),
			Name: "Porter Robinson",
		},
	}
	bArtist1 := &Artist{
		Model:     gorm.Model{ID: 1},
		Name:      "Linkin Park",
		SpotifyId: "spotify:linkin_park",
	}
	bArtist2 := &Artist{
		Model:     gorm.Model{ID: 2},
		Name:      "Porter Robinson",
		SpotifyId: "spotify:porter_robinson",
	}
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
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

func diffArtist(want, got *Artist) string {
	return cmp.Diff(want, got, cmpopts.IgnoreFields(Artist{}, "Model.CreatedAt", "Model.UpdatedAt", "Model.DeletedAt"))
}
