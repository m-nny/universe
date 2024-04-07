package xiter

import "iter"

func Map[T, U any](fn func(item T) U, seq iter.Seq[T]) iter.Seq[U] {
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
