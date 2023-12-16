package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
	"github.com/mikehelmick/go-functional/slice"
)

var (
	// helpers for reflections.
	leftMirror = map[string]string{
		twod.RIGHT: twod.DOWN, twod.DOWN: twod.RIGHT, twod.UP: twod.LEFT, twod.LEFT: twod.UP,
	}
	rightMirror = map[string]string{
		twod.RIGHT: twod.UP, twod.DOWN: twod.LEFT, twod.UP: twod.RIGHT, twod.LEFT: twod.DOWN,
	}
)

type Grid [][]string

type Light struct {
	Position  *twod.Pos
	Direction string
}

func (l *Light) Key() string {
	return fmt.Sprintf("%s-%s", l.Position.String(), l.Direction)
}

func (l *Light) Split(a, b string) (*Light, *Light) {
	return &Light{Position: l.Position.Clone(), Direction: a}, &Light{Position: l.Position.Clone(), Direction: b}
}

func (l *Light) Move() *Light {
	next := l.Position.Clone()
	next.Add(twod.Dirs[l.Direction])
	return &Light{
		Position:  next,
		Direction: l.Direction,
	}
}

func Print(g Grid, wf []*Light, e [][]bool) {
	m := make(map[string]*Light)
	for _, k := range wf {
		m[k.Position.String()] = k
	}

	for r, row := range g {
		for c, col := range row {
			if l, ok := m[fmt.Sprintf("{%v,%v}", r, c)]; ok {
				fmt.Printf(l.Direction)
			} else if e[r][c] {
				fmt.Printf("#")
			} else {
				fmt.Print(col)
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n**********\n")
}

func shootLasers(g Grid, e map[string]bool, start *Light) {
	isValid := func(l *Light) bool {
		p := l.Position
		return p.Row >= 0 && p.Col >= 0 && p.Row < len(g) && p.Col < len(g[0])
	}

	waveFront := []*Light{start}
	// This is for termination, the energized map is just energized.
	visitDirections := make(map[string]bool)

	for len(waveFront) > 0 {
		// for cool animations, uncomment next line.
		// Print(g, waveFront, e)
		nextWave := make([]*Light, 0)
		for _, l := range waveFront {
			// mark space visited
			e[l.Position.String()] = true
			visitDirections[l.Key()] = true
			space := g[l.Position.Row][l.Position.Col]
			switch space {
			case ".":
				nextWave = append(nextWave, l.Move())
			case "/":
				l.Direction = rightMirror[l.Direction]
				nextWave = append(nextWave, l.Move())
			case "\\":
				l.Direction = leftMirror[l.Direction]
				nextWave = append(nextWave, l.Move())
			case "|":
				if l.Direction == twod.LEFT || l.Direction == twod.RIGHT {
					a, b := l.Split(twod.DOWN, twod.UP)
					nextWave = append(nextWave, a.Move())
					nextWave = append(nextWave, b.Move())
				} else {
					nextWave = append(nextWave, l.Move())
				}
			case "-":
				if l.Direction == twod.UP || l.Direction == twod.DOWN {
					a, b := l.Split(twod.LEFT, twod.RIGHT)
					nextWave = append(nextWave, a.Move())
					nextWave = append(nextWave, b.Move())
				} else {
					nextWave = append(nextWave, l.Move())
				}
			default:
				panic(fmt.Sprintf("unexpected character: %q", space))
			}
		}
		// filter out lights that are now off the board
		// and filter out lights that we've seen before (Same space & direction)
		waveFront = slice.Filter(slice.Filter(nextWave, isValid),
			func(l *Light) bool {
				return !visitDirections[l.Key()]
			})
	}
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	g := make(Grid, 0)
	e := make(map[string]bool)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		row := strings.Split(line, "")
		g = append(g, row)
	}

	shootLasers(g, e, &Light{twod.NewPos(0, 0), twod.RIGHT})
	part1 := len(e)
	fmt.Printf("part1: %v\n", part1)

	starting := make([]*Light, 0)
	for r := 0; r < len(g); r++ {
		if r == 0 {
			for c := 0; c < len(g[0]); c++ {
				starting = append(starting, &Light{twod.NewPos(r, c), twod.DOWN})
			}
			starting = append(starting, &Light{twod.NewPos(r, 0), twod.RIGHT})
			starting = append(starting, &Light{twod.NewPos(r, len(g[0])-1), twod.LEFT})
		} else if r == len(g)-1 {
			for c := 0; c < len(g[0]); c++ {
				starting = append(starting, &Light{twod.NewPos(r, c), twod.UP})
			}
			starting = append(starting, &Light{twod.NewPos(r, 0), twod.RIGHT})
			starting = append(starting, &Light{twod.NewPos(r, len(g[0])-1), twod.LEFT})
		} else {
			starting = append(starting, &Light{twod.NewPos(r, 0), twod.RIGHT})
			starting = append(starting, &Light{twod.NewPos(r, len(g[r])-1), twod.LEFT})
		}
	}
	part2 := 0
	for _, s := range starting {
		e := make(map[string]bool)
		shootLasers(g, e, s)
		part2 = max(part2, len(e))
	}
	fmt.Printf("part2: %v\n", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
