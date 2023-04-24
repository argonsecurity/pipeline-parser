package utils

func ToSlice[T any](v any) ([]T, bool) {
	anySlice, ok := v.([]any)
	if !ok {
		return []T{}, false
	}

	result := make([]T, len(anySlice))
	for i, item := range anySlice {
		result[i], ok = item.(T)
		if !ok {
			return []T{}, false
		}
	}
	return result, true
}

func Map[T any, U any](s []T, cb func(v T) U) []U {
	result := make([]U, len(s))
	for i, item := range s {
		result[i] = cb(item)
	}
	return result
}

func MapWithIndex[T any, U any](s []T, cb func(v T, i int) U) []U {
	result := make([]U, len(s))
	for i, item := range s {
		result[i] = cb(item, i)
	}
	return result
}

func Filter[T any](s []T, cb func(v T) bool) []T {
	result := make([]T, 0)
	for _, item := range s {
		if cb(item) {
			result = append(result, item)
		}
	}
	return result
}

func SliceContains[T comparable](s []T, v T) bool {
	for _, item := range s {
		if item == v {
			return true
		}
	}
	return false
}

func SliceToMap[T comparable, U any](s []T, cb func(v T) U) map[T]U {
	result := make(map[T]U)
	for _, item := range s {
		result[item] = cb(item)
	}
	return result
}

func SliceContainsBy[T comparable](s []T, v T, cb func(y, u T) bool) bool {
	for _, item := range s {
		if cb(item, v) {
			return true
		}
	}
	return false
}
