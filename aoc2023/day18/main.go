package main

import (
	"bufio"
	"context"
	"fmt"
	"image"
	"math/big"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Plan struct {
	Dir   string
	Amt   int
	Color string
}

// For part2, decode color
func (p Plan) Distance() int {
	d, err := strconv.ParseInt(p.Color[0:5], 16, 64)
	if err != nil {
		panic(err)
	}
	return int(d)
}

// For part2, decode color
func (p Plan) Direction() string {
	// 0 means R, 1 means D, 2 means L, and 3 means U.
	switch p.Color[5:] {
	case "0":
		return "R"
	case "1":
		return "D"
	case "2":
		return "L"
	case "3":
		return "U"
	}
	panic("invalid color")
}

var dir = map[string]image.Point{
	"U": {-1, 0},
	"R": {0, 1},
	"D": {1, 0},
	"L": {0, -1},
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	plan := make([]Plan, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, " ")
		amt, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		plan = append(plan, Plan{parts[0], amt, strings.TrimSuffix(strings.TrimPrefix(parts[2], "(#"), ")")})
	}
	//fmt.Printf("%+v\n", plan)

	digger := image.Point{0, 0}
	points := make([]image.Point, len(plan))
	intPoints := int64(1)
	for _, p := range plan {
		digger = digger.Add(dir[p.Dir].Mul(p.Amt))
		points = append(points, digger)
		intPoints += int64(p.Amt)
	}
	inside := shoelace(points)
	//part1 := inside.Add(inside, perimiter)
	//fmt.Printf("part1: %+v\n", inside)
	fmt.Printf("part1: %+v\n", prick(inside, intPoints))

	// part 2 is... bigger
	digger = image.Point{0, 0}
	points = make([]image.Point, len(plan))
	intPoints = int64(1)
	for _, p := range plan {
		//fmt.Printf("digging: %v %v\n", p.Distance(), p.Direction())
		digger = digger.Add(dir[p.Direction()].Mul(p.Distance()))
		points = append(points, digger)
		intPoints += int64(p.Distance())
	}
	part2 := shoelace(points)
	//fmt.Printf("part2: %+v\n", part2)
	fmt.Printf("part2: %+v\n", prick(part2, intPoints))

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}

// prick's theorem
func prick(inside *big.Int, points int64) *big.Int {
	pb := big.NewInt(1 + (points-1)/2)
	return inside.Add(inside, pb)
}

// shoelace algorithm
// found by searching for algorithm for area of convex hull.
func shoelace(poly []image.Point) *big.Int {
	sum := big.NewInt(0)
	p0 := poly[len(poly)-1]
	for _, p1 := range poly {
		p0y := big.NewInt(int64(p0.Y))
		p1x := big.NewInt(int64(p1.X))
		p0x := big.NewInt(int64(p0.X))
		p1y := big.NewInt(int64(p1.Y))

		p0y = p0y.Mul(p0y, p1x)
		p0x = p0x.Mul(p0x, p1y)

		sum = sum.Add(sum, p0y.Sub(p0y, p0x))
		p0 = p1
	}
	return sum.Div(sum, big.NewInt(2))
}

/*
Part 1 was implemented w/ a map + flood fill, but that was too big for part2.

func flood(path map[image.Point]bool, start image.Point) {
	wave := []image.Point{start}
	for len(wave) > 0 {
		next := make([]image.Point, 0)
		for _, p := range wave {
			for _, d := range dir {
				nextP := p.Add(d)
				if !path[nextP] {
					path[nextP] = true
					next = append(next, nextP)
				}
			}
		}
		wave = next
	}
}
*/
