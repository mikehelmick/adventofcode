package threed

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/mikehelmick/go-functional/slice"
)

var (
	Adj = []*Pos{
		NewPos(-1, 0, 0),
		NewPos(1, 0, 0),
		NewPos(0, -1, 0),
		NewPos(0, 1, 0),
		NewPos(0, 0, -1),
		NewPos(0, 0, 1),
	}
)

type Pos struct {
	X int
	Y int
	Z int
}

func ParsePos(s string) *Pos {
	parts := strings.Split(s, ",")
	ints := slice.Map(parts,
		func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			return i
		})
	return NewPos(ints[0], ints[1], ints[2])
}

func NewPos(x, y, z int) *Pos {
	return &Pos{
		X: x, Y: y, Z: z,
	}
}

func (p *Pos) String() string {
	return fmt.Sprintf("{%v,%v,%v}", p.X, p.Y, p.Z)
}

type ValidFunc func(p *Pos) bool

func (p *Pos) Equals(o *Pos) bool {
	return p.X == o.X && p.Y == o.Y && p.Z == o.Z
}

func (p *Pos) Neighbors(f ValidFunc) []*Pos {
	neighbors := make([]*Pos, 0, len(Adj))
	for _, d := range Adj {
		n := p.Clone()
		n.Add(d)
		if f(n) {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (p *Pos) Clone() *Pos {
	return NewPos(p.X, p.Y, p.Z)
}

func (p *Pos) Add(o *Pos) {
	p.X += o.X
	p.Y += o.Y
	p.Z += o.Z
}
