package utils

func GetMapKeys[T comparable, U any](m map[T]U) []T {
	keys := make([]T, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i += 1
	}
	return keys
}

func MapToSlice[T any, U any, K comparable](m map[K]T, cb func(k K, v T) U) []U {
	result := make([]U, len(m))
	var i int
	for k, v := range m {
		item := cb(k, v)
		result[i] = item
		i += 1
	}
	return result
}

func MapToSliceErr[T any, U any, K comparable](m map[K]T, cb func(k K, v T) (U, error)) ([]U, error) {
	result := make([]U, len(m))
	var i int
	for k, v := range m {
		item, err := cb(k, v)
		if err != nil {
			return nil, err
		}
		result[i] = item
		i += 1
	}
	return result, nil
}
