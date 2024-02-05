package utils

import "context"

func SliceMapCtxErr[T, D any](ctx context.Context, arr []T, fn func(ctx context.Context, item T) (D, error)) ([]D, error) {
	var res []D
	for _, item := range arr {
		val, err := fn(ctx, item)
		if err != nil {
			return nil, err
		}
		res = append(res, val)
	}
	return res, nil
}

func SliceMap[T, D any](arr []T, fn func(item T) D) []D {
	var res []D
	for _, item := range arr {
		res = append(res, fn(item))
	}
	return res
}

func SliceFlatMap[T, D any](arr []T, fn func(item T) []D) []D {
	var res []D
	for _, item := range arr {
		res = append(res, fn(item)...)
	}
	return res
}

func SliceUniqe[T any](arr []T, fn func(item T) string) []T {
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
