package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func doLines(filename string, do func(line string) error) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if err := do(line); err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func run() error {
	strs := strings.Split(os.Args[1], ",")

	seq := make([]int, 0, 2020)
	numTimes := make(map[int]int)

	turn := 1
	for _, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		seq = append(seq, n)
		numTimes[n] = 1
		turn++
	}

	for ; turn <= 2020; turn++ {
		prev := seq[len(seq)-1]
		nT := numTimes[prev]
		n := 0
		if nT > 1 {
			var i int
			for i = len(seq) - 2; i >= 0 && seq[i] != prev; i-- {
				// pass
			}
			lastTime := i + 1
			n = len(seq) - lastTime
		}
		seq = append(seq, n)
		numTimes[n]++
	}

	fmt.Println(seq[len(seq)-1])

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
