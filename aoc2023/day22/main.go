package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/threed"
)

type Chamber [][][]int

func (c Chamber) Print() {
	fmt.Printf(" -- THE CHAMBER --\n")
	for z := len(c) - 1; z > 0; z-- {
		fmt.Printf("Z=%03d\n", z)
		for x := 0; x < len(c[z]); x++ {
			for y := 0; y < len(c[z][x]); y++ {
				fmt.Printf("%5d", c[z][x][y])
			}
			fmt.Printf("\n")
		}
	}
	fmt.Printf(" -- END -- \n")
}

func (c Chamber) Get(p *threed.Pos) int {
	return c[p.Z][p.X][p.Y]
}

func (c Chamber) Set(p *threed.Pos, id int) {
	c[p.Z][p.X][p.Y] = id
}

type Brick struct {
	ID      int
	A       *threed.Pos
	B       *threed.Pos
	Points  []*threed.Pos
	RestsOn map[int]bool
}

func NewBrick(id int, l string) *Brick {
	parts := strings.Split(l, "~")
	b := &Brick{
		ID:      id,
		A:       threed.ParsePos(parts[0]),
		B:       threed.ParsePos(parts[1]),
		Points:  make([]*threed.Pos, 0),
		RestsOn: make(map[int]bool),
	}
	b.CalculatePoints()
	return b
}

// which brick does this brick support?
func (b *Brick) Supports(c Chamber) map[int]bool {
	supports := make(map[int]bool)
	for _, p := range b.Points {
		cand := p.Clone()
		cand.Z++
		if above := c.Get(cand); above > 0 && above != b.ID {
			supports[above] = true
		}
	}
	return supports
}

// CalculatePoints materializes all of the points that this brick occupies.
// This takes advantage of an unwritten artifact of the input that the x,y,z
// start and end of a brick are >= from left to right.
func (b *Brick) CalculatePoints() {
	b.Points = make([]*threed.Pos, 0)
	for z := b.A.Z; z <= b.B.Z; z++ {
		for x := b.A.X; x <= b.B.X; x++ {
			for y := b.A.Y; y <= b.B.Y; y++ {
				b.Points = append(b.Points, threed.NewPos(x, y, z))
			}
		}
	}
}

// Place materializes this brick's ID in to the spaces it occupies in the chamber.
// Panics if it cannot be placed.
func (b *Brick) Place(chamber Chamber) {
	for _, p := range b.Points {
		if chamber.Get(p) != 0 {
			panic(fmt.Sprintf("Attempting to overwrite (%v) from %d to %d", p, chamber.Get(p), b.ID))
		}
		chamber.Set(p, b.ID)
	}
}

// Fall attempts to move the brick -1 Z.
func (b *Brick) Fall(chamber Chamber) bool {
	safe := true
	// check every point
	for _, p := range b.Points {
		if p.Z == 1 {
			safe = false
			break
		}
		cand := p.Clone()
		cand.Z--
		if v := chamber.Get(cand); !(v == 0 || v == b.ID) {
			safe = false
			break
		}
	}
	if safe {
		for _, p := range b.Points {
			chamber.Set(p, 0)
		}
		b.A.Z--
		b.B.Z--
		b.CalculatePoints() // calculate new points
		b.Place(chamber)
	}

	return safe
}

// Maxes returns the max x,y,z coordinate based on what we know
// and the context of this brick.
func (b *Brick) Maxes(x, y, z int) (int, int, int) {
	return max(x, b.A.X, b.B.X), max(y, b.A.Y, b.B.Y), max(z, b.A.Z, b.B.Z)
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	bricks := make([]*Brick, 0)
	mX, mY, mZ := 0, 0, 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		id := len(bricks) + 1
		brick := NewBrick(id, line)
		bricks = append(bricks, brick)
		mX, mY, mZ = brick.Maxes(mX, mY, mZ)
	}
	log.Infow("loaded", "maxX", mX, "maxY", mY, "mazZ", mZ, "bricks", len(bricks))

	// Create a 3D array and place all of the bricks into the chamber.
	chamber := make(Chamber, mZ+2)
	for z := 0; z <= mZ+1; z++ {
		grid := make([][]int, 0, mX+1)
		for x := 0; x <= mX; x++ {
			grid = append(grid, make([]int, mY+1))
		}
		chamber[z] = grid
	}
	for _, b := range bricks {
		b.Place(chamber)
	}
	if logging.IsDebug() {
		chamber.Print()
	}

	// sort by initial z to make the falling more efficient.
	sort.Slice(bricks, func(i, j int) bool {
		aMinZ := min(bricks[i].A.Z, bricks[i].B.Z)
		bMinZ := min(bricks[j].A.Z, bricks[j].B.Z)
		return aMinZ < bMinZ
	})

	// until nothing moves
	for {
		anyFell := false
		// attempt to drop each brick -1 Z coordinate.
		for _, b := range bricks {
			anyFell = b.Fall(chamber) || anyFell
		}
		if !anyFell {
			break
		}
	}
	if logging.IsDebug() {
		chamber.Print()
	}

	// Calculate supports and supported by; Index and inverse index.
	// Everyone I support is supported by me.
	supports := make(map[int]map[int]bool)
	supportedBy := make(map[int]map[int]bool)
	for _, b := range bricks {
		iSupport := b.Supports(chamber)
		supports[b.ID] = iSupport
		// invert this index.
		for s := range iSupport {
			if _, ok := supportedBy[s]; !ok {
				supportedBy[s] = make(map[int]bool)
			}
			supportedBy[s][b.ID] = true
		}
	}

	// Part 1. count removable bricks
	part1 := 0
	for _, b := range bricks {
		// easy case, doesn't support anything.
		if len(supports[b.ID]) == 0 {
			log.Debugw("remove brick that doesn't support anything", "id", b.ID)
			part1++
			continue
		}

		// if everything this brick supports is also supported by another brick, then it could be removed
		canRemove := true
		for iSupport := range supports[b.ID] {
			if len(supportedBy[iSupport]) == 1 {
				canRemove = false
				break
			}
		}
		if canRemove {
			log.Debugw("remove brick because everything it supports has another support", "id", b.ID)
			part1++
		}
	}
	log.Infow("answer", "part1", part1)

	part2 := 0
	// For every brick, calculate what would fall if just thar brick was removed
	// not inclusive of the removed brick.
	for _, b := range bricks {
		// doesn't support anything, useless.
		if len(supports[b.ID]) == 0 {
			continue
		}

		// chain is all bricks that would fall if this was removed.
		chain := make(map[int]bool)
		wave := make(map[int]bool)
		for s := range supports[b.ID] {
			// if I am the only support for a brick above, chain reaction
			if len(supportedBy[s]) == 1 {
				chain[s] = true
				wave[s] = true
			}
		}

		for len(wave) > 0 {
			next := make(map[int]bool)
			for w := range wave {
				for s := range supports[w] {
					// s falls if it is only supported by bricks that have also fallen already
					if allIn(chain, supportedBy[s]) {
						chain[s] = true
						next[s] = true
					}
				}
			}
			wave = next
		}

		if len(chain) > 0 {
			log.Debugw("removing", "id", b.ID, "fells", len(chain))
		}
		part2 += len(chain)
	}
	log.Infow("answer", "part2", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}

func allIn(set map[int]bool, subset map[int]bool) bool {
	for s := range subset {
		if !set[s] {
			return false
		}
	}
	return true
}
