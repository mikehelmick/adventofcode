package twod

import (
	"fmt"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/straid"
)

const (
	RIGHT = "R"
	UP    = "U"
	LEFT  = "L"
	DOWN  = "D"
)

var (
	Dirs = map[string]*Pos{
		"R": NewPos(0, 1),
		"U": NewPos(-1, 0),
		"L": NewPos(0, -1),
		"D": NewPos(1, 0),
	}

	TurnRight = map[string]string{
		"U": "R",
		"R": "D",
		"D": "L",
		"L": "U",
	}

	DirArrows = map[string]*Pos{
		">": NewPos(0, 1),
		"^": NewPos(-1, 0),
		"<": NewPos(0, -1),
		"V": NewPos(1, 0),
	}

	Manhattan = []*Pos{
		{Row: 0, Col: 1},
		{Row: 1, Col: 0},
		{Row: 0, Col: -1},
		{Row: -1, Col: 0},
	}

	Diags = []*Pos{
		{Row: 1, Col: 1},
		{Row: -1, Col: 1},
		{Row: -1, Col: -1},
		{Row: 1, Col: -1},
	}

	Adjacent = []*Pos{
		{Row: -1, Col: -1}, {Row: -1, Col: 0}, {Row: -1, Col: 1},
		{Row: 0, Col: -1}, {Row: 0, Col: 1},
		{Row: 1, Col: -1}, {Row: 1, Col: 0}, {Row: 1, Col: 1},
	}
)

type Pos struct {
	Row int
	Col int
}

func NewPos(r, c int) *Pos {
	return &Pos{
		Row: r,
		Col: c,
	}
}

func FromString(s string) *Pos {
	s = strings.ReplaceAll(s, "{", "")
	s = strings.ReplaceAll(s, "}", "")
	parts := strings.Split(s, ",")

	r := int(straid.AsInt(parts[0]))
	c := int(straid.AsInt(parts[1]))
	return NewPos(r, c)
}

func (p *Pos) String() string {
	return fmt.Sprintf("{%v,%v}", p.Row, p.Col)
}

type ValidFunc func(p *Pos) bool

func (p *Pos) Dist(o *Pos) int {
	x := p.Col - o.Col
	if x <= 0 {
		x *= -1
	}
	y := p.Row - o.Row
	if y <= 0 {
		y *= -1
	}
	return x + y
}

func (p *Pos) Equals(o *Pos) bool {
	return p.Row == o.Row && p.Col == o.Col
}

func (p *Pos) Follow(f ValidFunc, adj []*Pos) []*Pos {
	neighbors := make([]*Pos, 0, len(adj))
	for _, d := range adj {
		n := p.Clone()
		n.Add(d)
		if f(n) {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (p *Pos) ManhattanNeighbors(f ValidFunc) []*Pos {
	neighbors := make([]*Pos, 0, len(Manhattan))
	for _, d := range Manhattan {
		n := p.Clone()
		n.Add(d)
		if f(n) {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (p *Pos) Neighbors(f ValidFunc) []*Pos {
	neighbors := make([]*Pos, 0, len(Dirs))
	for _, d := range Dirs {
		n := p.Clone()
		n.Add(d)
		if f(n) {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (p *Pos) Adjacent(f ValidFunc) []*Pos {
	neighbors := make([]*Pos, 0, len(Adjacent))
	for _, d := range Adjacent {
		n := p.Clone()
		n.Add(d)
		if f(n) {
			neighbors = append(neighbors, n)
		}
	}
	return neighbors
}

func (p *Pos) Clone() *Pos {
	return &Pos{
		Row: p.Row,
		Col: p.Col,
	}
}

func (p *Pos) Add(o *Pos) {
	p.Row += o.Row
	p.Col += o.Col
}
