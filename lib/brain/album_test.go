package brain

import (
	"testing"

	"github.com/zmb3/spotify/v2"
)

func Test_SaveAlbumsSqlx(t *testing.T) {
	t.Run("returns same ID for same spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if want, got := 0, checkNMetaAlbumsSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT, bMetaAlbumN}
		got1, err := brain.SaveAlbumsSqlx([]*spotify.FullAlbum{sFullAlbumHT, sFullAlbumHT20, sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Fatalf("SaveAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		got2, err := brain.SaveAlbumsSqlx([]*spotify.FullAlbum{sFullAlbumHT, sFullAlbumHT20, sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got2); diff != "" {
			t.Fatalf("SaveAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 2, checkNMetaAlbumsSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}
	})
	t.Run("returns different ID for different spotify ID", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if want, got := 0, checkNMetaAlbumsSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}

		want1 := []*MetaAlbum{bMetaAlbumHT}
		got1, err := brain.SaveAlbumsSqlx([]*spotify.FullAlbum{sFullAlbumHT})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Fatalf("SaveAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want2 := []*MetaAlbum{bMetaAlbumHT}
		got2, err := brain.SaveAlbumsSqlx([]*spotify.FullAlbum{sFullAlbumHT20})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want2, got2); diff != "" {
			t.Fatalf("SaveAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		want3 := []*MetaAlbum{bMetaAlbumN}
		got3, err := brain.SaveAlbumsSqlx([]*spotify.FullAlbum{sFullAlbumN})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want3, got3); diff != "" {
			t.Fatalf("SaveAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}
		if want, got := 2, checkNMetaAlbumsSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}
	})
	t.Run("handles 0 args", func(t *testing.T) {
		brain := getInmemoryBrain(t)
		if want, got := 0, checkNMetaAlbumsSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}

		want1 := []*MetaAlbum{}
		got1, err := brain.SaveAlbumsSqlx([]*spotify.FullAlbum{})
		if err != nil {
			t.Fatalf("got Error: %v", err)
		}
		if diff := diffMetaAlbums(want1, got1); diff != "" {
			t.Fatalf("SaveAlbumsSqlx() mismatch (-want +got):\n%s", diff)
		}

		if want, got := 0, checkNMetaAlbumsSqlx(t, brain.sqlxDb); got != want {
			t.Fatalf("sqlx has %d meta albums, but want %d rows", got, want)
		}
	})
}
