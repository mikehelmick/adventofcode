package main

/** Note:
This won't just generally work on any input.

I converted the input into graphviz format and used the neato
layout engine. This made it trivial to identify the 3 edges to remove.

I manually removed them from the input file and then this works.

If you uncomment the printf statements in main, you can generate the graph
and repeat the procedure.
*/

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yourbasic/graph"
)

type Node struct {
	ID    string
	Edges []string
}

func NewNode(n string) *Node {
	return &Node{
		ID:    n,
		Edges: make([]string, 0),
	}
}

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	index := make(map[string]int)
	next := 0
	nodes := make(map[string]*Node)
	//fmt.Printf("graph {\n")
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, ":")
		to := strings.Split(strings.TrimSpace(parts[1]), " ")

		if _, ok := nodes[parts[0]]; !ok {
			nodes[parts[0]] = NewNode(parts[0])
		}
		for _, t := range to {
			if _, ok := nodes[t]; !ok {
				nodes[t] = NewNode(t)
			}
			//fmt.Printf("  %s -- %s\n", parts[0], t)
			nodes[parts[0]].Edges = append(nodes[parts[0]].Edges, t)
			nodes[t].Edges = append(nodes[t].Edges, parts[0])
		}

		if _, ok := index[parts[0]]; !ok {
			index[parts[0]] = next
			next++
		}
		for _, t := range to {
			if _, ok := index[t]; !ok {
				index[t] = next
				next++
			}
		}
	}
	//fmt.Printf("}\n")
	fmt.Printf("total notes: %+v\n", len(nodes))

	g := graph.New(len(nodes))
	for _, n := range nodes {
		for _, e := range n.Edges {
			g.AddBoth(index[n.ID], index[e])
		}
	}

	com := graph.Components(g)
	if len(com) != 2 {
		panic("done messed up")
	}

	fmt.Printf("answer: %v\n", len(com[0])*len(com[1]))
}
