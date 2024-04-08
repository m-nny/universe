package xiter

import (
	"iter"
	"strings"
)

func Map[T, U any](seq iter.Seq[T], fn func(item T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for item := range seq {
			if !yield(fn(item)) {
				return
			}
		}
	}
}

func Iter[T any](seq []T) iter.Seq[T] {
	return func(yield func(T) bool) {
		for _, item := range seq {
			if !yield(item) {
				return
			}
		}
	}
}

func Join(seq iter.Seq[string], sep string) string {
	var builder strings.Builder
	for item := range seq {
		if builder.Len() > 0 {
			builder.WriteString(sep)
		}
		builder.WriteString(item)
	}
	return builder.String()
}

func JoinFn[T any](seq iter.Seq[T], sep string, mapper func(item T) string) string {
	return Join(Map(seq, mapper), sep)
}

func SliceJoinFn[T any](seq []T, sep string, mapper func(item T) string) string {
	return Join(Map(Iter(seq), mapper), sep)
}

func Unique[T any, K comparable](arr iter.Seq[T], fn func(item T) K) iter.Seq[T] {
	return func(yield func(T) bool) {
		has := make(map[K]bool)
		for item := range arr {
			id := fn(item)
			if has[id] {
				continue
			}
			has[id] = true
			if !yield(item) {
				return
			}
		}
	}
}
