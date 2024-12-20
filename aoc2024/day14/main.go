package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/straid"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

type Robot struct {
	Pos      *twod.Pos
	Velocity *twod.Pos
}

func (r *Robot) Clone() *Robot {
	return &Robot{
		Pos:      r.Pos.Clone(),
		Velocity: r.Velocity.Clone(),
	}
}

func (r *Robot) Move(times int, width int, height int) {
	for i := 0; i < times; i++ {
		r.Pos.Add(r.Velocity)
	}
	r.Pos.Row = r.Pos.Row % height
	if r.Pos.Row < 0 {
		r.Pos.Row += height
	}
	r.Pos.Col = r.Pos.Col % width
	if r.Pos.Col < 0 {
		r.Pos.Col += width
	}
}

func (r *Robot) Quadrant(width int, height int) int {
	// Top left
	if r.Pos.Row < height/2 && r.Pos.Col < width/2 {
		return 1
	}
	// Top right
	if r.Pos.Row < height/2 && r.Pos.Col > width/2 {
		return 2
	}
	// Bottom left
	if r.Pos.Row > height/2 && r.Pos.Col < width/2 {
		return 3
	}
	if r.Pos.Row > height/2 && r.Pos.Col > width/2 {
		return 4
	}
	return 0
}

func NewRobot(line string) *Robot {
	line = strings.ReplaceAll(line, " v=", ",")
	line = strings.ReplaceAll(line, "p=", "")
	parts := strings.Split(line, ",")

	return &Robot{
		Pos:      twod.NewPos(straid.AsInt32(parts[1]), straid.AsInt32(parts[0])),
		Velocity: twod.NewPos(straid.AsInt32(parts[3]), straid.AsInt32(parts[2])),
	}
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	dims := scanner.Text()
	parts := strings.Split(dims, ",")
	height := straid.AsInt32(parts[0])
	width := straid.AsInt32(parts[1])

	robots := make([]*Robot, 0)
	p2robots := make([]*Robot, 0)
	for scanner.Scan() {
		line := scanner.Text()
		log.Debugf("Line: %s", line)
		robot := NewRobot(line)
		log.Debugw("loaded robot", "position", robot.Pos, "velocity", robot.Velocity)
		robots = append(robots, robot)
		p2robots = append(p2robots, robot.Clone())
	}

	{
		quads := make(map[int]int)
		for _, robot := range robots {
			robot.Move(100, width, height)
			quad := robot.Quadrant(width, height)
			quads[quad]++
		}
		log.Debugw("quads", "quads", quads)

		total := 0
		for _, v := range quads {
			total += v
		}
		log.Debugw("robots", "in", len(robots), "out", total)
		part1 := quads[1] * quads[2] * quads[3] * quads[4]
		log.Infow("part1", "part1", part1)
	}

	for i := 1; i <= 10000; i++ {
		for _, robot := range p2robots {
			robot.Move(1, width, height)
		}
		if Draw(i, p2robots, width, height) {
			log.Infow("part2", "part2", i)
			break
		}
	}
}

func Draw(iter int, robots []*Robot, width int, height int) bool {
	out := strings.Builder{}

	rMap := make(map[twod.Pos]int)
	for _, r := range robots {
		rMap[*r.Pos]++
	}

	for r := 0; r < height; r++ {
		line := strings.Builder{}
		for c := 0; c < width; c++ {
			if rMap[twod.Pos{Row: r, Col: c}] > 0 {
				line.WriteString("#")
			} else {
				line.WriteString(".")
			}
		}
		out.WriteString(line.String())
		out.WriteString("\n")
	}

	thisItr := out.String()
	if strings.Contains(thisItr, "########") {
		fmt.Println("Iter", iter)
		fmt.Println(out.String())
		return true
	}
	return false
}
