package utils

func Map[R any, T any](list []R, transformer func(R) T) []T {
	items := make([]T, len(list))
	for i, element := range list {
		items[i] = transformer(element)
	}
	return items
}
