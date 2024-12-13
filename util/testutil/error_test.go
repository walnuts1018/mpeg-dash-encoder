package testutil

import (
	"reflect"
	"testing"
)

func TestIgnoreError(t *testing.T) {
	t.Run("func() (T, error)", func(t *testing.T) {
		want := 1

		testFunc := func() (int, error) {
			return want, nil
		}

		if got := IgnoreError(testFunc()); !reflect.DeepEqual(got, want) {
			t.Errorf("IgnoreError() = %v, want %v", got, want)
		}
	})

	t.Run("func() (T1, T2, error)", func(t *testing.T) {
		want := 1

		testFunc := func() (int, int, error) {
			return want, want, nil
		}

		if got1, got2 := IgnoreError2(testFunc()); !reflect.DeepEqual(got1, want) || !reflect.DeepEqual(got2, want) {
			t.Errorf("IgnoreError2() = %v, %v, want %v, %v", got1, got2, want, want)
		}
	})
}
