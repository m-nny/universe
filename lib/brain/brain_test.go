package brain

import (
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func getInmemoryBrain(tb testing.TB) *Brain {
	brain, err := New("file::memory:" /*enableLogging=*/, false)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	return brain
}
