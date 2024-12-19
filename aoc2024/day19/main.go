package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/go-functional/slice"
)

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	line := scanner.Text()

	towels := strings.Split(line, ",")
	towels = slice.Map(towels, func(pattern string) string {
		return strings.TrimSpace(pattern)
	})
	log.Debugw("loaded towels", "towels", towels)

	patterns := make([]string, 0, 100)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		patterns = append(patterns, line)
	}
	log.Debug("loaded patterns", "patterns", patterns)

	part1 := 0
	part2 := 0
	cache := newCache()
	for _, pattern := range patterns {
		log.Debugw("checking pattern", "pattern", pattern)
		matches := canMatchPattern(cache, towels, pattern)
		if matches > 0 {
			log.Infow("pattern matches", "pattern", pattern)
			part2 += matches
			part1++
		} else {
			log.Infow("pattern does not match", "pattern", pattern)
		}

	}
	log.Infow("answer", "part1", part1)
	log.Infow("answer", "part2", part2)
}

type patternCache struct {
	mu    sync.RWMutex
	cache map[string]int
}

func newCache() *patternCache {
	return &patternCache{
		cache: make(map[string]int),
	}
}

func (p *patternCache) Check(pattern string) (int, error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if v, ok := p.cache[pattern]; ok {
		return v, nil
	}
	return 0, fmt.Errorf("unknown)")
}

func (p *patternCache) Set(pattern string, value int) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.cache[pattern] = value
}

func canMatchPattern(patternCache *patternCache, towels []string, pattern string) int {
	if pattern == "" {
		return 0
	}
	if res, err := patternCache.Check(pattern); err == nil {
		return res
	}

	sum := 0
	for _, towel := range towels {
		if pattern == towel {
			patternCache.Set(pattern, 1)
			sum++
		} else if strings.HasPrefix(pattern, towel) {
			res := canMatchPattern(patternCache, towels, strings.TrimPrefix(pattern, towel))
			patternCache.Set(pattern, res)
			sum += res
		}
	}
	patternCache.Set(pattern, sum)
	return sum
}
