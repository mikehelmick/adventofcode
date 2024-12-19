package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mikehelmick/adventofcode/pkg/logging"
	"github.com/mikehelmick/adventofcode/pkg/straid"
)

type operator int

const (
	ADD operator = iota
	MULTIPLY
	CONCAT
	WRAP
)

type OperatorVec []operator

func (o OperatorVec) Increment(allowConcat bool) error {
	canWrap := false
	for _, op := range o {
		if allowConcat && op != CONCAT {
			canWrap = true
			break
		} else if !allowConcat && op != MULTIPLY {
			canWrap = true
			break
		}
	}
	if !canWrap {
		return fmt.Errorf("Cannot increment")
	}

	for i := len(o) - 1; i >= 0; i-- {
		if o[i] == ADD {
			o[i] = MULTIPLY
			return nil
		}
		if o[i] == MULTIPLY {
			if allowConcat {
				o[i] = CONCAT
				return nil
			} else {
				o[i] = ADD
				continue
			}
		}
		if o[i] == CONCAT {
			o[i] = ADD
		}
	}
	return nil
}

func (o OperatorVec) String() string {
	str := ""
	for _, op := range o {
		switch op {
		case ADD:
			str += "+"
		case MULTIPLY:
			str += "*"
		}
	}
	return str
}

type Equation struct {
	Values   []int64
	Answer   int64
	Solution OperatorVec
}

func NewEquation(line string) Equation {
	eq := Equation{
		Values:   make([]int64, 0),
		Solution: make(OperatorVec, 0),
	}
	answerAndValues := strings.Split(line, ":")
	eq.Answer = straid.AsInt(answerAndValues[0])

	values := strings.Split(strings.TrimSpace(answerAndValues[1]), " ")
	for _, v := range values {
		eq.Values = append(eq.Values, straid.AsInt(v))
	}
	return eq
}

func (e *Equation) String() string {
	return fmt.Sprintf("%v: %v || %+v", e.Answer, e.Values, e.Solution)
}

func (e *Equation) Solvable(allowConcat bool) bool {
	operators := make(OperatorVec, len(e.Values)-1)

	for {
		acc := e.Values[0]
		for i, op := range operators {
			switch op {
			case ADD:
				acc += e.Values[i+1]
			case MULTIPLY:
				acc *= e.Values[i+1]
			case CONCAT:
				acc = straid.AsInt(fmt.Sprintf("%d%d", acc, e.Values[i+1]))
			}
			if acc > e.Answer {
				break
			}
		}

		if acc == e.Answer {
			e.Solution = operators
			return true
		}

		if err := operators.Increment(allowConcat); err != nil {
			return false
		}
	}
}

func main() {
	log := logging.DefaultLogger()
	scanner := bufio.NewScanner(os.Stdin)

	equations := make([]Equation, 0)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			break
		}
		equations = append(equations, NewEquation(line))
	}

	var part1 int64
	for _, eq := range equations {
		log.Debugw("Equation", "eq", eq)
		if eq.Solvable(false) {
			log.Infow("Solution", "eq", eq)
			part1 += eq.Answer
		}
	}
	log.Infow("answer", "part1", part1)

	var part2 int64
	for _, eq := range equations {
		log.Debugw("Equation", "eq", eq)
		if eq.Solvable(true) {
			log.Infow("Solution", "eq", eq)
			part2 += eq.Answer
		}
	}
	log.Infow("answer", "part2", part2)
}
