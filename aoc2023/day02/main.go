package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/mikehelmick/adventofcode/aoc2023/pkg/logging"
	"github.com/mikehelmick/go-functional/maps"
	"github.com/mikehelmick/go-functional/slice"
)

// Show is an individual set of cubes you were shown.
type Show struct {
	Cubes map[string]int
}

// Given a set of stones, was this show possible?
func (s *Show) Possible(bag map[string]int) bool {
	log := logging.DefaultLogger()
	for color, amt := range s.Cubes {
		// if we were shown more than the bag amount, not possible.
		if amt > bag[color] {
			log.Debugw("miss", "color", color, "shown", amt, "bag", bag[color])
			return false
		}
	}
	return true
}

// Game is an individual game.
type Game struct {
	ID    int
	Shows []Show
}

// Power is for part 2, the fewest stones of each color multiplied together.
func (g *Game) Power() int {
	return slice.FoldL(maps.ToSlice(g.Fewest()), 1,
		func(v int, acc int) int {
			return v * acc
		})
}

// Fewest is the fewest of each color stone that makes the game still possible.
func (g *Game) Fewest() map[string]int {
	return slice.FoldL(g.Shows, map[string]int{},
		func(s Show, rtn map[string]int) map[string]int {
			for color, count := range s.Cubes {
				if count > rtn[color] {
					rtn[color] = count
				}
			}
			return rtn
		})
}

// Possible says if a game is possible w/ the given amount of cubes.
func (g *Game) Possible(bag map[string]int) bool {
	return slice.All(g.Shows, func(s Show) bool { return s.Possible(bag) })
}

func (g *Game) String() string {
	return fmt.Sprintf("Game %d: %+v", g.ID, g.Shows)
}

func ParseLine(l string) *Game {
	parts := strings.Split(l, ":")

	gid := strings.Split(parts[0], " ")
	id, err := strconv.Atoi(gid[1])
	if err != nil {
		panic(err)
	}

	game := &Game{
		ID:    id,
		Shows: make([]Show, 0),
	}

	shown := strings.Split(parts[1], ";")
	for _, s := range shown {
		s = strings.TrimSpace(s)
		countColors := strings.Split(s, ",")

		show := Show{
			Cubes: make(map[string]int),
		}

		for _, cc := range countColors {
			cc = strings.TrimSpace(cc)
			countColor := strings.Split(cc, " ")
			count, err := strconv.Atoi(countColor[0])
			if err != nil {
				panic(err)
			}
			show.Cubes[countColor[1]] = count
		}
		game.Shows = append(game.Shows, show)
	}

	return game
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	games := make([]*Game, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		game := ParseLine(line)
		log.Debugw("loaded", "game", game)
		games = append(games, game)
	}

	// Part one, find which games are possible w/ this number of cubes.
	bag := map[string]int{"red": 12, "green": 13, "blue": 14}
	part1 := slice.FoldL(games, 0, func(game *Game, acc int) int {
		if game.Possible(bag) {
			return acc + game.ID
		}
		return acc
	})
	log.Infow("part1", "answer", part1)

	// part 2
	part2 := slice.FoldL(games, 0,
		func(g *Game, acc int) int {
			return acc + g.Power()
		})
	log.Infow("part2", "answer", part2)

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
