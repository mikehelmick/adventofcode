package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type pos struct {
	x, y int
}

var offsets = []pos{{-1, -1}, {-1, 0}, {-1, 1}, {0, 1}, {1, 1}, {1, 0}, {1, -1}, {0, -1}}

func (p pos) add(o pos) pos {
	return pos{p.x + o.x, p.y + o.y}
}

func (p pos) String() string {
	return fmt.Sprintf("{%d, %d}", p.x, p.y)
}

type Grid [][]string

type Dict interface {
	IsWord(s string) bool
	IsPrefix(s string) bool
}

type XMASDict struct {
}

// The world's worst trie :)
func (x XMASDict) IsWord(s string) bool {
	return s == "XMAS"
}

func (x XMASDict) IsPrefix(s string) bool {
	return s == "X" || s == "XM" || s == "XMA"
}

type MASDict struct {
}

func (x MASDict) IsWord(s string) bool {
	return s == "MAS"
}

func (x MASDict) IsPrefix(s string) bool {
	return s == "M" || s == "MA"
}

func (g Grid) CountOuccrences() int {
	log := logging.DefaultLogger()
	count := 0
	for r, row := range g {
		for c := range row {
			for _, o := range offsets {
				log.Debugw("origin", "pos", pos{r, c}, "offset", o)
				count += g.search("", pos{r, c}, o, XMASDict{})
			}
		}
	}
	return count
}

func (g Grid) search(prefix string, p pos, dir pos, dict Dict) int {
	log := logging.DefaultLogger()

	candidate := prefix + g[p.x][p.y]
	if dict.IsWord(candidate) {
		log.Debugw("found", "candidate", candidate, "pos", p)
		return 1
	}
	if !dict.IsPrefix(candidate) {
		return 0
	}

	log.Debugw("searching", "prefix", prefix, "candidate", candidate, "pos", p)

	// on a valid prefix.
	next := p.add(dir)
	if next.x < 0 || next.x >= len(g) || next.y < 0 || next.y >= len(g[0]) {
		return 0
	}
	return g.search(candidate, next, dir, dict)
}

func (g Grid) findMas() int {
	// only the diagonals are valid
	offsets := []pos{{-1, -1}, {-1, 1}, {1, 1}, {1, -1}}

	log := logging.DefaultLogger()

	centers := make(map[pos]bool)
	count := 0
	for r, row := range g {
		for c := range row {
			for _, o := range offsets {
				p := pos{r, c}
				log.Debugw("origin", "pos", p, "offset", o)
				if (g.search("", p, o, MASDict{})) > 0 {
					center := p.add(o)
					if _, ok := centers[center]; !ok {
						centers[center] = true
					} else {
						count++
					}
				}
			}
		}
	}
	return count
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
	log.Infow("count", "count", grid.CountOuccrences())
	log.Infow("mascount", "count", grid.findMas())
}
