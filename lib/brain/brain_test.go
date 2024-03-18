package brain

import (
	"testing"
)

func getInmemoryBrain(tb testing.TB) *Brain {
	dbName := "file::memory:"
	brain, err := New(dbName /*enableLogging=*/, false)
	if err != nil {
		tb.Fatalf("err: %v", err)
	}
	return brain
}
