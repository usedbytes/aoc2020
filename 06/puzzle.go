package main

import (
	"bufio"
	"fmt"
	"os"
)

type Form [26]bool

func (f *Form) Answer(r rune) {
	idx := r - 'a'
	f[idx] = true
}

func (f *Form) Sum() int {
	sum := 0
	for _, v := range f {
		if v {
			sum++
		}
	}

	return sum
}

type CombineFunc func(a, b bool) bool

func And(a, b bool) bool {
	return a && b
}

func Or(a, b bool) bool {
	return a || b
}

func (f *Form) Combine(other Form, op CombineFunc) {
	for i, v := range other {
		f[i] = op(f[i], v)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return fmt.Errorf("Usage: %s INPUT [PART2]", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	combineOp := Or
	if len(os.Args) > 2 {
		combineOp = And
	}

	totalAnswers := 0

	var group *Form

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			totalAnswers += group.Sum()
			group = nil
		} else {
			var individual Form
			for _, a := range line {
				individual.Answer(a)
			}

			if group == nil {
				group = &individual
			} else {
				group.Combine(individual, combineOp)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// If there's no blank line at EOF
	totalAnswers += group.Sum()

	fmt.Println(totalAnswers)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
