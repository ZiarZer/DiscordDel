package utils

func MakePointer[T any](value T) *T {
	return &value
}
