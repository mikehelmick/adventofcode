package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Report struct {
	Values []int
}

func NewReport(line string) *Report {
	report := &Report{
		Values: make([]int, 0, 5),
	}

	parts := strings.Split(line, " ")
	for _, part := range parts {
		val, err := strconv.Atoi(part)
		if err != nil {
			panic(err)
		}
		report.Values = append(report.Values, val)
	}
	return report
}

func NewDampenedReport(r Report, pos int) *Report {
	dampened := &Report{
		Values: make([]int, 0, len(r.Values)-1),
	}
	for i, v := range r.Values {
		if i == pos {
			continue
		}
		dampened.Values = append(dampened.Values, v)
	}
	return dampened
}

func (r Report) safe() bool {
	diffs := make([]int, 0, 4)
	allPositive := true
	allNegative := true
	for i := 0; i < len(r.Values)-1; i++ {
		diff := r.Values[i+1] - r.Values[i]
		allPositive = allPositive && diff > 0
		allNegative = allNegative && diff < 0
		diffs = append(diffs, diff)
	}
	if !(allPositive || allNegative) {
		return false
	}
	// check magnitude of diffs
	for _, d := range diffs {
		if d < 0 {
			d = -d
		}
		if d < 1 || d > 3 {
			return false
		}
	}
	return true
}

func (r Report) damperSafe() bool {
	if r.safe() {
		return true
	}
	for i := 0; i < len(r.Values); i++ {
		dampened := NewDampenedReport(r, i)
		if dampened.safe() {
			return true
		}
	}
	return false
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	reports := make([]*Report, 0)
	for scanner.Scan() {
		line := scanner.Text()
		reports = append(reports, NewReport(line))
	}

	part1 := 0
	part2 := 0
	for _, report := range reports {
		if report.safe() {
			part1++
		}
		if report.damperSafe() {
			part2++
		}
	}
	log.Infow("answer", "part1", part1)
	log.Infow("answer", "part2", part2)
}
