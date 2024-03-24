package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/zmb3/spotify/v2"
)

var (
	sArtistLP = spotify.SimpleArtist{
		ID:   "spotify:linkin_park",
		Name: "Linkin Park",
	}
	sArtistPR = spotify.SimpleArtist{
		ID:   "spotify:porter_robinson",
		Name: "Porter Robinson",
	}
	bArtistLP = &Artist{
		ID:        1,
		Name:      "Linkin Park",
		SpotifyId: "spotify:linkin_park",
	}
	bArtistPR = &Artist{
		ID:        2,
		Name:      "Porter Robinson",
		SpotifyId: "spotify:porter_robinson",
	}
)

func Test_saveArtists(t *testing.T) {
	t.Run("returns same ID, when called multiple times with same SpotifyId", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nArtists := logAllArtists(t, brain); nArtists != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*Artist{bArtistLP, bArtistPR}
		got1, err := brain._saveArtists([]spotify.SimpleArtist{sArtistLP, sArtistPR})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("_saveArtistss() mismatch (-want +got):\n%s", diff)
		}
		logAllArtists(t, brain)

		got2, err := brain._saveArtists([]spotify.SimpleArtist{sArtistLP, sArtistPR})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got2); diff != "" {
			t.Errorf("_saveArtistss() mismatch (-want +got):\n%s", diff)
		}
		logAllArtists(t, brain)
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nArtists := logAllArtists(t, brain); nArtists != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*Artist{bArtistLP}
		got1, err := brain._saveArtists([]spotify.SimpleArtist{sArtistLP})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want1, got1); diff != "" {
			t.Errorf("_saveArtistss() mismatch (-want +got):\n%s", diff)
		}
		logAllArtists(t, brain)

		want2 := []*Artist{bArtistLP, bArtistPR}
		got2, err := brain._saveArtists([]spotify.SimpleArtist{sArtistLP, sArtistLP, sArtistPR, sArtistPR})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffArtists(want2, got2); diff != "" {
			t.Errorf("_saveArtistss() mismatch (-want +got):\n%s", diff)
		}
		logAllArtists(t, brain)
	})
}

func diffArtist(want, got *Artist) string {
	return cmp.Diff(want, got)
}

func diffArtists(want, got []*Artist) string {
	return cmp.Diff(want, got)
}

func logAllArtists(tb testing.TB, brain *Brain) int {
	var allArtists []Artist
	if err := brain.gormDb.Find(&allArtists).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	tb.Logf("There are %d artists in db:\n", len(allArtists))
	for idx, item := range allArtists {
		tb.Logf("[%d/%d] artist: %+v", idx+1, len(allArtists), item)
	}
	return len(allArtists)
}
