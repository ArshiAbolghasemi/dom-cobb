package utils

func ToAnySlice[T any](entries []T) []any {
	interfaces := make([]any, len(entries))
	for i, v := range entries {
		interfaces[i] = v
	}
	return interfaces
}
