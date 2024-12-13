package testutil

func IgnoreError[T any](t T, err error) T {
	return t
}

func IgnoreError2[T1 any, T2 any](t1 T1, t2 T2, err error) (T1, T2) {
	return t1, t2
}
