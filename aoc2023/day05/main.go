// Day 5 solution.
// I suspect this isn't perfect, but solved my input.

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/go-functional/slice"
)

// Represents a source->dest mapping range from the input
type Range struct {
	Destination int64
	Source      int64
	Length      int64
}

// True if a single int is in the range
func (r *Range) IsMapped(src int64) bool {
	return src >= r.Source && src < r.Source+r.Length
}

// Transform a single int
func (r *Range) Map(src int64) int64 {
	if r.IsMapped(src) {
		return r.Destination + (src - r.Source)
	}
	// not mapped
	return src
}

func (r *Range) String() string {
	return fmt.Sprintf("dest: %d src: %d range: %d", r.Destination, r.Source, r.Length)
}

// Parser for input lines that represent ranges.
func NewRange(s string) *Range {
	parts := strings.Split(strings.TrimSpace(s), " ")

	dest, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	src, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	rng, err := strconv.Atoi(parts[2])
	if err != nil {
		panic(err)
	}

	return &Range{
		Destination: int64(dest),
		Source:      int64(src),
		Length:      int64(rng),
	}
}

// Walker is a struct that represents a range of data moving through the system.
type Walker struct {
	Start  int64
	Length int64
}

// Talk an input walker and transform it to a set of outputs.
// Will be split as it intersects with input ranges.
func (w *Walker) Map(ranges []*Range) []Walker {
	log := logging.DefaultLogger()
	rtn := make([]Walker, 0)

	toSplit := Walker{Start: w.Start, Length: w.Length}

	log.Debugw("splitting range:", "range", toSplit)
	for toSplit.Length > 0 {
		orig := toSplit

		for i, r := range ranges {
			if toSplit.Start >= r.Source+r.Length {
				if i == len(ranges)-1 {
					// above all the ranges
					rtn = append(rtn, toSplit)
					toSplit = Walker{Length: 0}
					break
				}
				// the toSplit is above this range
				continue
			}
			// check to see if it is below the range.
			if toSplit.Start < r.Source {
				// There is an unmapped section.
				potentialLength := (r.Source) - w.Start
				if potentialLength > toSplit.Length {
					// whole range is unmapped
					rtn = append(rtn, toSplit)
					break
				}
				// take part of the range
				rtn = append(rtn, Walker{Start: toSplit.Start, Length: potentialLength})
				toSplit = Walker{Start: r.Source, Length: toSplit.Length - potentialLength}
				break
			}
			// isn't above the range, isn't below the range, so starts in the range.
			if toSplit.Start+toSplit.Length <= r.Source+r.Length {
				// totally in the range
				rtn = append(rtn, Walker{Start: r.Map(toSplit.Start), Length: toSplit.Length})
				toSplit = Walker{Length: 0} // taken the whole range.
				break
			}
			if toSplit.Start == r.Source {
				// exact start
				if toSplit.Length <= r.Length {
					// map entire range
					rtn = append(rtn, toSplit)
					toSplit = Walker{Length: 0}
					break
				}
			}

			// else - starts in a range and goes above the range.
			toTake := (r.Source + r.Length) - toSplit.Start
			// add the mapped part of the range to the output
			newWalker := Walker{Start: r.Map(toSplit.Start), Length: toTake}
			rtn = append(rtn, newWalker)
			// keep the remaining part of the input for next.
			toSplit = Walker{Start: r.Source + r.Length, Length: toSplit.Length - toTake}
		}

		if orig.Start == toSplit.Start && orig.Length == toSplit.Length {
			rtn = append(rtn, toSplit)
			break
		}
	}
	log.Debug("done splitting range", "results", rtn)
	return rtn
}

func (w *Walker) String() string {
	return fmt.Sprintf("s: %d len: %d", w.Start, w.Length)
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	seeds := make([]int64, 0)
	ranges := make(map[string][]*Range)
	conversions := make(map[string]string)
	current := ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "seeds:") {
			seedsS := strings.Split(line[7:], " ")
			seeds = slice.Map(seedsS, func(s string) int64 {
				v, err := strconv.Atoi(s)
				if err != nil {
					panic(err)
				}
				return int64(v)
			})
			log.Debugw("loaded seeds", "seeds", seeds)
			continue
		}

		if strings.HasSuffix(line, " map:") {
			parts := strings.Split(line, " ")
			current = parts[0]
			ranges[current] = make([]*Range, 0)

			parts = strings.Split(current, "-to-")
			conversions[parts[0]] = parts[1]
			continue
		}

		// Range to parse.
		ranges[current] = append(ranges[current], NewRange(line))
	}
	for k, r := range ranges {
		sort.Slice(r,
			func(i, j int) bool {
				return r[i].Source <= r[j].Source
			})
		ranges[k] = r
	}

	for from, to := range conversions {
		log.Debugw("conversion", "from", from, "to", to)
		conv := from + "-to-" + to
		for _, r := range ranges[conv] {
			log.Debugw("range", "range", r)
		}
	}
	Calculate("part1", seeds, conversions, ranges)

	walkers := make([]Walker, 0)
	for i := 0; i < len(seeds); i += 2 {
		walkers = append(walkers, Walker{Start: seeds[i], Length: seeds[i+1]})
	}
	Calculate2(walkers, conversions, ranges)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}

func Calculate2(walkers []Walker, conversions map[string]string, ranges map[string][]*Range) {
	log := logging.DefaultLogger()

	var numbers int64
	for _, w := range walkers {
		numbers += w.Length
	}

	current := "seed"
	for {
		mapTo := conversions[current]
		log.Debugw("Mapping", "from", current, "to", mapTo, "input", walkers)
		conv := current + "-to-" + mapTo

		next := make([]Walker, 0, len(walkers))
		for _, w := range walkers {
			next = append(next, w.Map(ranges[conv])...)
		}
		walkers = next

		var nextCount int64
		for _, w := range walkers {
			nextCount += w.Length
		}
		if nextCount != numbers {
			for _, w := range next {
				log.Debugw("LOST LOST", "walker", w)
			}
			log.Errorw("lost numbers", "starting", numbers, "next", nextCount)
		}

		current = mapTo
		if current == "location" {
			break
		}
	}
	sort.Slice(walkers, func(i, j int) bool {
		return walkers[i].Start <= walkers[j].Start
	})
	log.Infow("part2", "answer", walkers[0].Start)
}

func Calculate(part string, seeds []int64, conversions map[string]string, ranges map[string][]*Range) {
	log := logging.DefaultLogger()

	current := "seed"
	values := make([]int64, len(seeds))
	copy(values, seeds)
	for {
		mapTo := conversions[current]
		log.Debugw("Mapping", "from", current, "to", mapTo, "input", values)
		conv := current + "-to-" + mapTo

		next := slice.Map(values, func(in int64) int64 {
			for _, r := range ranges[conv] {
				if r.IsMapped(in) {
					return r.Map(in)
				}
			}
			return in
		})
		values = next

		current = mapTo
		if current == "location" {
			break
		}
	}
	sort.Slice(values, func(i, j int) bool { return values[i] <= values[j] })
	log.Infow(part, "answer", values[0])
}
