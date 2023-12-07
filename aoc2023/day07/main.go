package main

import (
	"bufio"
	"context"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

var tieBreak = map[string]int{}
var tieBreakPart2 = map[string]int{}

func init() {
	in := "AKQJT98765432"
	in2 := "AKQT98765432J"
	for i := 0; i < len(in); i++ {
		tieBreak[in[i:i+1]] = i
		tieBreakPart2[in2[i:i+1]] = i
	}
}

type HandType int

const (
	// Represents the hand types, stronger is lower num.
	FIVE_KIND  = iota //0
	FOUR_KIND         // 1
	FULL_HOUSE        // 2
	THREE_KIND        // 3
	TWO_PAIR          // 4
	ONE_PAIR          // 5
	HIGH_CARD         // 6
)

type Camel struct {
	Hand string
	Bid  int
}

func (c *Camel) CardAt(i int) string {
	return c.Hand[i : i+1]
}

func (c *Camel) Hits() map[string]int {
	hits := make(map[string]int)
	for i := 0; i < len(c.Hand); i++ {
		card := c.Hand[i : i+1]
		if v, ok := hits[card]; !ok {
			hits[card] = 1
		} else {
			hits[card] = v + 1
		}
	}
	return hits
}

// Creates a camel hand from the input line
func NewCamel(line string) *Camel {
	parts := strings.Split(line, " ")
	bid, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return &Camel{Hand: parts[0], Bid: bid}
}

func Score2(c *Camel) HandType {
	jokers := 0
	for i := 0; i < len(c.Hand); i++ {
		if c.Hand[i:i+1] == "J" {
			jokers++
		}
	}
	if jokers == 0 || jokers == 5 {
		return Score(c)
	}
	if jokers == 4 {
		return FIVE_KIND
	}

	hits := c.Hits()

	if jokers == 3 {
		if len(hits) == 3 {
			// 3 jokers, 2 singles
			return FOUR_KIND
		}
		if len(hits) == 2 {
			// this is already a full house, but it is now five kind
			return FIVE_KIND
		}
	}
	if jokers == 2 {
		if len(hits) == 2 {
			// 2 jokers, 3 of a kind, upgraded to
			return FIVE_KIND
		}
		if len(hits) == 3 {
			return FOUR_KIND
		}
		if len(hits) == 4 {
			return THREE_KIND
		}
	}
	if jokers == 1 {
		if len(hits) == 2 {
			return FIVE_KIND
		}
		if len(hits) == 3 {
			for k, v := range hits {
				if k == "J" {
					continue
				}
				if v == 3 {
					return FOUR_KIND
				}
			}
			return FULL_HOUSE
		}
		if len(hits) == 4 {
			return THREE_KIND
		}
		if len(hits) == 5 {
			return ONE_PAIR
		}
	}

	panic("this is not possible")
}

func Score(c *Camel) HandType {
	hits := c.Hits()

	// easy cases
	if len(hits) == 5 {
		return HIGH_CARD
	}
	if len(hits) == 1 {
		return FIVE_KIND
	}

	if len(hits) == 2 {
		// could be four_kind, or full_house
		for _, v := range hits {
			if v == 4 {
				return FOUR_KIND
			}
		}
		return FULL_HOUSE
	}

	if len(hits) == 3 {
		// must be 3 of a kind or 2 pair
		for _, v := range hits {
			if v == 3 {
				return THREE_KIND
			}
		}
		return TWO_PAIR
	}

	// only thing left
	return ONE_PAIR
}

type ScoreFn func(*Camel) HandType

// Used to sort Camel hands based on a scoring function and associated tieBreak map.
func (c *Camel) StrongerThan(o *Camel, score ScoreFn, tieBreak map[string]int) bool {
	cScore := score(c)
	oScore := score(o)

	if cScore < oScore {
		return true
	}

	if cScore == oScore {
		for i := 0; i < len(c.Hand); i++ {
			cCard := tieBreak[c.CardAt(i)]
			oCard := tieBreak[o.CardAt(i)]
			if cCard == oCard {
				continue
			}
			if cCard < oCard {
				return true
			}
			return false
		}
	}

	return false
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	// Load input
	hands := make([]*Camel, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		hands = append(hands, NewCamel(line))
	}
	log.Debugw("loaded hands", "hands", hands)

	// Part 1, using first scoring function.
	sort.Slice(hands, func(i, j int) bool {
		return !hands[i].StrongerThan(hands[j], Score, tieBreak)
	})
	part1 := 0
	for i, h := range hands {
		log.Debugw("output", "rank", i+1, "hand", h.Hand, "bid", h.Bid, "score", Score(h))
		part1 += ((i + 1) * h.Bid)
	}
	log.Infow("answer", "part1", part1)

	// Part 2 using second scoring function.
	sort.Slice(hands, func(i, j int) bool {
		return !hands[i].StrongerThan(hands[j], Score2, tieBreakPart2)
	})
	part2 := 0
	for i, h := range hands {
		log.Debugw("output", "rank", i+1, "hand", h.Hand, "bid", h.Bid, "score", Score2(h))
		part2 += ((i + 1) * h.Bid)
	}
	log.Infow("answer", "part2", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
