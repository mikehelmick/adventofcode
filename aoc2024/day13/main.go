package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/twod"
)

type Machine struct {
	ButtonA twod.Pos
	ButtonB twod.Pos
	Prize   twod.Pos
}

func (m Machine) Solve(offset int64) int64 {
	x1 := int64(m.ButtonA.Row)
	y1 := int64(m.ButtonA.Col)
	x2 := int64(m.ButtonB.Row)
	y2 := int64(m.ButtonB.Col)

	c := int64(m.Prize.Row) + offset
	d := int64(m.Prize.Col) + offset

	a := (c*y2 - d*x2) / (x1*y2 - x2*y1)
	b := (d*x1 - c*y1) / (x1*y2 - x2*y1)

	log := logging.DefaultLogger()
	log.Debugw("solving machine", "a", a, "b", b)

	if a*x1+b*x2 == c && a*y1+b*y2 == d {
		return 3*a + b
	}
	return 0
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	machines := make([]Machine, 0, 100)
	for scanner.Scan() {
		buttonA := scanner.Text()
		if !strings.HasPrefix(buttonA, "Button") {
			continue
		}
		scanner.Scan()
		buttonB := scanner.Text()
		scanner.Scan()
		prize := scanner.Text()

		machines = append(machines, Machine{
			ButtonA: parse(buttonA, "+"),
			ButtonB: parse(buttonB, "+"),
			Prize:   parse(prize, "="),
		})
	}
	log.Infow("loaded machines", "machines", machines)

	var part1 int64
	var part2 int64
	for i, m := range machines {
		ans := m.Solve(0)
		ans2 := m.Solve(10000000000000)
		log.Infow("machine solved", "i", i+1, "ans", ans, "ans2", ans2)
		part1 += ans
		part2 += ans2
	}
	log.Infow("part1", "part1", part1)
	log.Infow("part2", "part2", part2)
}

func parse(s string, sep string) twod.Pos {
	first := strings.Index(s, "X"+sep)
	second := strings.Index(s, ", Y"+sep)

	xCord := s[first+2 : second]
	yCord := s[second+4:]

	x, err := strconv.Atoi(xCord)
	if err != nil {
		panic(err)
	}
	y, err := strconv.Atoi(yCord)
	if err != nil {
		panic(err)
	}
	return twod.Pos{Row: x, Col: y}
}
