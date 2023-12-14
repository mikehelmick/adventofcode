package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
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

func (g Grid) TiltNorth() {
	for r := 1; r < len(g); r++ {
		for c := 0; c < len(g[0]); c++ {
			if g[r][c] == "O" {
				for cr := r; cr >= 1; cr-- {
					if g[cr-1][c] == "." {
						g[cr-1][c] = "O"
						g[cr][c] = "."
						continue
					}
					break
				}
			}
		}
	}
}

func (g Grid) TiltWest() {
	for c := 0; c < len(g[0]); c++ {
		for r := 0; r < len(g); r++ {
			if g[r][c] == "O" {
				for rc := c; rc >= 1; rc-- {
					if g[r][rc-1] == "." {
						g[r][rc-1] = "O"
						g[r][rc] = "."
						continue
					}
					break
				}
			}
		}
	}
}

func (g Grid) TiltSouth() {
	for r := len(g) - 1; r >= 0; r-- {
		for c := 0; c < len(g[0]); c++ {
			if g[r][c] == "O" {
				for cr := r; cr < len(g)-1; cr++ {
					if g[cr+1][c] == "." {
						g[cr+1][c] = "O"
						g[cr][c] = "."
						continue
					}
					break
				}
			}
		}
	}
}

func (g Grid) TiltEast() {
	for c := len(g[0]) - 1; c >= 0; c-- {
		for r := 0; r < len(g); r++ {
			if g[r][c] == "O" {
				for rc := c; rc < len(g[0])-1; rc++ {
					if g[r][rc+1] == "." {
						g[r][rc+1] = "O"
						g[r][rc] = "."
						continue
					}
					break
				}
			}
		}
	}
}

func (g Grid) Cycle() {
	g.TiltNorth()
	g.TiltWest()
	g.TiltSouth()
	g.TiltEast()
}

func (g Grid) Weight() int {
	l := len(g)
	weight := 0
	for r, row := range g {
		for _, v := range row {
			if v == "O" {
				weight += (l - r)
			}
		}
	}
	return weight
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	grid := make(Grid, 0)
	grid2 := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		row := strings.Split(line, "")
		grid = append(grid, row)
		// second grid since first one gets a single north tilt
		grid2 = append(grid2, row)
	}

	grid.TiltNorth()
	fmt.Printf("answer part 1: %v\n", grid.Weight())

	grid = nil

	cycles := make(map[string]int)
	cycles[grid2.String()] = 0

	initial := 0
	cycleWeight := 0
	// assume there will be a cycle before 10k
	for i := 1; i <= 10000; i++ {
		grid2.Cycle()
		key := grid2.String()
		if v, ok := cycles[key]; ok {
			fmt.Printf("MATCH: %+v to %+v cycles\n", v, i)
			initial = v
			cycleWeight = i - v
			break
		}
		cycles[key] = i
	}
	// cycle has been found - subtract the items before the first cycle
	// and the mod of the cycle length is how many cycles to 1B.
	toDo := (1_000_000_000 - initial) % cycleWeight
	for i := 0; i < toDo; i++ {
		grid2.Cycle()
	}

	fmt.Printf("part2: %+v\n", grid2.Weight())
}
