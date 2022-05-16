package utils

func GetPtr[T any](v T) *T {
	return &v
}

func GetPtrOrNil[T comparable](v T) *T {
	var zeroValue T
	if zeroValue == v {
		return nil
	}

	return &v
}
