package slices

func Copy[T any](ar []T) []T {
	ret := make([]T, len(ar))
	copy(ret, ar)
	return ret
}
