package main

import (
	"bufio"
	"os"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/straid"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

type Grid [][]int64

func doDFS(grid Grid, start *twod.Pos) int {
	if grid[start.Row][start.Col] == 9 {
		return 1
	}

	val := grid[start.Row][start.Col]
	paths := 0
	candidates := start.Neighbors(
		func(p *twod.Pos) bool {
			return p.Row >= 0 && p.Row < len(grid) && p.Col >= 0 && p.Col < len(grid[0])
		})

	for _, cand := range candidates {
		if nextV := grid[cand.Row][cand.Col]; nextV == val+1 {
			paths += doDFS(grid, cand)
		}
	}

	return paths
}

func doBFS(grid Grid, start twod.Pos) int {
	waveFront := []twod.Pos{start}

	nines := make(map[twod.Pos]bool)

	for len(waveFront) > 0 {
		next := make(map[twod.Pos]bool)
		for _, pos := range waveFront {
			if grid[pos.Row][pos.Col] == 9 {
				nines[pos] = true
				continue
			}

			val := grid[pos.Row][pos.Col]
			candidates := pos.Neighbors(
				func(p *twod.Pos) bool {
					return p.Row >= 0 && p.Row < len(grid) && p.Col >= 0 && p.Col < len(grid[0])
				})
			for _, cand := range candidates {
				if nextV := grid[cand.Row][cand.Col]; nextV == val+1 {
					next[*cand] = true
				}
			}
		}
		waveFront = make([]twod.Pos, 0, len(next))
		for pos := range next {
			waveFront = append(waveFront, pos)
		}
	}

	return len(nines)
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	starts := make([]*twod.Pos, 0)
	grid := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]int64, 0, len(line))
		for i, c := range line {
			row = append(row, straid.AsInt(string(c)))
			if c == '0' {
				starts = append(starts, twod.NewPos(len(grid), i))
			}
		}
		grid = append(grid, row)
	}

	part1 := 0
	part2 := 0
	for _, start := range starts {
		part1 += doBFS(grid, *start)
		part2 += doDFS(grid, start)
	}

	log.Debugw("loaded", "starts", starts, "grid", grid)
	log.Infow("trailhead scores", "part1", part1)
	log.Infow("trailhead scores", "part2", part2)
}
