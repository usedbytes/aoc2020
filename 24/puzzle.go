package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

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
var dirs map[string][2][2]int = map[string][2][2]int{ // [(dx_even, dy_even), (dx_odd, dy_odd)]
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

func CountNeighbours(floor map[[2]int]bool, coord [2]int, black bool) int {
	count := 0
	for _, dxdys := range dirs {
		dxdy := dxdys[0]
		// Odd y coord
		if coord[1]%2 != 0 {
			dxdy = dxdys[1]
		}

		neighbour := [2]int{coord[0]+dxdy[0], coord[1]+dxdy[1]}
		//fmt.Println("Neighbour", neighbour, floor[neighbour])
		if floor[neighbour] == black {
			count++
		}
	}

	return count
}

func ScanAndFlip(floor map[[2]int]bool) map[[2]int]bool{
	visited := map[[2]int]bool{}
	newFloor := map[[2]int]bool{}

	for coord, black := range floor {
		if visited[coord] {
			continue
		}
		visited[coord] = true

		if !black {
			continue
		}

		// Count adjacent, if 0 or >2, flip
		//fmt.Println("black coord", coord)
		numBlack := CountNeighbours(floor, coord, true)
		//fmt.Println("coord", coord, "numBlack", numBlack)
		if !((numBlack == 0) || (numBlack > 2)) {
			// Stays black
			newFloor[coord] = true
		} // else - don't add white tiles to newFloor


		// Then visit all the non-visited adjacent ones, to see
		// if they're white. We're only interested in white tiles
		// which are touching black ones, so this should find them all
		for _, dxdys := range dirs {
			dxdy := dxdys[0]
			// Odd y coord
			if coord[1]%2 != 0 {
				dxdy = dxdys[1]
			}

			neighbour := [2]int{coord[0]+dxdy[0], coord[1]+dxdy[1]}

			if visited[neighbour] {
				// Skip if already visited
				continue
			}
			visited[neighbour] = true

			black := floor[neighbour]

			// Skip if black
			if black {
				continue
			}

			// If white, check the white rules
			//fmt.Println("white neighbour", neighbour)
			numBlack := CountNeighbours(floor, neighbour, true)
			//fmt.Println("neighbour", neighbour, "numBlack", numBlack)
			if numBlack == 2 {
				newFloor[neighbour] = true
			} // else - don't add white tiles to newFloor
		}
	}

	return newFloor
}

func run() error {

	lobby := map[[2]int]bool{}
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
	fmt.Println("Part1:", count)

	for day := 0; day < 100; day++ {
		lobby = ScanAndFlip(lobby)

		count := 0
		for _, v := range lobby {
			if v {
				count++
			}
		}
		fmt.Println("Day", day+1, count)
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
