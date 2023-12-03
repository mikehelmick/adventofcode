package main

import (
	"bufio"
	"os"
	"sort"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

var (
	partOne = map[string]int{"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}
	partTwo = map[string]int{
		"one": 1, "two": 2, "three": 3, "four": 4, "five": 5, "six": 6, "seven": 7, "eight": 8, "nine": 9,
		"1": 1, "2": 2, "3": 3, "4": 4, "5": 5, "6": 6, "7": 7, "8": 8, "9": 9}
)

type Digit struct {
	Pos   int
	Value int
}

type Index []Digit

func indexDigits(s string, m map[string]int) Index {
	rtn := make(Index, 0, 2)
	for k, v := range m {
		r := s
		firstPos := strings.Index(r, k)
		if firstPos >= 0 {
			rtn = append(rtn, Digit{firstPos, v})
		}
		if lastPos := strings.LastIndex(r, k); lastPos >= 0 && lastPos != firstPos {
			rtn = append(rtn, Digit{lastPos, v})
		}
	}
	sort.Slice(rtn, func(i, j int) bool {
		return rtn[i].Pos <= rtn[j].Pos
	})
	return rtn
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	part1 := 0
	part2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		p1idx := indexDigits(line, partOne)
		if len(p1idx) > 0 { // some examples in the p2 example don't parse
			part1 += (p1idx[0].Value*10 + p1idx[len(p1idx)-1].Value)
		}

		p2idx := indexDigits(line, partTwo)
		log.Debugf("%v -> %+v", line, p2idx)
		part2 += (p2idx[0].Value*10 + p2idx[len(p2idx)-1].Value)
	}
	if err := scanner.Err(); err != nil {
		log.Errorw("io error", "error", err)
	}

	log.Infow("part1", "answer", part1)
	log.Infow("part2", "answer", part2)
}
