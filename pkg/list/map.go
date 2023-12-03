package list

func Map[T any, O any](in []T, f func(T) O) []O {
	out := make([]O, len(in))
	for i, v := range in {
		out[i] = f(v)
	}
	return out
}
