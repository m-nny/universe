package sliceutils

func Unique[T any, K comparable](arr []T, fn func(item T) K) []T {
	var result []T
	has := make(map[K]bool)
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

func Sum[T any](arr []T, fn func(item T) int) int {
	val := 0
	for _, item := range arr {
		val += fn(item)
	}
	return val
}
