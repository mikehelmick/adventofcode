package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	left := make(sort.IntSlice, 0, 1000)
	right := make(sort.IntSlice, 0, 1000)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		inParts := strings.Split(line, " ")
		parts := make([]string, 0, 2)
		for _, p := range inParts {
			if p != "" {
				parts = append(parts, p)
			}
		}

		if len(parts) != 2 {
			log.Errorw("invalid input", "line", line, "parts", parts)
			panic("invalid input")
		}

		lInt, err := strconv.Atoi(parts[0])
		if err != nil {
			panic(err)
		}
		left = append(left, lInt)
		rInt, err := strconv.Atoi(parts[1])
		if err != nil {
			panic(err)
		}
		right = append(right, rInt)
	}
	if err := scanner.Err(); err != nil {
		log.Errorw("io error", "error", err)
	}

	left.Sort()
	right.Sort()
	if len(left) != len(right) {
		log.Errorw("invalid input", "left", len(left), "right", len(right))
		panic("invalid input")
	}

	part1 := 0
	for i, l := range left {
		r := right[i]
		if diff := l - r; diff > 0 {
			part1 += diff
		} else {
			part1 += (-1 * diff)
		}
	}
	log.Infof("Part 1: %d", part1)

	rightIndex := indexList(right)
	part2 := 0
	for _, lValue := range left {
		if rCount, ok := rightIndex[lValue]; ok {
			part2 += (lValue * rCount)
		}
	}
	log.Infof("Part 2: %d", part2)

}

func indexList(list []int) map[int]int {
	index := make(map[int]int)
	for _, v := range list {
		if _, ok := index[v]; !ok {
			index[v] = 0
		}
		index[v]++
	}
	return index
}
