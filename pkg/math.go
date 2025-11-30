package utils

// Map function — pure functional helper
func Map[T any, R any](data []T, fn func(T) R) []R {
	result := make([]R, len(data))
	for i, v := range data {
		result[i] = fn(v)
	}
	return result
}
