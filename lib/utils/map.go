package utils

import "slices"

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
