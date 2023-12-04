package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/go-functional/slice"
)

type Card struct {
	ID      int
	Numbers []string
	Winners map[string]bool
}

// Matches tells you how many matches there were for part 1.
func (c *Card) Matches() int {
	points := 0
	for _, n := range c.Numbers {
		if c.Winners[n] {
			points++
		}
	}
	return points
}

// Points is the value for part 1.
func (c *Card) Points() int {
	points := 0
	for _, n := range c.Numbers {
		if c.Winners[n] {
			if points == 0 {
				points = 1
				continue
			}
			points = points * 2
		}
	}
	return points
}

// String turns a card into a string for debugging.
func (c *Card) String() string {
	w := make([]string, 0, len(c.Winners))
	for k := range c.Winners {
		w = append(w, k)
	}
	return fmt.Sprintf("Card: %v n: %v w: %v p: %v", c.ID,
		strings.Join(c.Numbers, ","),
		strings.Join(w, ","), c.Points())
}

// NewCard pares a line of text into a card.
func NewCard(s string) *Card {
	p := strings.Split(s, ":")
	nums := strings.Split(p[1], "|")

	id := strings.TrimSpace(strings.Split(p[0], "d")[1])

	numbers := strings.Split(strings.TrimSpace(nums[0]), " ")
	winners := strings.Split(strings.TrimSpace(nums[1]), " ")

	winMap := make(map[string]bool, len(winners))
	for _, w := range winners {
		w = strings.TrimSpace(w)
		if w != " " {
			winMap[w] = true
		}
	}

	keep := make([]string, 0, len(numbers))
	for _, n := range numbers {
		if n != "" {
			keep = append(keep, n)
		}
	}
	numbers = slice.Map(keep, func(s string) string {
		return strings.TrimSpace(s)
	})

	idNum, err := strconv.Atoi(id)
	if err != nil {
		panic(err)
	}

	return &Card{
		ID:      idNum,
		Numbers: numbers,
		Winners: winMap,
	}
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	cards := make([]*Card, 0)
	cardMap := make(map[int]*Card)
	copies := make(map[int]int)
	part1 := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		card := NewCard(line)
		log.Debugw("loaded", "card", card.String())
		cards = append(cards, card)
		cardMap[card.ID] = card
		copies[card.ID] = 1
		part1 += card.Points()
	}

	log.Infow("answer", "part1", part1)

	cardCount := len(cards)
	// while the card count map isn't empty, seeded w/ 1 copy of every card
	for len(copies) > 0 {
		log.Debugw("scratching", "cards", len(copies))

		// Build the next wave from this wave.
		next := make(map[int]int)
		// For every card in the current wave, scratch it and seed that many cards for the next wave.
		for cid, count := range copies {
			// Scratch all the instances of a card that we have at once (count added in below)
			card := cardMap[cid]
			wins := card.Matches()
			// Add the wins from this scratch to the next round
			for i := 0; i < wins; i++ {
				// and increment cardCount
				cardCount += count
				wonCard := card.ID + 1 + i
				if cur, ok := next[wonCard]; !ok {
					next[wonCard] = count
				} else {
					next[wonCard] = cur + count
				}
			}

		}
		copies = next
	}
	log.Infow("answer", "part2", cardCount)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
