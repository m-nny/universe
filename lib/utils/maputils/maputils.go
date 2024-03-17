package maputils

func Keys[K comparable, V any](m map[K]V) []K {
	var keys []K
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

func Values[K comparable, V any](m map[K]V) []V {
	var values []V
	for _, value := range m {
		values = append(values, value)
	}
	return values
}
