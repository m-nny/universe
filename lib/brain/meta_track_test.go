package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

var (
	sSimpleTrackOS = spotify.SimpleTrack{
		ID:          "spotify:one_step",
		Name:        "One Step Closer",
		TrackNumber: 2,
		Artists:     []spotify.SimpleArtist{sArtistLP},
		Album:       sSimpleAlbumHT,
	}
	sSimpleTrackITE = spotify.SimpleTrack{
		ID:          "spotify:in_the_end",
		Name:        "In the end",
		TrackNumber: 8,
		Artists:     []spotify.SimpleArtist{sArtistLP},
		Album:       sSimpleAlbumHT,
	}
	sSimpleTrackSC = spotify.SimpleTrack{
		ID:          "spotify:something_comforting",
		Name:        "Something Comforting",
		TrackNumber: 11,
		Artists:     []spotify.SimpleArtist{sArtistPR},
		Album:       sSimpleAlbumN,
	}
	sSavedTrackOS = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: sSimpleTrackOS,
			Album:       sSimpleAlbumHT,
		},
	}
	sSavedTrackITE = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: sSimpleTrackITE,
			Album:       sSimpleAlbumHT,
		},
	}
	sSavedTrackSC = spotify.SavedTrack{
		FullTrack: spotify.FullTrack{
			SimpleTrack: sSimpleTrackSC,
			Album:       sSimpleAlbumN,
		},
	}
	bMetaTrackOS = &MetaTrack{
		ID:             1,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: bMetaAlbumHT.SimplifiedName + " - 02.  one step closer",
		MetaAlbumID:    bMetaAlbumHT.ID,
		MetaAlbum:      bMetaAlbumHT,
	}
	bMetaTrackITE = &MetaTrack{
		ID:             2,
		Artists:        []*Artist{bArtistLP},
		SimplifiedName: bMetaAlbumHT.SimplifiedName + " - 08.  in the end",
		MetaAlbumID:    bMetaAlbumHT.ID,
		MetaAlbum:      bMetaAlbumHT,
	}
	bMetaTrackSC = &MetaTrack{
		ID:             3,
		Artists:        []*Artist{bArtistPR},
		SimplifiedName: bMetaAlbumN.SimplifiedName + " - 11.  something comforting",
		MetaAlbumID:    bMetaAlbumN.ID,
		MetaAlbum:      bMetaAlbumN,
	}
)

func Test_SaveTracks(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if nTracks := checkNMetaTracksGorm(t, brain.gormDb); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaTrack{bMetaTrackOS, bMetaTrackITE}
		got1, err := brain.SaveTracks([]spotify.SavedTrack{sSavedTrackOS, sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveTracks([]spotify.SavedTrack{sSavedTrackOS, sSavedTrackITE}, username)
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
		if nTracks := checkNMetaTracksGorm(t, brain.gormDb); nTracks != 0 {
			t.Fatalf("sqlite db is not clean")
		}

		want1 := []*MetaTrack{bMetaTrackOS}
		got1, err := brain.SaveTracks([]spotify.SavedTrack{sSavedTrackOS}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bMetaTrackITE}
		got2, err := brain.SaveTracks([]spotify.SavedTrack{sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaTrack{bMetaTrackSC}
		got3, err := brain.SaveTracks([]spotify.SavedTrack{sSavedTrackSC}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want3, got3); diff != "" {
			t.Errorf("SaveTracks() mismatch (-want +got):\n%s", diff)
		}
	})
}

func Test_upsertMetaTracksGorm(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		gormDb := getInmemoryBrain(t).gormDb
		// username := "test_username"
		if want, got := 0, checkNMetaTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsGorm(gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaTrack{bMetaTrackOS, bMetaTrackITE}
		got1, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("upsertMetaTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got2); diff != "" {
			t.Errorf("upsertMetaTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNMetaTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		gormDb := getInmemoryBrain(t).gormDb
		// username := "test_username"
		if want, got := 0, checkNMetaTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsGorm(gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*MetaTrack{bMetaTrackOS}
		got1, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("upsertMetaTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bMetaTrackITE}
		got2, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("upsertMetaTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaTrack{bMetaTrackSC}
		got3, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackSC}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want3, got3); diff != "" {
			t.Errorf("upsertMetaTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 3, checkNMetaTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}
	})
}

var IGNORE_META_TRACK_FIELDS = cmpopts.IgnoreFields(MetaTrack{}, "MetaAlbum")

func diffMetaTracks(want, got []*MetaTrack) string {
	return cmp.Diff(want, got, IGNORE_META_ALBUM_FIELDS, IGNORE_META_TRACK_FIELDS)
}

func checkNMetaTracksGorm(tb testing.TB, db *gorm.DB) int {
	var allTracks []MetaTrack
	if err := db.Find(&allTracks).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	return len(allTracks)
}
