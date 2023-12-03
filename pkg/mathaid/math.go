package mathaid

func Min[K int | float32 | int64 | float64](a K, b K) K {
	if a <= b {
		return a
	}
	return b
}

func Max[K int | float32 | int64 | float64](a K, b K) K {
	if a >= b {
		return a
	}
	return b
}
