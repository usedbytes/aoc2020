package main

import (
	"bufio"
	"fmt"
	"os"
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

	lobby := map[[2]int]bool{}
	// It's a hexagonal grid, so I think we can just treat it as
	// a rectilinear grid, with "strange" dx/dys, and they need to be different
	// for odd vs even rows.
	//
	//  *     *     *     *     *     *
	//    \NW     /NE    E       W
	//     *     *     *---->*<----*
	// Let's define (dx, dy) for even rows:
	//  NW: ( 0, -1)
	//  NE: ( 1, -1)
	//  E:  ( 1,  0)
	//  W:  (-1,  0)
	//  SW: ( 0,  1)
	//  SE: ( 1,  1)
	//
	// Then for odd rows, the dxs need to be different, because NW->NE should
	// result in the same X coord
	//  NW: ( -1, -1)
	//  NE: ( 0, -1)
	//  E:  ( 1,  0)
	//  W:  (-1,  0)
	//  SW: ( -1,  1)
	//  SE: ( 0,  1)
	dirs := map[string][2][2]int{ // (dx, dy)
		"nw": [2][2]int{
			{0, -1},
			{-1, -1},
		},
		"ne": [2][2]int{
			{1, -1},
			{0, -1},
		},
		"e": [2][2]int{
			{1, 0},
			{1, 0},
		},
		"w": [2][2]int{
			{-1, 0},
			{-1, 0},
		},
		"sw": [2][2]int{
			{0, 1},
			{-1, 1},
		},
		"se": [2][2]int{
			{1, 1},
			{0, 1},
		},
	}

	if err := doLines(os.Args[1], func(line string) error {
		coord := [2]int{}
		for len(line) > 0 {
			for k, vs := range dirs {
				v := vs[0]
				// Odd y coord
				if coord[1]%2 != 0 {
					v = vs[1]
				}
				if strings.HasPrefix(line, k) {
					coord[0], coord[1] = coord[0]+v[0], coord[1]+v[1]
					line = line[len(k):]
				}
			}
		}
		current := lobby[coord]
		lobby[coord] = !current
		return nil
	}); err != nil {
		return err
	}

	count := 0
	for _, v := range lobby {
		if v {
			count++
		}
	}
	fmt.Println(count)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
