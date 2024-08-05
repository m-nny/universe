package logutils

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"time"
)

// Infof is an example of a user-defined logging function that wraps slog.
// The log record contains the source position of the caller of Infof.
func Infof(format string, args ...any) {
	logger := slog.Default()
	if !logger.Enabled(context.Background(), slog.LevelInfo) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	r := slog.NewRecord(time.Now(), slog.LevelInfo, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

// Debugf is an example of a user-defined logging function that wraps slog.
// The log record contains the source position of the caller of Infof.
func Debugf(format string, args ...any) {
	logger := slog.Default()
	if !logger.Enabled(context.Background(), slog.LevelDebug) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Debugf]
	r := slog.NewRecord(time.Now(), slog.LevelDebug, fmt.Sprintf(format, args...), pcs[0])
	_ = logger.Handler().Handle(context.Background(), r)
}

// Fatalf is an example of a user-defined logging function that wraps slog.
// The log record contains the source position of the caller of Infof.
func Fatal(msg string, args ...any) {
	logger := slog.Default()
	if !logger.Enabled(context.Background(), slog.LevelError) {
		return
	}
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Fatalf]
	r := slog.NewRecord(time.Now(), slog.LevelError, msg, pcs[0])
	r.Add(args...)
	_ = logger.Handler().Handle(context.Background(), r)
	os.Exit(1)
}
