package main

import (
	"bufio"
	"context"
	"os"
	"strconv"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

var (
	digits = map[string]bool{
		"0": true, "1": true, "2": true, "3": true, "4": true, "5": true, "6": true, "7": true, "8": true, "9": true,
	}
)

type Number struct {
	Pos    *twod.Pos
	Length int
	Value  int64
}

// Points returns all of the points that make up this number.
func (n *Number) Points() []*twod.Pos {
	rtn := make([]*twod.Pos, 0, n.Length)
	for i := 0; i < n.Length; i++ {
		rtn = append(rtn, twod.NewPos(n.Pos.Row, n.Pos.Col+i))
	}
	return rtn
}

// Returns all of the stars (as points) that are adjacent to this number
func (n *Number) StarAdjacent(g []string) []twod.Pos {
	points := n.Points()
	vFunc := ValidFunc(g)

	stars := make(map[twod.Pos]bool)
	for _, p := range points {
		for _, adj := range p.Adjacent(vFunc) {
			char := g[adj.Row][adj.Col : adj.Col+1]
			if char == "*" {
				stars[*adj] = true
			}
		}
	}
	rtn := make([]twod.Pos, 0, len(stars))
	for p := range stars {
		rtn = append(rtn, p)
	}
	return rtn
}

// SymbolAdjacent returns true if this number is adjacent
// to a symbol (non number, non .) on the grid.
func (n *Number) SymbolAdjacent(g []string) bool {
	points := n.Points()
	vFunc := ValidFunc(g)

	for _, p := range points {
		for _, adj := range p.Adjacent(vFunc) {
			char := g[adj.Row][adj.Col : adj.Col+1]
			if char != "." && !digits[char] {
				// must be a symbol
				return true
			}
		}
	}
	return false
}

// FindNumbers find all of the numbers in the grid.
func FindNumbers(g []string) []Number {
	rtn := make([]Number, 0)
	for row, line := range g {
		for col := 0; col < len(line); col++ {
			if digits[g[row][col:col+1]] {
				// start of a number
				num := Number{Pos: twod.NewPos(row, col), Length: 0}
				for col < len(line) && digits[g[row][col:col+1]] {
					col++
					num.Length++
				}
				var err error
				num.Value, err = strconv.ParseInt(g[row][num.Pos.Col:num.Pos.Col+num.Length], 10, 64)
				if err != nil {
					panic(err)
				}
				rtn = append(rtn, num)
			}
		}
	}
	return rtn
}

// ValidFunc creates a twod.ValidFunc for checking valid points on the grid.
func ValidFunc(grid []string) twod.ValidFunc {
	rows := len(grid)
	cols := len(grid[0])
	return func(p *twod.Pos) bool {
		return p.Row >= 0 && p.Row < rows && p.Col >= 0 && p.Col < cols
	}
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	grid := make([]string, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		grid = append(grid, line)
		log.Debugw("parsed line", "line", line)
	}

	numbers := FindNumbers(grid)
	log.Debugw("found numbers", "numbers", numbers)

	var part1 int64
	// Just add up all the numbers that are symbol adjacent.
	for _, n := range numbers {
		if n.SymbolAdjacent(grid) {
			log.Debugw("symbol adjacent", "row", n.Pos.Row, "col", n.Pos.Col, "value", n.Value)
			part1 += n.Value
		}
	}
	log.Infow("part 1", "answer", part1)

	// part 2
	// Create a map of all the stars to the numbers they are adjacent to.
	starMap := make(map[twod.Pos][]Number)
	// Do this by going over every number
	for _, n := range numbers {
		// and finding all the stars it is adjacent to (could be more than one)
		stars := n.StarAdjacent(grid)
		for _, s := range stars {
			cur, ok := starMap[s]
			if !ok {
				cur = make([]Number, 0, 1)
			}
			cur = append(cur, n)
			starMap[s] = cur
		}
	}
	log.Debug("starmap", "stars", starMap)

	var part2 int64
	// For all the stars that are adjacent to exactly 2 numbers.
	for _, adjNum := range starMap {
		if len(adjNum) != 2 {
			continue
		}
		// Multiply the numbers and add them to the total.
		ratio := adjNum[0].Value * adjNum[1].Value
		part2 += ratio
	}
	log.Infow("part 2", "answer", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
