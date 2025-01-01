package gomap

func ComparableEqualFunc[T comparable](k1, k2 T) bool {
	return k1 == k2
}
