package main

import (
	"bufio"
	"context"
	"fmt"
	"maps"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

// For part 1
type Part struct {
	Values map[string]int
}

func (p *Part) Sum() int {
	s := 0
	for _, v := range p.Values {
		s += v
	}
	return s
}

func NewPart(s string) *Part {
	p := Part{
		Values: make(map[string]int),
	}
	s = strings.TrimPrefix(s, "{")
	s = strings.TrimSuffix(s, "}")

	parts := strings.Split(s, ",")
	for _, sub := range parts {
		nval := strings.Split(sub, "=")
		val, err := strconv.Atoi(nval[1])
		if err != nil {
			panic(err)
		}
		p.Values[nval[0]] = val
	}
	return &p
}

// This was a very clever optimization for part 1, that did
// nothing for part 2. oh well.
type RuleFn func(*Part) (bool, string)

type Workflow struct {
	Name       string
	Rules      []RuleFn
	RuleString []string
}

func (wf *Workflow) Sort(p *Part) string {
	for _, r := range wf.Rules {
		match, dest := r(p)
		if match {
			return dest
		}
	}
	panic("part did not match any rules")
}

func writeGtLtFun(r string, op string) RuleFn {
	exAndDest := strings.Split(r, ":")
	dest := exAndDest[1]
	operands := strings.Split(exAndDest[0], op)
	val, err := strconv.Atoi(operands[1])
	xmas := operands[0]
	if err != nil {
		panic(err)
	}
	if op == ">" {
		return func(p *Part) (bool, string) {
			if p.Values[xmas] > val {
				return true, dest
			}
			return false, ""
		}
	}
	return func(p *Part) (bool, string) {
		if p.Values[xmas] < val {
			return true, dest
		}
		return false, ""
	}
}

func NewWorkflow(s string) *Workflow {
	nameAndRule := strings.Split(s, "{")
	wf := Workflow{
		Name:       nameAndRule[0],
		Rules:      make([]RuleFn, 0),
		RuleString: make([]string, 0), // preserve for part 2
	}

	parts := strings.Split(strings.TrimSuffix(nameAndRule[1], "}"), ",")
	for _, r := range parts {
		wf.RuleString = append(wf.RuleString, r)
		var rf RuleFn
		if strings.Contains(r, ">") {
			rf = writeGtLtFun(r, ">")
		} else if strings.Contains(r, "<") {
			rf = writeGtLtFun(r, "<")
		} else {
			// write static fun
			rf = func(p *Part) (bool, string) {
				return true, r
			}
		}
		wf.Rules = append(wf.Rules, rf)
	}

	return &wf
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	workflows := make(map[string]*Workflow)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		wf := NewWorkflow(line)
		log.Debugw("loaded", "workflow", wf.Name, "rules", len(wf.Rules))
		workflows[wf.Name] = wf
	}
	log.Infow("loaded workflows", "n", len(workflows))

	parts := make([]*Part, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		parts = append(parts, NewPart(line))
	}
	log.Infow("loaded parts", "n", len(parts))

	buckets := make(map[string][]*Part)
	buckets["in"] = parts

	// part 1
	// progressively process every bucket of parts through the assigned workflow
	accepted := make([]*Part, 0)
	for len(buckets) > 0 {
		next := make(map[string][]*Part)
		for wfName, parts := range buckets {
			if wfName == "R" {
				continue
			}
			if wfName == "A" {
				for _, p := range parts {
					p := p
					accepted = append(accepted, p)
				}
				continue
			}

			for _, part := range parts {
				part := part
				dest := workflows[wfName].Sort(part)
				if next[dest] == nil {
					next[dest] = make([]*Part, 0)
				}
				next[dest] = append(next[dest], part)
			}
		}
		buckets = next
	}

	part1 := 0
	for _, a := range accepted {
		part1 += a.Sum()
	}
	log.Infow("answer", "part1", part1)

	part2(ctx, workflows)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}

// Range is the key to part 2.
// We start with XMAS all 1-4000
// and the range will be split depending on rules.
// At the end, we can sum all the parts in the accepted ranges.
type Range struct {
	Min map[string]int
	Max map[string]int
}

func (r *Range) Possibilities() int64 {
	p := int64(1)
	for k, v := range r.Max {
		p *= int64(v - r.Min[k] + 1) // +1 since inclusive on both ends
	}
	return p
}

