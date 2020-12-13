package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func isUpper(r rune) bool {
	return r == 'B' || r == 'R'
}

func BinarySegment(segments string, setSize int) (int, error) {
	start := 0
	for _, s := range segments {
		setSize /= 2
		if isUpper(s) {
			start += setSize
		}
	}

	if setSize > 1 {
		return 0, fmt.Errorf("not enough segments to product answer")
	} else if setSize < 1 {
		return 0, fmt.Errorf("too many segments to product answer")
	}

	return start, nil
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

	seats := make([]int, 0, 1000)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		row, err := BinarySegment(line[:7], 128)
		if err != nil {
			return err
		}

		col, err := BinarySegment(line[7:10], 8)
		if err != nil {
			return err
		}

		seatId := row*8 + col
		seats = append(seats, seatId)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	sort.Ints(seats)
	fmt.Println("Number of seats:", seats[len(seats)-1])

	for i, s := range seats {
		if s+1 != seats[i+1] {
			fmt.Println("My seat:", s+1)
			break
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
