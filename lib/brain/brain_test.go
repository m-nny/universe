package brain

import (
	"testing"
)

func getInmemoryBrain(tb testing.TB) *Brain {
	dbName := "file::memory:"
	brain, err := New(dbName)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	if nArtists := logAllArists(tb, brain); nArtists != 0 {
		tb.Fatalf("sqlite db is not clean")
	}
	return brain
}
