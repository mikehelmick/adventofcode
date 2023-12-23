package main

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Maze [][]string

func (m Maze) isValid(p image.Point) bool {
	return p.X >= 0 && p.Y >= 0 && p.X < len(m) && p.Y < len(m[0])
}

func (m Maze) Start() image.Point {
	for c := 0; c < len(m[0]); c++ {
		if m[0][c] == "." {
			return image.Point{X: 0, Y: c}
		}
	}
	panic("no start")
}

func (m Maze) Print(visited [][]bool, length int) {
	fmt.Printf("maze: %d v: %d\n", length, len(visited))
	for r, row := range m {
		for c, v := range row {
			if visited[r][c] {
				fmt.Printf("o")
			} else {
				fmt.Printf("%s", v)
			}
		}
		fmt.Printf("\n")
	}
}

func (m Maze) End() image.Point {
	for c := 0; c < len(m[0]); c++ {
		if m[len(m)-1][c] == "." {
			return image.Point{X: len(m) - 1, Y: c}
		}
	}
	panic("no start")
}

func LogestPath(ctx context.Context, m Maze, slopes bool) int {
	start := m.Start()
	end := m.End()

	answers := map[int]bool{}
	visited := make([][]bool, 0, len(m))
	for _, row := range m {
		visited = append(visited, make([]bool, len(row)))
	}

	term := Terminator{
		MaxSeen:       0,
		PathsSinceMax: 0,
		Terminate:     false,
	}

	dfs(ctx, m, start, end, visited, &term, answers, 0, slopes)

	logging.FromContext(ctx).Debugw("dfs complete", "answers", answers, slopes)
	ans := 0
	for k := range answers {
		ans = max(ans, k)
	}
	return ans
}

var dirs = []image.Point{image.Pt(-1, 0), image.Pt(0, 1), image.Pt(1, 0), image.Pt(0, -1)}

type Terminator struct {
	MaxSeen       int
	PathsSinceMax int
	Terminate     bool
}

func dfs(ctx context.Context, m Maze, from image.Point, to image.Point, visited [][]bool, term *Terminator, answers map[int]bool, cur int, slopes bool) {
	//log := logging.FromContext(ctx)
	if from.Eq(to) {
		if cur > term.MaxSeen {
			term.MaxSeen = cur
			term.PathsSinceMax = 0
		} else {
			term.PathsSinceMax++
		}
		// This is just a heuristic that if we haven't found longer paths in 15k paths to the exit, bail.
		if term.PathsSinceMax > 15000 {
			term.Terminate = true
		}
		logging.FromContext(ctx).Debugw("possibility", "length", cur, "max", term.MaxSeen, "times", term.PathsSinceMax)
		answers[cur] = true
		return
	}

	for _, dir := range dirs {
		if term.Terminate {
			return
		}

		cand := from.Add(dir)
		if !m.isValid(cand) {
			continue
		}
		if visited[cand.X][cand.Y] {
			continue
		}
		space := m[cand.X][cand.Y]
		if space == "#" {
			continue
		}

		if !slopes {
			visited[cand.X][cand.Y] = true
			dfs(ctx, m, cand, to, visited, term, answers, cur+1, slopes)
			visited[cand.X][cand.Y] = false
			continue
		}

		switch space {
		case "^":
			panic("there are none of these  in the input :)")
		case ">":
			if dir.X == 0 && dir.Y == -1 {
				continue // invalid step
			}

			visited[cand.X][cand.Y] = true
			next := cand.Add(image.Pt(0, 1))
			visited[next.X][next.Y] = true
			dfs(ctx, m, next, to, visited, term, answers, cur+2, slopes)
			visited[next.X][next.Y] = false
			visited[cand.X][cand.Y] = false
			continue
		case "v":
			if dir.X == -1 && dir.Y == 0 {
				continue // invalid step
			}

			visited[cand.X][cand.Y] = true
			next := cand.Add(image.Pt(1, 0))
			visited[next.X][next.Y] = true
			dfs(ctx, m, next, to, visited, term, answers, cur+2, slopes)
			visited[next.X][next.Y] = false
			visited[cand.X][cand.Y] = false
			continue
		case "<":
			if dir.X == 0 && dir.Y == 1 {
				continue // invalid step
			}

			visited[cand.X][cand.Y] = true
			next := cand.Add(image.Pt(0, -1))
			visited[next.X][next.Y] = true
			dfs(ctx, m, next, to, visited, term, answers, cur+2, slopes)
			visited[next.X][next.Y] = false
			visited[cand.X][cand.Y] = false
			continue
		case ".":
			visited[cand.X][cand.Y] = true
			dfs(ctx, m, cand, to, visited, term, answers, cur+1, slopes)
			visited[cand.X][cand.Y] = false
			continue
		case "#":
			continue // dead end.
		}
	}
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	m := make(Maze, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		row := strings.Split(line, "")
		m = append(m, row)
	}
	log.Infow("answer", "part1", LogestPath(ctx, m, true))

	log.Infow("answer", "part2", LogestPath(ctx, m, false))
}
