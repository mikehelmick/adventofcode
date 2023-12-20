package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/mathaid"
)

type Signal int

const (
	NONE Signal = iota
	LOW
	HIGH
)

type Module interface {
	Reset()
	Receive(Signal, string) Signal
	GetDestinations() []string
}

type FlipFlop struct {
	Name         string
	On           bool
	Destinations []string
}

func NewFlipFlop(name string, dest []string) *FlipFlop {
	return &FlipFlop{
		Name:         name,
		On:           false,
		Destinations: dest,
	}
}

func (f *FlipFlop) Reset() {
	f.On = false
}

func (f *FlipFlop) GetDestinations() []string {
	return f.Destinations
}

func (f *FlipFlop) Receive(s Signal, from string) Signal {
	if s == HIGH {
		return NONE
	}
	f.On = !f.On
	if f.On {
		return HIGH
	}
	return LOW
}

type Conjunction struct {
	Name         string
	Inputs       map[string]Signal
	Destinations []string
}

func (c *Conjunction) Reset() {
	for k := range c.Inputs {
		c.Inputs[k] = LOW
	}
}

func (c *Conjunction) AddInput(in string) {
	c.Inputs[in] = LOW
}

func NewConjunction(name string, outputs []string) *Conjunction {
	return &Conjunction{
		Name:         name,
		Inputs:       make(map[string]Signal),
		Destinations: outputs,
	}
}

func (c *Conjunction) GetDestinations() []string {
	return c.Destinations
}

func (c *Conjunction) Receive(s Signal, from string) Signal {
	if _, ok := c.Inputs[from]; !ok {
		panic(fmt.Sprintf("%s received input from %s which is unexpected", c.Name, from))
	}
	c.Inputs[from] = s
	output := LOW
	for _, v := range c.Inputs {
		if v == LOW {
			output = HIGH
			break
		}
	}
	return output
}

type Broadcast struct {
	Name         string
	Destinations []string
}

func (b *Broadcast) Reset() {
}

func NewBroadcast(dest []string) *Broadcast {
	return &Broadcast{
		Name:         "broadcaster",
		Destinations: dest,
	}
}

func (b *Broadcast) GetDestinations() []string {
	return b.Destinations
}

func (b *Broadcast) Receive(s Signal, from string) Signal {
	return s
}

func ParseModule(line string) (string, Module, []string) {
	if strings.HasPrefix(line, "broadcaster") {
		parts := strings.Split(line, "->")
		sendsTo := strings.Split(strings.ReplaceAll(parts[1], " ", ""), ",")
		return "broadcaster", NewBroadcast(sendsTo), sendsTo
	}

	line = strings.ReplaceAll(line, " ", "")
	if strings.HasPrefix(line, "%") {
		line = strings.TrimPrefix(line, "%")
		parts := strings.Split(line, "->")
		sendsTo := strings.Split(parts[1], ",")
		return parts[0], NewFlipFlop(parts[0], sendsTo), sendsTo
	}
	if strings.HasPrefix(line, "&") {
		line = strings.TrimPrefix(line, "&")
		parts := strings.Split(line, "->")
		sendsTo := strings.Split(parts[1], ",")
		return parts[0], NewConjunction(parts[0], sendsTo), sendsTo
	}
	panic("invalid input: " + line)
}

type Send struct {
	From  string
	To    string
	Pulse Signal
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	modules := make(map[string]Module)
	inputs := make(map[string][]string)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		name, mod, outputs := ParseModule(line)
		modules[name] = mod
		// create a reverse index so that we can tell all conjunction modules
		// what their inputs are!
		for _, o := range outputs {
			if inputs[o] == nil {
				inputs[o] = make([]string, 0, 2)
			}
			inputs[o] = append(inputs[o], name)
		}
	}
	for n, m := range modules {
		if cm, ok := m.(*Conjunction); ok {
			for _, in := range inputs[n] {
				cm.AddInput(in)
			}
		}
	}
	log.Infow("loaded", "modules", len(modules))
	for name, mod := range modules {
		log.Debugw("module", "name", name, "mod", mod)
	}

	lowSent := 0
	highSent := 0
	for i := 0; i < 1000; i++ {
		r := pressButton(ctx, i, modules, nil)
		lowSent += r.LowSent
		highSent += r.HighSent
	}

	log.Infow("signals", "low", lowSent, "high", highSent)
	part1 := lowSent * highSent
	log.Infow("answer", "part1", part1)

	for _, m := range modules {
		m.Reset()
	}

	// Using some observations from input set.
	report := Report{
		From: map[string]bool{},
	}
	target := ""
	for n, mod := range modules {
		if slices.Contains(mod.GetDestinations(), "rx") {
			target = n
		}
	}
	log.Infow("part2 target", "module", target)
	// then find all of the conjunctions that feed into target
	cycles := map[string]int{}
	for n, mod := range modules {
		if slices.Contains(mod.GetDestinations(), target) {
			cycles[n] = 0
			report.From[n] = true
		}
	}
	log.Infow("part2 cycles", "modules", report.From)
	report.Target = target

	presses := 0
	for {
		presses++
		allFound := true
		res := pressButton(ctx, presses, modules, &report)
		for d, h := range res.ReportHigh {
			if h {
				if v := cycles[d]; v == 0 {
					cycles[d] = presses
				}
			}
		}
		for _, v := range cycles {
			if v == 0 {
				allFound = false
			}
		}
		if allFound {
			break
		}
	}

	cyclesAt := make([]int64, 0)
	for m, v := range cycles {
		log.Infow("cycles", "mod", m, "count", v)
		cyclesAt = append(cyclesAt, int64(v))
	}
	part2 := mathaid.LowestCommonMultiple(cyclesAt[0], cyclesAt[1], cyclesAt[2:]...)
	log.Infow("answer", "part2", part2)
}

type Report struct {
	Target string
	From   map[string]bool
}

type Result struct {
	LowSent    int
	HighSent   int
	ReportHigh map[string]bool
}

func pressButton(ctx context.Context, press int, modules map[string]Module, report *Report) *Result {
	log := logging.FromContext(ctx)

	r := &Result{
		ReportHigh: make(map[string]bool, 0),
	}
	if report != nil {
		for m := range report.From {
			r.ReportHigh[m] = false
		}
	}

	// push the button
	r.LowSent++
	queue := []Send{
		{
			From:  "button",
			To:    "broadcaster",
			Pulse: LOW,
		},
	}

	for len(queue) > 0 {
		msg := queue[0]
		queue = queue[1:]

		if msg.To == "rx" {
			continue
		}

		mod, ok := modules[msg.To]
		if !ok {
			panic("invalid module: " + msg.To)
		}
		out := mod.Receive(msg.Pulse, msg.From)
		if out == NONE {
			continue
		}
		for _, d := range mod.GetDestinations() {
			if out == LOW {

				r.LowSent++
			} else {
				if report != nil {
					if report.Target == d && report.From[msg.To] {
						log.Debugw("high signal", "to", d, "from", msg.To, "presses", press)
						r.ReportHigh[msg.To] = true
					}
				}
				r.HighSent++
			}
			queue = append(queue,
				Send{
					From:  msg.To,
					To:    d,
					Pulse: out,
				})
		}
	}

	return r
}
