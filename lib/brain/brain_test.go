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
	return brain
}
