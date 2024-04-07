package brain

import (
	"testing"

	"github.com/zmb3/spotify/v2"
)

func Test_SaveTracksGorm(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if want, got := 0, checkNMetaTracksGorm(t, brain.gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}

		want1 := []*MetaTrack{bMetaTrackOS, bMetaTrackITE}
		got1, err := brain.SaveTracksGorm([]spotify.SavedTrack{sSavedTrackOS, sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveTracksGorm([]spotify.SavedTrack{sSavedTrackOS, sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got2); diff != "" {
			t.Errorf("SaveTracksGorm() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 2, checkNMetaTracksGorm(t, brain.gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if want, got := 0, checkNMetaTracksGorm(t, brain.gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}

		want1 := []*MetaTrack{bMetaTrackOS}
		got1, err := brain.SaveTracksGorm([]spotify.SavedTrack{sSavedTrackOS}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bMetaTrackITE}
		got2, err := brain.SaveTracksGorm([]spotify.SavedTrack{sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("SaveTracksGorm() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaTrack{bMetaTrackSC}
		got3, err := brain.SaveTracksGorm([]spotify.SavedTrack{sSavedTrackSC}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want3, got3); diff != "" {
			t.Errorf("SaveTracksGorm() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 3, checkNMetaTracksGorm(t, brain.gormDb); got != want {
			t.Fatalf("gorm has %d meta tracks, but want %d rows", got, want)
		}
	})
}

func Test_SaveTracksSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if want, got := 0, checkNMetaTracksSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}

		want1 := []*MetaTrack{bMetaTrackOS, bMetaTrackITE}
		got1, err := brain.SaveTracksSqlx([]spotify.SavedTrack{sSavedTrackOS, sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveTracksSqlx([]spotify.SavedTrack{sSavedTrackOS, sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got2); diff != "" {
			t.Errorf("SaveTracksSqlx() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 2, checkNMetaTracksSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		username := "test_username"
		if want, got := 0, checkNMetaTracksSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}

		want1 := []*MetaTrack{bMetaTrackOS}
		got1, err := brain.SaveTracksSqlx([]spotify.SavedTrack{sSavedTrackOS}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want1, got1); diff != "" {
			t.Errorf("SaveTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaTrack{bMetaTrackITE}
		got2, err := brain.SaveTracksSqlx([]spotify.SavedTrack{sSavedTrackITE}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want2, got2); diff != "" {
			t.Errorf("SaveTracksSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaTrack{bMetaTrackSC}
		got3, err := brain.SaveTracksSqlx([]spotify.SavedTrack{sSavedTrackSC}, username)
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaTracks(want3, got3); diff != "" {
			t.Errorf("SaveTracksSqlx() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 3, checkNMetaTracksSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta tracks, but want %d rows", got, want)
		}
	})
}
