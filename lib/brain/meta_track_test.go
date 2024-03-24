package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
)

var (
	sTrack1 = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:          ("spotify:one_step"),
				Name:        "One Step Closer",
				TrackNumber: 2,
				Artists:     []spotify.SimpleArtist{sArtistLP},
				Album:       sAlbumHT.SimpleAlbum,
			},
		},
	}
	sTrack2 = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:          ("spotify:in_the_end"),
				Name:        "In the end",
				TrackNumber: 8,
				Artists:     []spotify.SimpleArtist{sArtistLP},
				Album:       sAlbumHT.SimpleAlbum,
			},
		},
	}
	sTrack3 = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:          ("spotify:something_comforting"),
				Name:        "Something Comforting",
				TrackNumber: 11,
				Artists:     []spotify.SimpleArtist{sArtistPR},
				Album:       sAlbumN.SimpleAlbum,
			},
		},
	}
	bTrack1 = &MetaTrack{
		ID: 1,
		// Name:           "One Step Closer",
		// SpotifyId:      "spotify:one_step",
		Artists: []*Artist{bArtistLP},
		// SpotifyAlbum:   bAlbum1,
		// SpotifyAlbumId: bAlbum1.ID,
		// MetaTrackId:    1,
	}
	bTrack2 = &MetaTrack{
		ID: 2,
		// Name:           "In the end",
		// SpotifyId:      "spotify:in_the_end",
		Artists: []*Artist{bArtistLP},
		// SpotifyAlbum:   bAlbum1,
		// SpotifyAlbumId: bAlbum1.ID,
		// MetaTrackId:    2,
	}
	bTrack3 = &MetaTrack{
		ID: 3,
		// Name:           "Something Comforting",
		// SpotifyId:      "spotify:something_comforting",
		Artists: []*Artist{bArtistPR},
		// SpotifyAlbum:   bAlbum2,
		// SpotifyAlbumId: bAlbum2.ID,
		// MetaTrackId:    3,
	}
)

func TestSaveTracks(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nTracks := logAllMetaTracks(t, brain); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaTrack{bTrack1, bTrack2}
		got1, err := brain.SaveTracks([]spotify.SavedTrack{sTrack1, sTrack2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveTracks([]spotify.SavedTrack{sTrack1, sTrack2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got2); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if nTracks := logAllMetaTracks(t, brain); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaTrack{bTrack1}
		got1, err := brain.SaveTracks([]spotify.SavedTrack{sTrack1})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bTrack2}
		got2, err := brain.SaveTracks([]spotify.SavedTrack{sTrack2})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_META_TRACK_FIELDS = cmpopts.IgnoreFields(MetaTrack{}, "MetaAlbum")

func diffMetaTrack(want, got *MetaTrack) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS, IGNORE_META_TRACK_FIELDS)
}

func diffMetaTracks(want, got []*MetaTrack) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS, IGNORE_META_TRACK_FIELDS)
}

func logAllMetaTracks(tb testing.TB, brain *Brain) int {
	var allTracks []MetaTrack
	if err := brain.gormDb.Find(&allTracks).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	tb.Logf("There are %d tracks in db:\n", len(allTracks))
	for idx, item := range allTracks {
		tb.Logf("[%d/%d] track: %+v", idx+1, len(allTracks), item)
	}
	return len(allTracks)
}
