package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func scan(vals []int, accum, levels, target int, results []int) (bool, []int) {
	for i, v := range vals[:len(vals)-(levels-1)] {
		if levels > 1 {
			found, results := scan(vals[i+1:], accum+v, levels-1, target, results)
			if found {
				return true, append(results, v)
			}
		} else {
			if accum+v == target {
				return true, append(results, v)
			}
		}
	}

	return false, results
}

func run() error {
	if len(os.Args) != 3 {
		return fmt.Errorf("Usage: %s N_NUMBERS INPUT", os.Args[0])
	}

	n, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return err
	}

	f, err := os.Open(os.Args[2])
	if err != nil {
		return err
	}
	defer f.Close()

	vals := make([]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		val, err := strconv.Atoi(line)
		if err != nil {
			return err
		}

		vals = append(vals, val)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	found, results := scan(vals, 0, n, 2020, nil)
	if found {
		r := 1
		for _, v := range results {
			r *= v
		}
		fmt.Println(results, r)
		return nil
	}

	return fmt.Errorf("no results found")
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
