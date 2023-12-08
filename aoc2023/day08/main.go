package main

import (
	"bufio"
	"context"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/mathaid"
)

type Node struct {
	Name  string
	Left  string
	Right string
}

func NewNode(s string) *Node {
	parts := strings.Split(s, "=")
	name := strings.TrimSpace(parts[0])

	rest := strings.ReplaceAll(parts[1], " ", "")
	rest = strings.ReplaceAll(rest, "(", "")
	rest = strings.ReplaceAll(rest, ")", "")

	parts = strings.Split(rest, ",")

	return &Node{
		Name:  name,
		Left:  parts[0],
		Right: parts[1],
	}
}

type FinishFn func(string) bool

func Traverse(start string, isFinished FinishFn, directions string, nodes map[string]*Node) int {
	steps := 0
	dir := 0
	for !isFinished(start) {
		if dir >= len(directions) {
			dir = 0
		}
		choose := directions[dir : dir+1]
		if choose == "L" {
			start = nodes[start].Left
		} else {
			start = nodes[start].Right
		}
		steps++
		dir++
	}
	return steps
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	directions := ""

	nodes := make(map[string]*Node)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if !strings.Contains(line, "=") {
			directions = line
		} else {
			node := NewNode(line)
			log.Debugw("newNode", "node", node)
			nodes[node.Name] = node
		}
	}
	log.Debugw("loaded", "directions", directions)

	// part1
	steps := Traverse("AAA", func(s string) bool { return s == "ZZZ" }, directions, nodes)
	log.Infow("answer", "part1", steps)

	// part 2
	// start to finish for all starts
	starts := make([]string, 0)
	for _, n := range nodes {
		if strings.HasSuffix(n.Name, "A") {
			starts = append(starts, n.Name)
		}
	}
	log.Debugw("Part 2", "starts", starts)

	stepsToZ := make([]int64, len(starts))
	endsInZ := func(s string) bool { return strings.HasSuffix(s, "Z") }
	for wI, start := range starts {
		stepsToZ[wI] = int64(Traverse(start, endsInZ, directions, nodes))
	}
	log.Debugw("stepsToZ", "steps", stepsToZ)

	answer := mathaid.LowestCommonMultiple(stepsToZ[0], stepsToZ[1], stepsToZ[2:]...)
	log.Infow("answer", "part2", answer)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
