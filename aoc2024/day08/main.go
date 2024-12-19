package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mikehelmick/adventofcode/pkg/combinatorics"
	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

type Grid [][]string

func (g Grid) String() string {
	str := ""
	for _, row := range g {
		for _, c := range row {
			str += c
		}
		str += "\n"
	}
	return str
}

func (g Grid) Copy() Grid {
	c := make(Grid, len(g))
	for i, row := range g {
		c[i] = make([]string, len(row))
		copy(c[i], row)
	}
	return c
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	antennae := make(map[string][]*twod.Pos)
	grid := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		row := make([]string, 0, len(line))
		for i, c := range line {
			row = append(row, string(c))

			if c != '.' {
				pos := twod.NewPos(len(grid), i)
				if _, ok := antennae[string(c)]; !ok {
					antennae[string(c)] = make([]*twod.Pos, 0, 2)
				}
				antennae[string(c)] = append(antennae[string(c)], pos)
			}
		}
		grid = append(grid, row)
	}
	log.Debugw("antennae", "antennae", antennae)

	p2Grid := grid.Copy()
	// part1
	{
		antinodes := make(map[twod.Pos]bool)
		for _, locs := range antennae {
			pairs := combinatorics.AllPairs(locs)
			for _, pair := range pairs {
				slopeRow := pair[0].Row - pair[1].Row
				slopeCol := pair[0].Col - pair[1].Col

				antinode := twod.NewPos(pair[0].Row+slopeRow, pair[0].Col+slopeCol)
				if isValid(antinode, grid) {
					grid[antinode.Row][antinode.Col] = "#"
					antinodes[*antinode] = true
				}
				antinode = twod.NewPos(pair[1].Row-slopeRow, pair[1].Col-slopeCol)
				if isValid(antinode, grid) {
					grid[antinode.Row][antinode.Col] = "#"
					antinodes[*antinode] = true
				}
			}
		}
		fmt.Printf("after:\n%s\n", grid.String())
		log.Infow("part1", "antinodes", len(antinodes))
	}

	// part2
	{
		grid = p2Grid
		antinodes := make(map[twod.Pos]bool)
		for _, locs := range antennae {
			pairs := combinatorics.AllPairs(locs)
			for _, pair := range pairs {
				slopeRow := pair[0].Row - pair[1].Row
				slopeCol := pair[0].Col - pair[1].Col

				antinodes[*pair[0]] = true
				antinodes[*pair[1]] = true

				antinode := twod.NewPos(pair[0].Row+slopeRow, pair[0].Col+slopeCol)
				for isValid(antinode, grid) {
					grid[antinode.Row][antinode.Col] = "#"
					antinodes[*antinode] = true
					antinode = twod.NewPos(antinode.Row+slopeRow, antinode.Col+slopeCol)
				}
				antinode = twod.NewPos(pair[1].Row-slopeRow, pair[1].Col-slopeCol)
				for isValid(antinode, grid) {
					grid[antinode.Row][antinode.Col] = "#"
					antinodes[*antinode] = true
					antinode = twod.NewPos(antinode.Row-slopeRow, antinode.Col-slopeCol)
				}
			}
		}
		fmt.Printf("after:\n%s\n", grid.String())
		log.Infow("part2", "antinodes", len(antinodes))
	}
}

func isValid(p *twod.Pos, grid Grid) bool {
	return p.Row >= 0 && p.Row < len(grid) && p.Col >= 0 && p.Col < len(grid[0])
}
