package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
	"github.com/mikehelmick/go-functional/slice"
)

const (
	WALL   = 0
	EMPTY  = 1
	OBJECT = 2
	ROBOT  = 3

	OBJECT_RIGHT = 4
)

type Grid [][]int

func (g Grid) Write(p2 bool) string {
	var s string
	for _, row := range g {
		for _, cell := range row {
			switch cell {
			case WALL:
				s += "#"
			case EMPTY:
				s += "."
			case OBJECT:
				if p2 {
					s += "["
				} else {
					s += "O"
				}
			case OBJECT_RIGHT:
				s += "]"
			case ROBOT:
				s += "@"
			}
		}
		s += "\n"
	}
	return s
}

func (g Grid) Move(robot *twod.Pos, command string, doMove bool) *twod.Pos {
	return g.moveInternal(robot, command, true)
}

func (g Grid) moveInternal(robot *twod.Pos, command string, doMove bool) *twod.Pos {
	current := g[robot.Row][robot.Col]

	cand := robot.Clone()
	cand.Add(twod.DirArrows[command])
	// easy case
	if g[cand.Row][cand.Col] == EMPTY {
		if doMove {
			g[robot.Row][robot.Col] = EMPTY
			g[cand.Row][cand.Col] = current
		}
		return cand
	}
	// second easy case, cannot move into wall
	if g[cand.Row][cand.Col] == WALL {
		return robot
	}

	if v := g[cand.Row][cand.Col]; (v == OBJECT || v == OBJECT_RIGHT) && (command == "<" || command == ">") {
		newObj := g.moveInternal(cand, command, doMove)
		if newObj.Equals(cand) {
			return robot
		}
		// can now move this object as if the next space was empty
		if g[cand.Row][cand.Col] != EMPTY {
			panic("inconsistency")
		}
		if doMove {
			g[robot.Row][robot.Col] = EMPTY
			g[cand.Row][cand.Col] = current
		}
		return cand
	}

	// moving wide boxes is, fun?
	// send a tracer on both sides of the box
	stack := make([][]*twod.Pos, 0)
	stack = append(stack, make([]*twod.Pos, 0, 2))
	offset := twod.DirArrows[command]
	otherCand := cand.Clone()
	if g[cand.Row][cand.Col] == OBJECT {
		otherCand.Add(twod.DirArrows[">"])
	} else {
		cand.Add(twod.DirArrows["<"])
	}
	stack[0] = append(stack[0], cand, otherCand)

	canMove := true
	for {
		next := slice.Map(stack[len(stack)-1], func(p *twod.Pos) *twod.Pos {
			c := p.Clone()
			c.Add(offset)
			return c
		})
		next = slice.Filter(next, func(p *twod.Pos) bool {
			return g[p.Row][p.Col] != EMPTY
		})

		if len(next) == 0 {
			break
		}

		// if they are all empty, we can move
		anyWalls := false
		for _, p := range next {
			if g[p.Row][p.Col] == WALL {
				anyWalls = true
			}
		}
		if anyWalls {
			canMove = false
			break
		}

		if g[next[0].Row][next[0].Col] == OBJECT_RIGHT {
			need := next[0].Clone()
			need.Add(twod.DirArrows["<"])
			next = append([]*twod.Pos{need}, next...)
		}
		if g[next[len(next)-1].Row][next[len(next)-1].Col] == OBJECT {
			need := next[len(next)-1].Clone()
			need.Add(twod.DirArrows[">"])
			next = append(next, need)
		}

		next = slice.Filter(next, func(p *twod.Pos) bool {
			return g[p.Row][p.Col] != EMPTY
		})

		stack = append(stack, next)
	}

	if canMove {
		for i := len(stack) - 1; i >= 0; i-- {
			for _, p := range stack[i] {
				val := g[p.Row][p.Col]
				g[p.Row][p.Col] = EMPTY
				p.Add(offset)
				g[p.Row][p.Col] = val
			}
		}

		g[robot.Row][robot.Col] = EMPTY
		robot.Add(offset)
		g[robot.Row][robot.Col] = current

		return robot
	}

	// cannot move
	return robot
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	p2 := flag.Bool("p2", false, "run part 2")
	flag.Parse()

	// load grid
	grid := make(Grid, 0)
	var robot *twod.Pos
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		row := make([]int, 0)
		w := 0
		for _, c := range line {
			switch c {
			case '#':
				row = append(row, WALL)
				if *p2 {
					row = append(row, WALL)
				}
			case '.':
				row = append(row, EMPTY)
				if *p2 {
					row = append(row, EMPTY)
				}
			case 'O':
				row = append(row, OBJECT)
				if *p2 {
					row = append(row, OBJECT_RIGHT)
				}
			case '@':
				row = append(row, ROBOT)
				robot = &twod.Pos{Row: len(grid), Col: w}
				if *p2 {
					row = append(row, EMPTY)
				}
			}
			w++
			if *p2 {
				w++
			}
		}
		grid = append(grid, row)
	}
	log.Info("loaded grid", "robot", robot)
	fmt.Printf("LOADED\n%s", grid.Write(*p2))

	instr := ""
	for scanner.Scan() {
		instr += scanner.Text()
	}
	log.Infow("loaded instructions", "instructions", instr)
	instr = strings.ReplaceAll(instr, "v", "V")

	for _, c := range instr {
		log.Infow("moving", "robot", robot, "command", string(c))
		robot = grid.Move(robot, string(c), true)
		fmt.Printf("MOVED %s\n%s\n\n", string(c), grid.Write(*p2))
	}

	part1 := 0
	for r, row := range grid {
		for c, cell := range row {
			if cell == OBJECT {
				part1 += (100*r + c)
			}
		}
	}
	part := "part1"
	if *p2 {
		part = "part2"
	}
	log.Infow(part, "answer", part1)
}
