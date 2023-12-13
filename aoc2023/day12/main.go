package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/go-functional/slice"
)

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())

	log := logging.FromContext(ctx)
	scanner := bufio.NewScanner(os.Stdin)

	part1 := 0
	part2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		parts := strings.Split(line, " ")
		groupings := slice.Map[string, int](strings.Split(parts[1], ","), func(s string) int {
			i, err := strconv.Atoi(s)
			if err != nil {
				panic(err)
			}
			return i
		})

		part1Matches := findMatches(map[string]int{}, parts[0], groupings)
		part1 += part1Matches
		log.Debugw("line", "matches", part1Matches, "segments", parts[0], "groupings", groupings)

		// Expand input for part 2
		xGroupings := append(append(append(append(groupings, groupings...), groupings...), groupings...), groupings...)
		xSegments := strings.Join([]string{parts[0], parts[0], parts[0], parts[0], parts[0]}, "?")
		matches := findMatches(map[string]int{}, xSegments, groupings)
		log.Debugw("line", "matches", matches, "segments", xSegments, "groupings", xGroupings)
		part2 += matches
	}

	log.Infow("answer", "part1", part1)
	log.Infow("answer", "part2", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}

func findMatches(memo map[string]int, segments string, groupings []int) int {
	key := fmt.Sprintf("%s-%+v", segments, groupings)
	if v, ok := memo[key]; ok {
		return v
	}

	// nothing left to match
	if len(groupings) == 0 {
		// but there are for sure # in the tail, so not a match
		if strings.Contains(segments, "#") {
			return 0
		}
		return 1
	}

	// nothing left in the input, but still groups to match
	if len(segments) == 0 {
		return 0
	}

	// Consume any .
	if strings.HasPrefix(segments, ".") {
		segments = strings.TrimLeft(segments, ".")
		memo[key] = findMatches(memo, segments, groupings)
		return memo[key]
	}

	// Test both options for a ?, this is the only branch.
	if strings.HasPrefix(segments, "?") {
		// first one is a '.', but would get stripped out anyway
		memo[key] = findMatches(memo, segments[1:], groupings) + findMatches(memo, "#"+segments[1:], groupings)
		return memo[key]
	}

	// must start w/ #
	// not enough input left for the next grouping.
	if len(segments) < groupings[0] {
		return 0
	}

	// cannot fit the required length of the next group.
	if strings.Contains(segments[0:groupings[0]], ".") {
		return 0
	}

	// if more than one group, make sure there is a space after this one is consumed.
	if len(groupings) > 1 {
		if len(segments) < groupings[0]+1 || segments[groupings[0]:groupings[0]+1] == "#" {
			return 0 //cannot go right into another sequence
		}
		segments = segments[groupings[0]+1:]
		groupings = groupings[1:]
		memo[key] = findMatches(memo, segments, groupings)
		return memo[key]
	}

	// consume last group.
	segments = segments[groupings[0]:]
	groupings = groupings[1:]
	memo[key] = findMatches(memo, segments, groupings)
	return memo[key]
}
