package sliceequal

import (
	"reflect"
	"slices"
)

func Equal[S ~[]E, E comparable](s1, s2 S, cmp func(a, b E) int) bool {
	if len(s1) != len(s2) {
		return false
	}

	sortedS1 := append(S{}, s1...)
	sortedS2 := append(S{}, s2...)
	slices.SortFunc(sortedS1, cmp)
	slices.SortFunc(sortedS2, cmp)

	return reflect.DeepEqual(sortedS1, sortedS2)
}
