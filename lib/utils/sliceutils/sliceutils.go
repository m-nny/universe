package sliceutils

import "context"

func Map[T, R any](arr []T, fn func(item T) R) []R {
	var res []R
	for _, item := range arr {
		res = append(res, fn(item))
	}
	return res
}

func MapP[T any](arr []T) []*T {
	return Map[T, *T](arr, func(item T) *T { return &item })
}

func MapErr[T, R any](arr []T, fn func(item T) (R, error)) ([]R, error) {
	var res []R
	for _, item := range arr {
		val, err := fn(item)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func MapCtxErr[T, R any](ctx context.Context, arr []T, fn func(ctx context.Context, item T) (R, error)) ([]R, error) {
	var res []R
	for _, item := range arr {
		val, err := fn(ctx, item)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func Uniqe[T any](arr []T, fn func(item T) string) []T {
	var result []T
	has := make(map[string]bool)
	for _, item := range arr {
		id := fn(item)
		if has[id] {
			continue
		}
		has[id] = true
		result = append(result, item)
	}
	return result
}

func Cnt[T any](arr []T, fn func(item T) int) int {
	val := 0
	for _, item := range arr {
		val += fn(item)
	}
	return val
}

func Identity[T any](item T) T {
	return item
}

func ToMap[K comparable, V any](arr []V, fn func(item V) K) map[K]V {
	res := make(map[K]V, len(arr))
	for _, item := range arr {
		key := fn(item)
		res[key] = item
	}
	return res
}
