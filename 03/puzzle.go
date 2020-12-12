package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Slope struct {
	Right, Down int
}

func (s Slope) String() string {
	return fmt.Sprintf("[%d, %d]", s.Right, s.Down)
}

func run() error {
	if len(os.Args) < 3 {
		return fmt.Errorf("Usage: %s INPUT right,down...", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	slopes := make([]Slope, 0)

	for _, a := range os.Args[2:] {
		s := Slope{}

		n, err := fmt.Sscanf(a, "%d,%d", &s.Right, &s.Down)
		if n != 2 {
			return fmt.Errorf("couldn't parse argument as right,down pair: %s", a)
		} else if err != nil {
			return err
		}

		slopes = append(slopes, s)
	}

	product := 1

	for _, s := range slopes {
		f.Seek(0, 0)
		scanner := bufio.NewScanner(f)

		numTrees := 0
		x := 0
		skip := 0 // We always process the first line

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			if skip == 0 {
				tree := (line[x%len(line)] == '#')
				if tree {
					numTrees++
				}
				x += s.Right
				skip = s.Down
			}

			skip--
		}

		if err := scanner.Err(); err != nil {
			return err
		}

		fmt.Println(s, numTrees)

		product *= numTrees
	}

	fmt.Println(product)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
