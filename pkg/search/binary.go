package search

type Test[T int | int64] func(median T) T

func BinarySearch[T int | int64](low, high T, check Test[T]) T {
	for low <= high {
		median := (low + high) / 2

		d := check(median)
		if d == 0 {
			return median
		}
		if d > 0 {
			low = median + 1
		} else {
			high = median - 1
		}
	}
	return -1
}
