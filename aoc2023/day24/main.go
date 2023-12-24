package main

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/threed"
	"github.com/mikehelmick/go-functional/slice"
	"gonum.org/v1/gonum/stat/combin"
)

type Hail struct {
	Position *threed.Pos
	Vector   *threed.Pos
}

func (h *Hail) IsFuture(x, y float64) bool {
	if h.Vector.X > 0 && x < float64(h.Position.X) {
		return false
	}
	if h.Vector.X < 0 && x > float64(h.Position.X) {
		return false
	}
	if h.Vector.Y > 0 && y < float64(h.Position.Y) {
		return false
	}
	if h.Vector.Y < 0 && y > float64(h.Position.Y) {
		return false
	}
	return true
}

func NewHail(s string) *Hail {
	s = strings.ReplaceAll(s, " ", "")
	s = strings.ReplaceAll(s, "@", ",")
	parts := strings.Split(s, ",")
	ints := slice.Map(parts, func(s string) int {
		v, err := strconv.Atoi(s)
		if err != nil {
			panic("cannot parse int")
		}
		return v
	})
	return &Hail{
		Position: threed.NewPos(ints[0], ints[1], ints[2]),
		Vector:   threed.NewPos(ints[3], ints[4], ints[5]),
	}
}

// part2, just 2d linear intersection
func linearIntersection(a, b *Hail) (rX *float64, rY *float64) {
	// lazy handling of divide by zero :)
	defer func() {
		if r := recover(); r != nil {
			rX, rY = nil, nil
		}
	}()

	aEnd := a.Position.Clone()
	aEnd.Add(a.Vector)
	bEnd := b.Position.Clone()
	bEnd.Add(b.Vector)

	x1, y1, x2, y2 := float64(a.Position.X), float64(a.Position.Y), float64(aEnd.X), float64(aEnd.Y)
	x3, y3, x4, y4 := float64(b.Position.X), float64(b.Position.Y), float64(bEnd.X), float64(bEnd.Y)

	t := (((x1 - x3) * (y3 - y4)) - ((y1 - y3) * (x3 - x4))) /
		((x1-x2)*(y3-y4) - (y1-y2)*(x3-x4))

	iX := float64(x1) + t*(x2-x1)
	iY := float64(y1) + t*(y2-y1)

	rX = &iX
	rY = &iY

	return
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	stones := make([]*Hail, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		hail := NewHail(line)
		stones = append(stones, hail)
	}
	log.Debugw("Loaded hail", "hail", stones)

	var minCord float64 = 200000000000000
	var maxCord float64 = 400000000000000

	part1 := 0
	pairs := combin.Combinations(len(stones), 2)
	for _, p := range pairs {
		iX, iY := linearIntersection(stones[p[0]], stones[p[1]])
		if iX == nil {
			log.Debugw("NOPE", "p1", stones[p[0]].Position, "p2", stones[p[1]].Position, "X", iX, "Y", iY)
		} else {
			log.Debugw("intersction", "p1", stones[p[0]].Position, "p2", stones[p[1]].Position, "X", iX, "Y", iY, "p1future", stones[p[0]].IsFuture(*iX, *iY), "p2future", stones[p[1]].IsFuture(*iX, *iY))
			if *iX >= minCord && *iX < maxCord && *iY >= minCord && *iY < maxCord {
				log.Debugw("in range")
				if stones[p[0]].IsFuture(*iX, *iY) && stones[p[1]].IsFuture(*iX, *iY) {
					log.Debugw("colission in future")
					part1++
				}
			}

		}
	}
	log.Infow("answer", "part1", part1)
}
