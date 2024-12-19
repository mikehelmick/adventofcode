package combinatorics

func AllPairs[T any](items []T) [][]T {
	pairs := make([][]T, 0, len(items)*(len(items)-1)/2)
	for i, a := range items {
		for _, b := range items[i+1:] {
			pairs = append(pairs, []T{a, b})
		}
	}
	return pairs
}
