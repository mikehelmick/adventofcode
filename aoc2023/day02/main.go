package main

import (
	"bufio"
	"context"
	"os"

	"github.com/mikehelmick/adventofcode/aoc2023/pkg/logging"
)

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		log.Infow("read", "line", line)
		log.Debug("debug message")
	}
	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
