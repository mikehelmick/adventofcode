package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

const PROCESSED = "#"

type Grid [][]string

func (g Grid) CalculateFence(r int, c int) (int, int) {
	validFn := func(p *twod.Pos) bool {
		return p.Row >= 0 && p.Row < len(g) && p.Col >= 0 && p.Col < len(g[p.Row])
	}

	areaType := g[r][c]
	allPoints := make(map[twod.Pos]bool)
	perimiterPoints := make(map[twod.Pos]bool)
	perimiterCount := 0

	wavefront := []twod.Pos{{Row: r, Col: c}}
	for len(wavefront) > 0 {
		next := make([]twod.Pos, 0)
		for _, p := range wavefront {
			if _, ok := allPoints[p]; ok {
				continue
			}
			allPoints[p] = true
			neighbors := p.ManhattanNeighbors(validFn)
			perimiterCount += 4 - len(neighbors)

			if len(neighbors) < 4 {
				validNeighbors := make(map[twod.Pos]bool)
				for _, n := range neighbors {
					validNeighbors[*n] = true
				}
				for _, offset := range twod.Manhattan {
					c := p.Clone()
					c.Add(offset)
					if _, ok := validNeighbors[*c]; !ok {
						perimiterPoints[*c] = true
					}
				}
			}

			for _, n := range neighbors {
				if g[n.Row][n.Col] == areaType {
					next = append(next, *n)
				} else {
					perimiterPoints[*n] = true
					perimiterCount++
				}
			}
		}
		wavefront = next
	}

	corners := 0
	// Count the corners.
	for p := range allPoints {
		// e i e corner
		if perimiterPoints[twod.Pos{Row: p.Row, Col: p.Col - 1}] && perimiterPoints[twod.Pos{Row: p.Row - 1, Col: p.Col}] {
			corners++
		}
		if perimiterPoints[twod.Pos{Row: p.Row, Col: p.Col - 1}] && perimiterPoints[twod.Pos{Row: p.Row + 1, Col: p.Col}] {
			corners++
		}
		if perimiterPoints[twod.Pos{Row: p.Row, Col: p.Col + 1}] && perimiterPoints[twod.Pos{Row: p.Row - 1, Col: p.Col}] {
			corners++
		}
		if perimiterPoints[twod.Pos{Row: p.Row, Col: p.Col + 1}] && perimiterPoints[twod.Pos{Row: p.Row + 1, Col: p.Col}] {
			corners++
		}
		// i i i e corner
		if allPoints[twod.Pos{Row: p.Row - 1, Col: p.Col}] && allPoints[twod.Pos{Row: p.Row, Col: p.Col + 1}] && perimiterPoints[twod.Pos{Row: p.Row - 1, Col: p.Col + 1}] {
			corners++
		}
		if allPoints[twod.Pos{Row: p.Row + 1, Col: p.Col}] && allPoints[twod.Pos{Row: p.Row, Col: p.Col + 1}] && perimiterPoints[twod.Pos{Row: p.Row + 1, Col: p.Col + 1}] {
			corners++
		}
		if allPoints[twod.Pos{Row: p.Row - 1, Col: p.Col}] && allPoints[twod.Pos{Row: p.Row, Col: p.Col - 1}] && perimiterPoints[twod.Pos{Row: p.Row - 1, Col: p.Col - 1}] {
			corners++
		}
		if allPoints[twod.Pos{Row: p.Row + 1, Col: p.Col}] && allPoints[twod.Pos{Row: p.Row, Col: p.Col - 1}] && perimiterPoints[twod.Pos{Row: p.Row + 1, Col: p.Col - 1}] {
			corners++
		}
	}

	// rewrite all points to be processed
	for p := range allPoints {
		g[p.Row][p.Col] = PROCESSED
	}
	log := logging.DefaultLogger()
	log.Debugw("DEBUG", "numPerimiter", perimiterCount, "corners", corners, "allPoints", len(allPoints))

	return perimiterCount * len(allPoints), corners * len(allPoints)
}

type Corner struct {
	Pos []twod.Pos
}

func (c *Corner) String() string {
	return fmt.Sprintf("%v,%v,%v", c.Pos[0], c.Pos[1], c.Pos[2])
}

func NewCorner(p []twod.Pos) *Corner {
	if len(p) != 3 {
		panic("invalid corner")
	}
	sort.Slice(p, func(i, j int) bool {
		if p[i].Row == p[j].Row {
			return p[i].Col < p[j].Col
		}
		return p[i].Row < p[j].Row
	})
	return &Corner{Pos: p}
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	grid := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]string, 0, len(line))
		for _, c := range line {
			row = append(row, string(c))
		}
		grid = append(grid, row)
	}
	log.Infow("loaded grid", "grid", grid)

	part1 := 0
	part2 := 0
	for r, row := range grid {
		for c := range row {
			if grid[r][c] == PROCESSED {
				continue // we've already processed this
			}
			log.Infow("processing", "r", r, "c", c, "val", grid[r][c])
			cost, bulkCost := grid.CalculateFence(r, c)
			part1 += cost
			part2 += bulkCost
			log.Infow("cost", "cost", cost, "bulkCost", bulkCost, "total", part1)
		}
	}

	log.Infow("part1", "part1", part1)
	log.Infow("part2", "part2", part2)
}
