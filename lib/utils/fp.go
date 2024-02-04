package utils

func SliceMap[T, D any](arr []T, fn func(item T) D) []D {
	var res []D
	for _, item := range arr {
		res = append(res, fn(item))
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
