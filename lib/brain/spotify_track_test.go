package brain

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/zmb3/spotify/v2"
	"gorm.io/gorm"
)

// Expected outputs
var (
	bSpotifyTrackOS = &SpotifyTrack{
		Artists:        []*Artist{bArtistLP},
		ID:             1,
		MetaTrackId:    1,
		Name:           "One Step Closer",
		SpotifyAlbumId: 1,
		SpotifyId:      "spotify:one_step",
	}
	bSpotifyTrackITE = &SpotifyTrack{
		Artists:        []*Artist{bArtistLP},
		ID:             2,
		MetaTrackId:    2,
		Name:           "In the end",
		SpotifyAlbumId: 1,
		SpotifyId:      "spotify:in_the_end",
	}
	bSpotifyTrackSC = &SpotifyTrack{
		Artists:        []*Artist{bArtistPR},
		ID:             3,
		MetaTrackId:    3,
		Name:           "Something Comforting",
		SpotifyAlbumId: 2,
		SpotifyId:      "spotify:something_comforting",
	}
)

func Test_upsertSpotifyTracksGorm(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		gormDb := getInmemoryBrain(t).gormDb
		// username := "test_username"
		if want, got := 0, checkNSpotifyTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsGorm(gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertSpotifyAlbumsGorm(gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyTrack{bSpotifyTrackOS, bSpotifyTrackITE}
		got1, err := upsertSpotifyTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want1, got1); diff != "" {
			t.Errorf("upsertSpotifyTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		got2, err := upsertSpotifyTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want1, got2); diff != "" {
			t.Errorf("upsertSpotifyTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNSpotifyTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		gormDb := getInmemoryBrain(t).gormDb
		// username := "test_username"
		if want, got := 0, checkNSpotifyTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}

		bi := newBrainIndex()
		if _, err := upsertArtistsGorm(gormDb, []spotify.SimpleArtist{sArtistLP, sArtistPR}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaAlbumsGorm(gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertSpotifyAlbumsGorm(gormDb, []spotify.SimpleAlbum{sSimpleAlbumHT, sSimpleAlbumN}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if _, err := upsertMetaTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS, sSimpleTrackITE, sSimpleTrackSC}, bi); err != nil {
			t.Fatalf("got Error: %v", err)
		}

		want1 := []*SpotifyTrack{bSpotifyTrackOS}
		got1, err := upsertSpotifyTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackOS}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want1, got1); diff != "" {
			t.Errorf("upsertSpotifyTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*SpotifyTrack{bSpotifyTrackITE}
		got2, err := upsertSpotifyTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackITE}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want2, got2); diff != "" {
			t.Errorf("upsertSpotifyTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*SpotifyTrack{bSpotifyTrackSC}
		got3, err := upsertSpotifyTracksGorm(gormDb, []spotify.SimpleTrack{sSimpleTrackSC}, bi.Clone())
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffSpotifyTracks(want3, got3); diff != "" {
			t.Errorf("upsertSpotifyTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 3, checkNSpotifyTracksGorm(t, gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}
	})

}

var IGNORE_SPOTIFY_TRACK_FIELDS = cmpopts.IgnoreFields(SpotifyTrack{}, "Artists", "MetaTrack", "SpotifyAlbum")

func diffSpotifyTracks(want, got []*SpotifyTrack) string {
	return cmp.Diff(want, got, IGNORE_SPOTIFY_TRACK_FIELDS)
}

func checkNSpotifyTracksGorm(tb testing.TB, db *gorm.DB) int {
	var gormSpotifyTracks []SpotifyTrack
	if err := db.Find(&gormSpotifyTracks).Error; err != nil {
		tb.Fatalf("err: %v", err)
	}
	return len(gormSpotifyTracks)
}
