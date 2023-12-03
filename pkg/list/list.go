package list

type List[T comparable] []T

func New[T comparable](len int) List[T] {
	return make([]T, 0, len)
}

func (l List[T]) Find(v T) int {
	for i, d := range l {
		if d == v {
			return i
		}
	}
	return -1
}

func (l List[T]) RemoveAt(i int) (T, []T) {
	v := l[i]
	if i == 0 {
		nl := l[1:]
		return v, nl
	} else if i == len(l)-1 {
		nl := l[0 : len(l)-1]
		return v, nl
	}
	nl := l[0:i]
	nl = append(nl, l[i+1:]...)
	return v, nl
}

func (l List[T]) Add(i int, v T) []T {
	if i == 0 {
		nl := make([]T, len(l)+1)
		nl[0] = v
		copy(nl[1:], l)
		return nl
	} else if i == len(l) {
		nl := make([]T, len(l)+1)
		copy(nl, l)
		nl[len(nl)-1] = v
		return nl
	}
	nl := make([]T, len(l)+1)
	copy(nl, l[0:i])
	nl[i] = v
	copy(nl[i+1:], l[i:])
	return nl
}
