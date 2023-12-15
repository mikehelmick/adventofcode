package main

import (
	"bufio"
	"context"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/go-functional/slice"
)

type Hasher string

func (h Hasher) Hash() int {
	v := 0
	for _, r := range h {
		ascii := int(r)
		v += ascii
		v *= 17
		v %= 256
	}
	return v
}

func ToLens(h string) Lens {
	parts := strings.Split(string(h), "=")
	if len(parts) == 1 {
		label := strings.TrimSuffix(string(h), "-")
		return Lens{
			Label: Hasher(label),
			Value: -1, // none of the input is negative.
		}
	}
	value, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return Lens{
		Label: Hasher(parts[0]),
		Value: value,
	}
}

type Lens struct {
	Label Hasher
	Value int
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)
	var line string
	for scanner.Scan() {
		line = scanner.Text()
		if line == "" {
			continue
		}
	}
	log.Infow("line", "line", line)

	parts := strings.Split(line, ",")
	tot := 0
	boxes := make([][]Lens, 256)
	for i := range boxes {
		boxes[i] = make([]Lens, 0)
	}

	for _, p := range parts {
		h := Hasher(p)
		hash := h.Hash()
		log.Debugw("hashing", "in", h, "value", hash)
		tot += hash
	}
	log.Infow("answer", "part1", tot)

	for _, p := range parts {
		lens := ToLens(p)
		hash := lens.Label.Hash()
		log.Debugw("lens", "label", lens.Label, "hash", hash, "value", lens.Value)

		if lens.Value < 0 {
			if len(boxes[hash]) == 0 {
				continue // nothing to do.
			}
			boxes[hash] = slice.Filter(boxes[hash], func(keep Lens) bool {
				return keep.Label != lens.Label
			})
		} else {
			// placement operation
			idx := -1
			for i, l := range boxes[hash] {
				if l.Label == lens.Label {
					idx = i
					break
				}
			}
			if idx >= 0 {
				boxes[hash][idx] = lens // replace
			} else {
				boxes[hash] = append(boxes[hash], lens)
			}
		}
	}

	part2 := 0
	for i, b := range boxes {
		if len(b) > 0 {
			log.Debugw("box", "id", i, "lens", b)
		}
		for li, l := range b {
			part2 += ((i + 1) * ((li + 1) * l.Value))
		}
	}
	log.Infow("answer", "part2", part2)
}
