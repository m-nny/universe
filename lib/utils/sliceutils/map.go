package sliceutils

func ToMap[K comparable, V any](arr []V, fn func(item V) K) map[K]V {
	res := make(map[K]V, len(arr))
	for _, item := range arr {
		key := fn(item)
		res[key] = item
	}
	return res
}
