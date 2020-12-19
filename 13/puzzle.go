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
	start := -1
	buses := make([]int, 0)
	if err := doLines(os.Args[1], func(line string) error {
		if start == -1 {
			t, err := strconv.Atoi(line)
			if err != nil {
				return err
			}
			start = t
		} else {
			sbuses := strings.Split(line, ",")
			for _, s := range sbuses {
				if s == "x" {
					continue
				}
				b, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				buses = append(buses, b)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	min := start
	minB := 0
	for _, b := range buses {
		diff := b - (start % b)
		if diff < min {
			min = diff
			minB = b
		}
	}

	fmt.Println("Bus", minB, "leaves in", min, "minutes -", minB*min)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
