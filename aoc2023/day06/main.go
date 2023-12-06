package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Race struct {
	Time     int64
	Distance int64
}

func (r *Race) String() string {
	return fmt.Sprintf("time: %d distance: %d", r.Time, r.Distance)
}

func (r *Race) WaysToBeat() int64 {
	var wins int64
	for i := int64(1); i < r.Time; i++ {
		covered := (r.Time - i) * i
		if covered > r.Distance {
			wins++
		}
	}
	return wins
}

func getInts(line string) []int64 {
	parts := strings.Split(line, ":")
	numPart := strings.TrimSpace(parts[1])

	nums := strings.Split(numPart, " ")
	rtn := make([]int64, 0, len(nums))

	for _, n := range nums {
		n = strings.TrimSpace(n)
		if n == "" {
			continue
		}
		v, err := strconv.ParseInt(n, 10, 64)
		if err != nil {
			panic(fmt.Sprintf("cannot parse %q: %v", n, err))
		}
		rtn = append(rtn, v)
	}
	return rtn
}

func mergedValue(line string, kind string) int64 {
	numS := strings.ReplaceAll(line[len(kind):], " ", "")
	num, err := strconv.ParseInt(numS, 10, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	races := make([]*Race, 0)
	part2Race := &Race{}
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "Time:") {
			for _, n := range getInts(line) {
				races = append(races, &Race{Time: n})
			}
			part2Race.Time = mergedValue(line, "Time:")
		}
		if strings.HasPrefix(line, "Distance:") {
			for i, n := range getInts(line) {
				races[i].Distance = n
			}
			part2Race.Distance = mergedValue(line, "Distance:")
		}
	}
	log.Debugw("loaded races", "races", races)

	part1 := int64(1)
	for i, r := range races {
		wins := r.WaysToBeat()
		log.Debugw("ways to beat", "id", i, "race", r, "wins", wins)
		part1 *= wins
	}
	log.Infow("answer", "part1", part1)

	log.Debugw("merged race", "race", part2Race)
	part2 := part2Race.WaysToBeat()
	log.Infow("answer", "part2", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
