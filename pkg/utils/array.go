package utils

func IsArray(input any) bool {
	_, ok := input.([]any)
	return ok
}
