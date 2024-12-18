package main

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

func solveMul(s string) int64 {
	if !strings.HasPrefix(s, "mul(") || !strings.HasSuffix(s, ")") {
		panic("invalid input")
	}
	s = strings.TrimSuffix(strings.TrimPrefix(s, "mul("), ")")
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		panic("invalid input")
	}
	a, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	b, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return int64(a * b)
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	re, err := regexp.Compile(`(do\(\))|(don't\(\))|(mul\(\d{1,3},\d{1,3}\))`)
	if err != nil {
		panic(err)
	}

	part1 := int64(0)
	part2 := int64(0)
	enabled := true
	for scanner.Scan() {
		line := scanner.Text()

		matches := re.FindAllString(line, -1)
		if matches == nil {
			log.Errorw("no matches found", "line", line)
			continue
		}

		for _, match := range matches {
			if match == `do()` {
				log.Infow("enabling mul", "match", match)
				enabled = true
				continue
			} else if match == `don't()` {
				log.Infow("disabling mul", "match", match)
				enabled = false
				continue
			}
			log.Infow("solving match", "match", match)
			ans := solveMul(match)
			part1 += ans
			if enabled {
				part2 += ans
			}
		}
	}
	log.Infow("Part 1", "part1", part1)
	log.Infow("Part 1", "part1", part2)
}
