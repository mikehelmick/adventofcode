package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	ex1 := flag.Bool("e1", false, "-e1 to run example 1")
	ex2 := flag.Bool("e2", false, "-e2 to run example 2")
	part2 := flag.Bool("p2", false, "-p2 to pass part2 flag to binary (may or may not support)")

	day := flag.Int("d", 1, "-d X for day X")
	flag.Parse()

	goPath, err := exec.LookPath("go")
	if err != nil {
		panic(err)
	}

	dayPath := fmt.Sprintf("./day%02d", *day)
	baseArgs := []string{goPath, "run", dayPath}
	if *part2 {
		baseArgs = append(baseArgs, "-p2")
	}

	filePath := filepath.Join(dayPath, "input.txt")
	if *ex1 {
		filePath = filepath.Join(dayPath, "example1.txt")
	} else if *ex2 {
		filePath = filepath.Join(dayPath, "example2.txt")
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		panic(fmt.Sprintf("cannot read input file: %q %v", filePath, err))
	}
	input := strings.NewReader(string(data))

	cmd := exec.Cmd{
		Path:   goPath,
		Args:   baseArgs,
		Stdin:  input,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
