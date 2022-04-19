package utils

func GetMapKeys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}

func MapToSlice[T any, U any, K comparable](s map[K]T, cb func(k K, v T) U) []U {
	result := make([]U, len(s))
	var i int
	for k, v := range s {
		result[i] = cb(k, v)
		i++
	}
	return result
}
