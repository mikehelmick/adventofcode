package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/mathaid"
	. "github.com/mikehelmick/adventofcode/pkg/twod"
)

var connections = map[string][]*Pos{
	"|": {NewPos(-1, 0), NewPos(1, 0)},
	"-": {NewPos(0, -1), NewPos(0, 1)},
	"L": {NewPos(-1, 0), NewPos(0, 1)},
	"J": {NewPos(-1, 0), NewPos(0, -1)},
	"7": {NewPos(1, 0), NewPos(0, -1)},
	"F": {NewPos(1, 0), NewPos(0, 1)},
	".": {},
	"S": {NewPos(-1, 0), NewPos(0, -1)}, // my starting input is a J
	// You'd have to change this to the correct character for your input or examples.
}

func findStart(grid []string) *Pos {
	for r, row := range grid {
		if c := strings.Index(row, "S"); c >= 0 {
			return NewPos(r, c)
		}
	}
	panic("no start found")
}

// does a BFS from the starting point to find the farthest point in the loop.
func findFurthest(grid []string, dist [][]int) int {
	validFunc := func(p *Pos) bool {
		return p.Row >= 0 && p.Col >= 0 &&
			p.Row < len(grid) && p.Col < len(grid[0])
	}

	start := findStart(grid)

	maxSetDist := 0
	dist[start.Row][start.Col] = 0
	distance := 0
	wave := []*Pos{start}

	for len(wave) > 0 {
		distance++
		next := make([]*Pos, 0)
		for _, from := range wave {
			for _, cand := range from.Follow(validFunc, connections[grid[from.Row][from.Col:from.Col+1]]) {
				if dist[cand.Row][cand.Col] < 0 {
					dist[cand.Row][cand.Col] = distance
					maxSetDist = mathaid.Max(maxSetDist, distance)
					next = append(next, cand.Clone())
				}
			}
		}
		wave = next
	}

	return maxSetDist
}

func isInsideShape(r, c int, grid []string) bool {
	char := grid[r][c : c+1]
	if char == "|" || char == "J" || char == "L" || char == "S" {
		return true
	}
	return false
}

// Find the number of tiles that are NOT part of the loop (from part 1)
// a tile is inside if it has an odd number of vertical, J or L next to them (S is a J in my input).
func countInsides(grid []string, dist [][]int) int {
	insides := 0
	for r, row := range dist {
		insideShapes := 0
		for c := range row {
			if dist[r][c] >= 0 {
				// part of the loop
				if isInsideShape(r, c, grid) {
					insideShapes++
				}
				continue
			}
			if insideShapes%2 == 1 {
				insides++
			}
		}
	}
	return insides
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	grid := make([]string, 0)
	dist := make([][]int, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		grid = append(grid, line)
		row := make([]int, len(line))
		for i := range row {
			row[i] = -1
		}
		dist = append(dist, row)
	}
	log.Debugw("loaded grid", "grid", grid, "dist", dist)

	log.Infow("answer", "part1", findFurthest(grid, dist))

	log.Infow("answer", "part2", countInsides(grid, dist))

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
