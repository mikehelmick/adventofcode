package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/straid"
)

func process(in []int) []int {
	out := make([]int, 0, len(in)*2)

	for _, v := range in {
		if v == 0 {
			out = append(out, 1)
			continue
		}
		asStr := fmt.Sprintf("%d", v)
		if len(asStr)%2 == 0 {
			out = append(out, int(straid.AsInt(asStr[0:len(asStr)/2])))
			out = append(out, int(straid.AsInt(asStr[len(asStr)/2:])))
			continue
		}
		out = append(out, v*2024)
	}

	return out
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	line := scanner.Text()

	parts := strings.Split(line, " ")
	stones := make(map[int64]int64)
	for _, part := range parts {
		stones[straid.AsInt(part)] += 1
	}
	log.Debugw("stones", "stones", stones)

	for i := 0; i < 75; i++ {
		next := make(map[int64]int64)
		for stone, count := range stones {
			if stone == 0 {
				next[1] += count
				continue
			}
			asStr := fmt.Sprintf("%d", stone)
			if len(asStr)%2 == 0 {
				next[straid.AsInt(asStr[0:len(asStr)/2])] += count
				next[straid.AsInt(asStr[len(asStr)/2:])] += count
				continue
			}
			next[stone*2024] += count
		}
		stones = next
		log.Debugw("blink", "round", i+1, "stones", len(stones), "total", total(stones))
		if i == 24 {
			log.Infow("answer", "part1", total(stones))
		}
	}
	log.Infow("answer", "total", total(stones))
}

func total(stones map[int64]int64) int64 {
	sum := int64(0)
	for _, v := range stones {
		sum += v
	}
	return sum
}