func (r *Range) String() string {
	b := strings.Builder{}
	for _, k := range []string{"x", "m", "a", "s"} {
		b.WriteString(fmt.Sprintf("%s:%d-%d ", k, r.Min[k], r.Max[k]))
	}
	return b.String()
}

// Lots of defensive copies because of all the pointers.
func (r *Range) DeepCopy() *Range {
	return &Range{
		Min: maps.Clone(r.Min),
		Max: maps.Clone(r.Max),
	}
}

// i.e. "x" < val, x < 2001,  in: X: 1 .. 4000
// out match: x: 1 .. 2000
// out not match: 2001 .. 4000
func (r *Range) SplitRangeLT(key string, val int) (*Range, *Range) {
	matchR := r.DeepCopy()
	if matchR.Min[key] >= val {
		// range doesn't match
		matchR = nil
	} else {
		// min is in range, walk max down (maybe)
		if matchR.Max[key] >= val {
			matchR.Max[key] = val - 1
		}
	}

	notMatch := r.DeepCopy()
	if notMatch.Min[key] < val {
		notMatch.Min[key] = val
	}
	return matchR, notMatch
}

// i.e. "x" > val, x > 2001,  in: X: 1 .. 4000
// out match: x: 2002 .. 4000
// out not match: 1 .. 2001
func (r *Range) SplitRangeGT(key string, val int) (*Range, *Range) {
	matchR := r.DeepCopy()
	if matchR.Max[key] <= val {
		// range doesn't match
		matchR = nil
	} else {
		// max is in range, maybe walk min up
		if matchR.Min[key] <= val {
			matchR.Min[key] = val + 1
		}
	}

	notMatch := r.DeepCopy()
	if notMatch.Max[key] > val {
		notMatch.Max[key] = val
	}
	return matchR, notMatch
}

func DefaultRange() *Range {
	return &Range{
		Min: map[string]int{"x": 1, "m": 1, "a": 1, "s": 1},
		Max: map[string]int{"x": 4000, "m": 4000, "a": 4000, "s": 4000},
	}
}

func part2(ctx context.Context, workflows map[string]*Workflow) {
	log := logging.FromContext(ctx)

	accepted := make([]*Range, 0)

	buckets := map[string][]*Range{"in": {DefaultRange()}}
	for len(buckets) > 0 {
		next := make(map[string][]*Range)
		fmt.Printf("SPLITTING\n")
		for wf, ranges := range buckets {
			fmt.Printf("%v ->\n", wf)
			for _, r := range ranges {
				fmt.Printf("   -> %v\n", r)
			}
		}

		for wfName, ranges := range buckets {
			if wfName == "R" {
				log.Debugw("rejecting", "ranges", ranges)
				continue
			}
			if wfName == "A" {
				log.Debugw("accepting", "ranges", ranges)
				accepted = append(accepted, ranges...)
				continue
			}

			todo := make([]*Range, len(ranges))
			copy(todo, ranges)

			for len(todo) > 0 {
				rng := todo[0]
				todo = todo[1:]

				for _, ruleString := range workflows[wfName].RuleString {
					if strings.Contains(ruleString, ">") {
						exAndDest := strings.Split(ruleString, ":")
						dest := exAndDest[1]
						operands := strings.Split(exAndDest[0], ">")
						val, _ := strconv.Atoi(operands[1])
						xmas := operands[0]

						match, noMatch := rng.SplitRangeGT(xmas, val)
						if match == nil {
							rng = noMatch
						} else {
							if next[dest] == nil {
								next[dest] = []*Range{match}
							} else {
								next[dest] = append(next[dest], match)
							}
						}
						rng = noMatch

					} else if strings.Contains(ruleString, "<") {
						exAndDest := strings.Split(ruleString, ":")
						dest := exAndDest[1]
						operands := strings.Split(exAndDest[0], "<")
						val, _ := strconv.Atoi(operands[1])
						xmas := operands[0]

						match, noMatch := rng.SplitRangeLT(xmas, val)
						if match == nil {
							rng = noMatch
						} else {
							if next[dest] == nil {
								next[dest] = []*Range{match}
							} else {
								next[dest] = append(next[dest], match)
							}
						}
						rng = noMatch

					} else {
						if next[ruleString] == nil {
							next[ruleString] = []*Range{rng}
						} else {
							next[ruleString] = append(next[ruleString], rng)
						}
					}
				}
			}
		}
		buckets = next
	}

	log.Infow("ranges", "n", len(accepted))

	var answer int64
	for _, a := range accepted {
		answer += a.Possibilities()
	}

	log.Infow("answer", "part2", answer)
}
