package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

const EMPTY = 0
const WALL = 1

type Maze [][]int

func (m Maze) String(visited map[twod.Pos]map[string]bool) string {
	b := strings.Builder{}
	for r, row := range m {
		for c, cell := range row {
			if cell == EMPTY {
				if _, ok := visited[twod.Pos{Row: r, Col: c}]; ok {
					b.WriteRune('X')
				} else {
					b.WriteRune('.')
				}
			} else {
				b.WriteRune('#')
			}
		}
		b.WriteRune('\n')
	}
	return b.String()
}

func (m Maze) Clone() Maze {
	newMaze := make(Maze, len(m))
	for i, row := range m {
		newMaze[i] = make([]int, len(row))
		copy(newMaze[i], row)
	}
	return newMaze
}

type Guard struct {
	Position    *twod.Pos
	Dir         string
	Orientation *twod.Pos
}

func (g *Guard) Clone() *Guard {
	return &Guard{
		Position:    g.Position.Clone(),
		Dir:         g.Dir,
		Orientation: g.Orientation.Clone(),
	}
}

func (g *Guard) TurnRight() {
	if g.Dir == "^" {
		g.Dir = ">"
		g.Orientation = twod.DirArrows[">"]
	} else if g.Dir == ">" {
		g.Dir = "V"
		g.Orientation = twod.DirArrows["V"]
	} else if g.Dir == "V" {
		g.Dir = "<"
		g.Orientation = twod.DirArrows["<"]
	} else if g.Dir == "<" {
		g.Dir = "^"
		g.Orientation = twod.DirArrows["^"]
	}
}

func (g *Guard) String() string {
	return fmt.Sprintf("Guard{Position: %v, Dir: %s, Orientation: %v}", g.Position, g.Dir, g.Orientation)
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	maze := make(Maze, 0)
	var guard *Guard
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		row := make([]int, len(line))
		for i, c := range line {
			if c == '#' {
				row[i] = WALL
				continue
			}
			if c != '.' {
				// found the guard
				guard = &Guard{
					Position:    &twod.Pos{Row: len(maze), Col: i},
					Dir:         string(c),
					Orientation: twod.DirArrows[string(c)],
				}
			}
		}
		maze = append(maze, row)
	}
	starting := maze.Clone()
	startingGuard := guard.Clone()
	fmt.Printf("maze :\n%s\n", maze.String(nil))
	log.Infow("Guard", "guard", guard)

	visited, _ := traverse(maze, guard)
	fmt.Printf("maze :\n%s\n", maze.String(visited))
	log.Infow("part 1", "visited", len(visited))

	// Part 2
	part2 := 0
	for newBlock := range visited {
		maze := starting.Clone()
		maze[newBlock.Row][newBlock.Col] = WALL
		guard := startingGuard.Clone()

		if _, exited := traverse(maze, guard); !exited {
			part2++
		}
	}
	log.Infow("part 2", "part2", part2)
}

func traverse(maze Maze, guard *Guard) (map[twod.Pos]map[string]bool, bool) {
	exited := false
	visited := make(map[twod.Pos]map[string]bool)
	for {
		if _, ok := visited[*guard.Position]; !ok {
			visited[*guard.Position] = make(map[string]bool)
		} else {
			if visited[*guard.Position][guard.Dir] {
				break
			}
		}
		visited[*guard.Position][guard.Dir] = true

		next := guard.Position.Clone()
		next.Add(guard.Orientation)
		if next.Row < 0 || next.Row >= len(maze) || next.Col < 0 || next.Col >= len(maze[0]) {
			exited = true
			break
		}

		if maze[next.Row][next.Col] == WALL {
			guard.TurnRight()
			continue
		} else {
			// Move forward
			guard.Position = next
		}
	}
	return visited, exited
}
