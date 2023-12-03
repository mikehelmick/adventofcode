package list_test

import (
	"testing"

	"github.com/mikehelmick/adventofcode/pkg/list"
)

func TestMap(t *testing.T) {

	input := []int{1, 2, 3, 4, 5}
	out := list.Map(input, func(i int) int { return i * 2 })

	for i, v := range input {
		if v*2 != out[i] {
			t.Errorf("elem %v is wrong, want: %v got: %v", i, v*2, out[i])
		}
	}
}
