package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

type XMAS struct {
	history []int
	cursor  int
}

func (x *XMAS) Receive(val int) {
	x.history[x.cursor] = val
	x.cursor++
	if x.cursor >= len(x.history) {
		x.cursor = 0
	}
}

func (x *XMAS) Valid(val int) bool {
	for i, v1 := range x.history {
		for _, v2 := range x.history[i:] {
			if v1+v2 == val {
				return true
			}
		}
	}
	return false
}

func NewXMAS(preambleLen int) *XMAS {
	return &XMAS{
		history: make([]int, preambleLen),
	}
}

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: %s INPUT", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	preambleLen := 25
	x := NewXMAS(preambleLen)
	i := 0
	n := 0

	message := make([]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		n, err = strconv.Atoi(line)
		if err != nil {
			return err
		}

		if i > preambleLen {
			if !x.Valid(n) {
				fmt.Println("Invalid:", n)
				break
			}
		}
		x.Receive(n)
		message = append(message, n)
		i++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	for i, v1 := range message {
		sum := v1
		min := v1
		max := v1
		for _, v2 := range message[i+1:] {
			if v2 < min {
				min = v2
			}
			if v2 > max {
				max = v2
			}

			sum += v2
			if sum == n {
				fmt.Println("min, max, sum", min, max, min+max)
			} else if sum > n {
				break
			}
		}
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
