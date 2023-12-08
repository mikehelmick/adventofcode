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

func GreatestCommonDivisor(a, b int64) int64 {
	for b != 0 {
		t := b
		b = a % b
		a = t
	}
	return a
}

func LowestCommonMultiple(a, b int64, integers ...int64) int64 {
	result := a * b / GreatestCommonDivisor(a, b)
	for i := 0; i < len(integers); i++ {
		result = LowestCommonMultiple(result, integers[i])
	}
	return result
}
