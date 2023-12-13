package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
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

func rowEquals(a, b []string) bool {
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func opposite(s string) string {
	if s == "." {
		return "#"
	}
	return "."
}
func (g Grid) SmudgeValue(cannot int) int {
	skipH := -1
	skipV := -1
	if cannot >= 100 {
		skipH = cannot / 100
	} else {
		skipV = cannot
	}

	// horizontal smudges - swap each element and see if we can find
	// a different mirror line.
	for r := 0; r < len(g); r++ {
		for c := 0; c < len(g[r]); c++ {
			g[r][c] = opposite(g[r][c])
			if hm := g.HorizontalMirror(skipH); hm > 0 {
				return hm * 100
			}
			g[r][c] = opposite(g[r][c])
		}
	}

	ng := g.Transpose()
	// transpose once and then swap each element.
	for r := 0; r < len(ng); r++ {
		for c := 0; c < len(ng[r]); c++ {
			ng[r][c] = opposite(ng[r][c])
			if hm := ng.HorizontalMirror(skipV); hm > 0 {
				return hm
			}
			ng[r][c] = opposite(ng[r][c])
		}
	}

	panic("mirror cannot find smudge")
}

// Just make it so I only have to write horizontal scanner. Lazy.
func (g Grid) Transpose() Grid {
	newG := make(Grid, 0)
	for col := 0; col < len(g[0]); col++ {
		newRow := make([]string, 0)
		for r := 0; r < len(g); r++ {
			newRow = append(newRow, g[r][col])
		}
		newG = append(newG, newRow)
	}
	return newG
}

// find the rows above a horizontal mirror, -1 if can't be found.
// at least one (top or bottom) must be fully covered.
func (g Grid) HorizontalMirror(skip int) int {
	log := logging.DefaultLogger()
	// r is "before" row 1 (Between 0 / 1)
	for r := 1; r < len(g); r++ {
		if r == skip {
			continue
		}
		numAbove := r
		numBelow := len(g) - r
		// how many rows must be the same
		same := min(numAbove, numBelow)

		allSame := true
		for offset := 0; offset < same; offset++ {
			log.Debugw("comparing", "a", g[r-offset-1], "b", g[r+offset])
			if !rowEquals(g[r-offset-1], g[r+offset]) {
				allSame = false
				break
			}
		}
		if allSame {
			return r
		}
	}
	return -1
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	grids := make([]Grid, 0)
	cur := make(Grid, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			grids = append(grids, cur)
			cur = make(Grid, 0)
			continue
		}

		row := strings.Split(line, "")
		cur = append(cur, row)
	}
	if len(cur) > 0 {
		grids = append(grids, cur)
	}

	part1 := 0
	values := make([]int, 0)
	for i, g := range grids {

		gV := g.Transpose()
		if vert := gV.HorizontalMirror(-1); vert > 0 {
			log.Infow("vertical", "mirror", i, "toLeft", vert)
			values = append(values, vert)
			part1 += vert
			continue
		}
		// must be horizontal
		horiz := g.HorizontalMirror(-1)
		log.Infow("horizontal", "mirror", i, "above", horiz)
		if horiz < 0 {
			panic("didn't find a mirror " + g.String())
		}
		values = append(values, 100*horiz)
		part1 += (100 * horiz)
	}
	log.Infow("answer", "part1", part1)

	part2 := 0
	for i, g := range grids {
		sv := g.SmudgeValue(values[i])
		log.Infow("part2", "mirror", i, "value", sv)
		part2 += sv
	}
	log.Infow("answer", "part2", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
