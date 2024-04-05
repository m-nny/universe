package brain

import (
	"testing"
)

func getInmemoryBrain(tb testing.TB) *Brain {
	brain, err := New("file::memory:", "file::memory:" /*enableLogging=*/, false)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	return brain
}
