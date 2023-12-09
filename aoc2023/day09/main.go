package main

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Oasis struct {
	Rows [][]int
}

func (o *Oasis) FirstValFromFirstRow() int {
	return o.Rows[0][0]
}

func (o *Oasis) LastValFromFirstRow() int {
	return o.Rows[0][len(o.Rows[0])-1]
}

func (o *Oasis) Fill() {
	last := len(o.Rows) - 1

	cur := o.Rows[last]
	for !AllZeros(cur) {
		next := make([]int, len(cur)-1)
		for i := range next {
			next[i] = cur[i+1] - cur[i]
		}
		o.Rows = append(o.Rows, next)
		cur = next
	}
}

func (o *Oasis) ExpandLeft() {
	last := len(o.Rows) - 1
	if !AllZeros(o.Rows[last]) {
		panic("hasn't been filled")
	}
	// expand the last row
	o.Rows[last] = append(o.Rows[last], 0)

	for row := last - 1; row >= 0; row-- {
		// firstVal - x = y
		// -x = y - firstVal
		// x = -1 (y - firstVal)
		newVal := -1 * (o.Rows[row+1][0] - o.Rows[row][0])
		o.Rows[row] = append([]int{newVal}, o.Rows[row]...)
	}
}

func (o *Oasis) Expand() {
	last := len(o.Rows) - 1
	if !AllZeros(o.Rows[last]) {
		panic("hasn't been filled")
	}
	// expand the last row
	o.Rows[last] = append(o.Rows[last], 0)

	for row := last - 1; row >= 0; row-- {
		// x - lastVal = y
		// x = y + lastVal
		nextRow := row + 1
		o.Rows[row] = append(o.Rows[row],
			o.Rows[nextRow][len(o.Rows[nextRow])-1]+o.Rows[row][len(o.Rows[row])-1])
	}
}

func AllZeros(i []int) bool {
	for _, v := range i {
		if v != 0 {
			return false
		}
	}
	return true
}

func New(s string) *Oasis {
	parts := strings.Split(s, " ")

	rows := make([][]int, 0, 1)
	rows = append(rows, make([]int, len(parts)))

	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		rows[0][i] = v
	}

	return &Oasis{
		Rows: rows,
	}
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	lines := make([]*Oasis, 0)
	// kind of inefficient, but was faster to just create 2 copies of the input.
	p2lines := make([]*Oasis, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		lines = append(lines, New(line))
		p2lines = append(p2lines, New(line))
	}
	log.Debugw("loaded", "oasis", lines)

	part1 := 0
	for i, o := range lines {
		o.Fill()
		log.Debugw("filled", "i", i, "oasis", o)
		o.Expand()
		log.Debugw("expand", "i", i, "oasis", o)
		part1 += o.LastValFromFirstRow()
	}
	log.Infow("answer", "part1", part1)

	part2 := 0
	for i, o := range p2lines {
		o.Fill()
		log.Debugw("filled", "i", i, "oasis", o)
		o.ExpandLeft()
		log.Debugw("expand", "i", i, "oasis", o)
		part2 += o.FirstValFromFirstRow()
	}
	log.Infow("answer", "part2", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
