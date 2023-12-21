package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

type Grid [][]string

func (g Grid) String() string {
	b := strings.Builder{}
	for _, row := range g {
		b.WriteString(strings.Join(row, ""))
		b.WriteString("\n")
	}
	return b.String()
}

func (g Grid) FindStart() *twod.Pos {
	for r, row := range g {
		for c, v := range row {
			if v == "S" {
				return twod.NewPos(r, c)
			}
		}
	}
	panic("no start")
}

func (g Grid) GetPoint(row int, col int) string {
	row = row % len(g)
	if row < 0 {
		row = len(g) + row // add negative num
	}
	col = col % len(g[0])
	if col < 0 {
		col = len(g[0]) + col
	}
	return g[row][col]
}

func (g Grid) BFS(s *twod.Pos, steps int, infinite bool) (int, []int64) {
	output := make([]int64, 0)
	log := logging.DefaultLogger()
	isValid := func(p *twod.Pos) bool {
		return p.Row >= 0 && p.Col >= 0 && p.Row < len(g) && p.Col < len(g[0])
	}
	if infinite {
		isValid = func(p *twod.Pos) bool { return true }
	}

	visited := make(map[string]bool)
	wave := []*twod.Pos{s}
	for i := 0; i < steps && len(wave) > 0; i++ {
		visited = make(map[string]bool)
		next := make([]*twod.Pos, 0)
		for _, w := range wave {
			for _, n := range w.Neighbors(isValid) {
				if v := g.GetPoint(n.Row, n.Col); (v == "." || v == "S") && !visited[n.String()] {
					next = append(next, n)
					visited[n.String()] = true
				}
			}
		}
		wave = next
		if (i+i == 65) || (i-65+1)%131 == 0 {
			output = append(output, int64(len(visited)))
			log.Debugw("wave", "steps", i+1, "count", len(visited))
		}
	}

	return len(visited), output
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	g := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		row := strings.Split(line, "")
		g = append(g, row)
	}

	part1, _ := g.BFS(g.FindStart(), 64, false)
	log.Infow("answer", "part1", part1)

	/*
		  // Used to get the input values for the quadratic formula.
			_, output := g.BFS(g.FindStart(), 328, true)
	*/

	log.Infow("answer", "part2", part2(26501365/131, 3725, 32896, 91055))
}

func part2(goal uint64, a0, a1, a2 int64) uint64 {
	b0 := uint64(a0)
	b1 := uint64(a1 - a0)
	b2 := uint64(a2 - a1)
	return b0 + b1*goal + (goal*(goal-1)/2)*(b2-b1)
}
