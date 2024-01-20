package utils

func SliceMap[T, D any](arr []T, fn func(item T) D) []D {
	var res []D
	for _, item := range arr {
		res = append(res, fn(item))
	}
	return res
}
