package main

import (
	"bufio"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Rule struct {
	Before int
	After  int
}

func (r Rule) Validate(p1, p1idx, p2, p2idx int) bool {
	if p1 == r.Before && p2 == r.After {
		return p1idx < p2idx
	}
	if p1 == r.After && p2 == r.Before {
		return p1idx > p2idx
	}
	// pages don't match this rule.
	return true
}

func (r *Rule) AddToMap(m map[int][]*Rule) {
	if _, ok := m[r.Before]; !ok {
		m[r.Before] = make([]*Rule, 0, 10)
	}
	m[r.Before] = append(m[r.Before], r)

	if _, ok := m[r.After]; !ok {
		m[r.After] = make([]*Rule, 0, 10)
	}
	m[r.After] = append(m[r.After], r)
}

func NewRule(s string) *Rule {
	parts := strings.Split(s, "|")
	before, err := strconv.Atoi(strings.TrimSpace(parts[0]))
	if err != nil {
		panic(err)
	}
	after, err := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err != nil {
		panic(err)
	}
	return &Rule{Before: before, After: after}
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	rules := make([]*Rule, 0, 100)
	ruleMap := make(map[int][]*Rule)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		rule := NewRule(line)
		rules = append(rules, rule)
		rule.AddToMap(ruleMap)
	}

	part1 := 0
	part2 := 0
	for scanner.Scan() {
		line := scanner.Text()
		pages := getPages(line)

		valid := pageOrderValid(pages, rules)
		if valid {
			log.Infow("valid row", "pages", pages)
			part1 += pages[len(pages)/2]
		} else {
			log.Infow("invalid row, sorting", "pages", pages)
			sortPages(pages, ruleMap)
			if pageOrderValid(pages, rules) {
				part2 += pages[len(pages)/2]
			} else {
				log.Fatalw("invalid row after sorting", "pages", pages)
			}
		}
	}
	log.Infow("Part 1", "valid", part1)
	log.Infow("Part 2", "valid", part2)
}

func pageOrderValid(pages []int, rules []*Rule) bool {
	valid := true
	for i, before := range pages {
		for j := i + 1; j < len(pages); j++ {
			after := pages[j]
			for _, rule := range rules {
				if !rule.Validate(before, i, after, j) {
					valid = false
					break
				}
			}
		}
		if !valid {
			break
		}
	}
	return valid
}

func sortPages(pages []int, rules map[int][]*Rule) {
	sort.Slice(pages, func(i, j int) bool {
		rules, ok := rules[pages[i]]
		if !ok {
			// there are no rules about the first page, ordering doesn't matter
			return true
		}

		for _, rule := range rules {
			if rule.Before == pages[i] && rule.After == pages[j] {
				// already in the right order
				return true

			}
			if rule.After == pages[i] && rule.Before == pages[j] {
				return false
			}
		}
		return true
	})
}

func getPages(line string) []int {
	parts := strings.Split(line, ",")
	pages := make([]int, 0, len(parts))
	for _, p := range parts {
		page, err := strconv.Atoi(p)
		if err != nil {
			panic(err)
		}
		pages = append(pages, page)
	}
	return pages
}
