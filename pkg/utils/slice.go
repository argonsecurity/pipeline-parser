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
