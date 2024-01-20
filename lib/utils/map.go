package utils

import "slices"

func SliceMap[T, D any](arr []T, fn func(item T) D) []D {
	var res []D
	for _, item := range arr {
		res = append(res, fn(item))
	}
	return res
}

func MapKeys[K comparable, V comparable](m map[K]V) []K {
	var res []K
	for key := range m {
		res = append(res, key)
	}
	return res
}

func TopValuesMap[T comparable](m map[T]int) []T {
	keys := MapKeys(m)
	slices.SortFunc(keys, func(a, b T) int {
		return m[b] - m[a]
	})
	return keys
}
