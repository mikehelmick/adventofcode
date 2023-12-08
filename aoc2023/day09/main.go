package main

import (
	"bufio"
	"context"
	"os"

	"github.com/mikehelmick/adventofcode/pkg/logging"
)

type Data struct {
}

func main() {
	ctx := logging.WithLogger(context.Background(), logging.DefaultLogger())
	log := logging.FromContext(ctx)

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		log.Debugw("parsed line", "line", line)
	}

	if err := scanner.Err(); err != nil {
		log.Errorw("read error", "err", err)
	}
}
