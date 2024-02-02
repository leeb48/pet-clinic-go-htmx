package utils

func Filter[T comparable](arr []T, test func(T) bool) []T {

	filtered := []T{}

	for _, ele := range arr {
		if test(ele) {
			filtered = append(filtered, ele)
		}
	}
	return filtered
}
