package collections

func ContainsFn[T any](items []T, fn func(item T) bool) bool {
	for _, item := range items {
		if fn(item) {
			return true
		}
	}

	return false
}
