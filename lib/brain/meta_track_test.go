package brain

import (
	"log"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
)

var (
	sTrackOS = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:          "spotify:one_step",
				Name:        "One Step Closer",
				TrackNumber: 2,
				Artists:     []spotify.SimpleArtist{sArtistLP},
			},
			Album: sAlbumHT.SimpleAlbum,
		},
	}
	sTrackITE = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:          "spotify:in_the_end",
				Name:        "In the end",
				TrackNumber: 8,
				Artists:     []spotify.SimpleArtist{sArtistLP},
			},
			Album: sAlbumHT.SimpleAlbum,
		},
	}
	sTrackSC = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: spotify.SimpleTrack{
				ID:          "spotify:something_comforting",
				Name:        "Something Comforting",
				TrackNumber: 11,
				Artists:     []spotify.SimpleArtist{sArtistPR},
			},
			Album: sAlbumN.SimpleAlbum,
		},
	}
	bTrackOS = &MetaTrack{
		ID:             1,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: bMetaAlbumHT.SimplifiedName + " - 02.  one step closer",
		MetaAlbumID:    bMetaAlbumHT.ID,
		MetaAlbum:      bMetaAlbumHT,
	}
	bTrackITE = &MetaTrack{
		ID:             2,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: bMetaAlbumHT.SimplifiedName + " - 08.  in the end",
		MetaAlbumID:    bMetaAlbumHT.ID,
		MetaAlbum:      bMetaAlbumHT,
	}
	bTrackSC = &MetaTrack{
		ID:             3,
		Artists:        []*Artist{bArtistPR},
		SimplifiedName: bMetaAlbumN.SimplifiedName + " - 11.  something comforting",
		MetaAlbumID:    bMetaAlbumN.ID,
		MetaAlbum:      bMetaAlbumN,
	}
)

func TestSaveTracks(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if nTracks := logAllMetaTracks(t, brain); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaTrack{bTrackOS, bTrackITE}
		got1, err := brain.SaveTracks([]spotify.SavedTrack{sTrackOS, sTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveTracks([]spotify.SavedTrack{sTrackOS, sTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got2); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if nTracks := logAllMetaTracks(t, brain); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaTrack{bTrackOS}
		got1, err := brain.SaveTracks([]spotify.SavedTrack{sTrackOS}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bTrackITE}
		got2, err := brain.SaveTracks([]spotify.SavedTrack{sTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaTrack{bTrackSC}
		got3, err := brain.SaveTracks([]spotify.SavedTrack{sTrackSC}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		log.Printf("got1: %v", got1)
		log.Printf("got2: %v", got2)
		log.Printf("got3: %v", got3)
		if diff := diffMetaTracks(want3, got3); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}
	})
}

var IGNORE_META_TRACK_FIELDS = cmpopts.IgnoreFields(MetaTrack{}, "MetaAlbum")

// func diffMetaTrack(want, got *MetaTrack) string {
// 	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS, IGNORE_META_TRACK_FIELDS)
// }

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
