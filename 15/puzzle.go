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

	numTurns, err := strconv.Atoi(os.Args[2])
	if err != nil {
		return err
	}

	lastTimes := make(map[int]int)

	turn := 1
	prev := 0

	for _, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return err
		}

		lastTimes[n] = turn
		prev = n
		turn++
	}

	for ; turn <= numTurns; turn++ {
		n := 0
		if lT, ok := lastTimes[prev]; ok {
			n = (turn - 1) - lT
		}
		// Note: Only update the map _after_ we've calculated the gap
		lastTimes[prev] = turn - 1

		prev = n
	}

	fmt.Println(prev)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
