package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/aoc2023/pkg/logging"
)

var (
	partOne = map[string]string{
		"1": "1",
		"2": "2",
		"3": "3",
		"4": "4",
		"5": "5",
		"6": "6",
		"7": "7",
		"8": "8",
		"9": "9",
	}

	partTwo = map[string]string{
		"one":   "1",
		"two":   "2",
		"three": "3",
		"four":  "4",
		"five":  "5",
		"six":   "6",
		"seven": "7",
		"eight": "8",
		"nine":  "9",
		"1":     "1",
		"2":     "2",
		"3":     "3",
		"4":     "4",
		"5":     "5",
		"6":     "6",
		"7":     "7",
		"8":     "8",
		"9":     "9",
	}
)

func firstDigit(s string, m map[string]string) string {
	pos := len(s) + 1
	val := ""
	for k, v := range m {
		if p := strings.Index(s, k); p >= 0 && p < pos {
			pos = p
			val = v
		}
	}
	return val
}

func lastDigit(s string, m map[string]string) string {
	pos := -1
	val := ""
	for k, v := range m {
		if p := strings.LastIndex(s, k); p >= 0 && p > pos {
			pos = p
			val = v
		}
	}
	return val
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

		p1 := firstDigit(line, partOne) + lastDigit(line, partOne)
		p2 := firstDigit(line, partTwo) + lastDigit(line, partTwo)

		v, err := strconv.Atoi(p1)
		if err == nil {
			part1 += v
		}

		log.Debugf("%v -> %v", line, p2)
		p2v, err := strconv.Atoi(p2)
		if err != nil {
			log.Panicf("unable to parse %q input %q, err: %v", p2, line, err)
		}
		part2 += p2v
	}
	if err := scanner.Err(); err != nil {
		log.Errorw("io error", "error", err)
	}

	log.Infow("part1", "answer", part1)
	log.Infow("part2", "answer", part2)
}
