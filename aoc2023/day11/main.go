package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
	"gonum.org/v1/gonum/stat/combin"
)

type Grid [][]string

func (g Grid) String() string {
	b := strings.Builder{}
	for _, r := range g {
		b.WriteString(fmt.Sprintf("%+v\n", r))
	}
	return b.String()
}

func (g Grid) Points() []*twod.Pos {
	points := make([]*twod.Pos, 0)
	for r, row := range g {
		for c, char := range row {
			if char == "#" {
				points = append(points, twod.NewPos(r, c))
			}
		}
	}
	return points
}

func isEmpty(r []string) bool {
	for _, s := range r {
		if s == "#" {
			return false
		}
	}
	return true
}

func (g Grid) EmptyCols() []int {
	emptyCols := make([]int, 0)
	for col := 0; col < len(g[0]); col++ {
		isEmpty := true
		for row := 0; row < len(g) && isEmpty; row++ {
			isEmpty = g[row][col] != "#"
		}
		if isEmpty {
			emptyCols = append(emptyCols, col)
		}
	}
	return emptyCols
}

func (g Grid) EmptyRows() []int {
	r := make([]int, 0)
	for rN, row := range g {
		if isEmpty(row) {
			r = append(r, rN)
		}
	}
	return r
}

func expand(factor int, points []*twod.Pos, emptyRows, emptyCols []int) []*twod.Pos {
	newPoints := make([]*twod.Pos, 0, len(points))
	for _, p := range points {
		addR := 0
		addC := 0
		for _, er := range emptyRows {
			if p.Row > er {
				addR += factor
			}
		}
		for _, ec := range emptyCols {
			if p.Col > ec {
				addC += factor
			}
		}
		newPoints = append(newPoints, twod.NewPos(p.Row+addR, p.Col+addC))
	}
	return newPoints
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	grid := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		chars := strings.Split(line, "")
		grid = append(grid, chars)
	}
	log.Debugw("loaded", "grid", grid)

	points := grid.Points()
	log.Debugw("found points", "points", points)

	emptyRows := grid.EmptyRows()
	emptyCols := grid.EmptyCols()

	part1Points := expand(1, points, emptyRows, emptyCols)
	log.Debugw("found points", "points", part1Points)

	pairs := combin.Combinations(len(part1Points), 2)
	log.Debug("calculated pairs", "pairs", pairs)

	allDist := 0
	for _, pair := range pairs {
		allDist += part1Points[pair[0]].Dist(part1Points[pair[1]])
	}
	log.Infow("answer", "part1", allDist)

	part2Points := expand(999999, points, emptyRows, emptyCols)
	log.Debugw("found points", "part2points", part2Points)
	p2answer := int64(0)
	for _, pair := range pairs {
		p2answer += int64(part2Points[pair[0]].Dist(part2Points[pair[1]]))
	}
	log.Infow("answer", "part2", p2answer)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
