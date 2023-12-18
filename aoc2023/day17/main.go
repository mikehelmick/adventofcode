package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/twod"
)

type Entry struct {
	P      Path
	Weight int
}

type PQ []Entry

func (h PQ) Len() int {
	return len(h)
}

func (h PQ) Less(i, j int) bool {
	return h[i].Weight < h[j].Weight
}

func (h PQ) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *PQ) Push(x any) {
	*h = append(*h, x.(Entry))
}

func (h *PQ) Pop() (x any) {
	x, *h = (*h)[len(*h)-1], (*h)[:len(*h)-1]
	return x
}

type Path struct {
	Pos *twod.Pos
	Dir *twod.Pos
}

func (p *Path) String() string {
	return fmt.Sprintf("%v-%v", p.Pos, p.Dir)
}

func search(g Grid, minDir int, maxDir int) int {
	isValid := func(p *twod.Pos) bool {
		return p.Row >= 0 && p.Col >= 0 && p.Row < len(g) && p.Col < len(g[0])
	}

	queue := &PQ{}
	visited := map[string]bool{}

	heap.Push(queue, Entry{P: Path{twod.NewPos(0, 0), twod.NewPos(1, 0)}, Weight: 0})
	heap.Push(queue, Entry{P: Path{twod.NewPos(0, 0), twod.NewPos(0, 1)}, Weight: 0})
	end := twod.NewPos(len(g)-1, len(g[0])-1)

	for queue.Len() > 0 {
		entry := heap.Pop(queue).(Entry)

		if entry.P.Pos.Equals(end) {
			return entry.Weight
		}
		if visited[entry.P.String()] {
			continue
		}
		visited[entry.P.String()] = true
		heat := entry.Weight

		for _, d := range []*twod.Pos{
			twod.NewPos(entry.P.Dir.Col, entry.P.Dir.Row),
			twod.NewPos(-1*entry.P.Dir.Col, -1*entry.P.Dir.Row),
		} {
			for i := minDir; i <= maxDir; i++ {
				nextPos := entry.P.Pos.Clone()
				// add the direction
				for t := 0; t < i; t++ {
					nextPos.Add(d)
				}

				if isValid(nextPos) {
					hl := 0
					for j := 1; j <= i; j++ {
						nextPos := entry.P.Pos.Clone()
						for t := 0; t < j; t++ {
							nextPos.Add(d)
						}
						hl += g[nextPos.Row][nextPos.Col]
					}
					heap.Push(queue, Entry{P: Path{nextPos, d}, Weight: heat + hl})
				}
			}
		}
	}
	return -1
}

type Grid [][]int

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	g := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "")
		row := make([]int, 0, len(parts))
		for _, p := range parts {
			v, err := strconv.Atoi(p)
			if err != nil {
				panic(err)
			}
			row = append(row, v)
		}
		g = append(g, row)
	}

	minLoss := search(g, 1, 3)
	fmt.Printf("part1 %+v\n", minLoss)

	part2 := search(g, 4, 10)
	fmt.Printf("part2 %+v\n", part2)
}
