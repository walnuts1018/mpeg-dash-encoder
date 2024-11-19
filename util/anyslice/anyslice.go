package anyslice

import (
	"fmt"
)

func FromAny[T any](slice []any) ([]T, error) {
	results := make([]T, len(slice))
	for i, v := range slice {
		tmp, ok := v.(T)
		if !ok {
			return nil, fmt.Errorf("failed to parse %T: %#v", v, v)
		}
		results[i] = tmp
	}
	return results, nil
}

func ToAny[T any](slice []T) []any {
	results := make([]any, len(slice))
	for i, v := range slice {
		results[i] = v
	}
	return results
}
