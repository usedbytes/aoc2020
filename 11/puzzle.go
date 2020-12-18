package main

import (
	"bufio"
	"fmt"
	"os"
)

func doLines(filename string, do func(line string)) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		do(line)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type Grid struct {
	Cells [][]rune
	Next  [][]rune
}

func (g *Grid) CountAround(c, r int, val rune) int {
	count := 0
	for y := r - 1; y <= r+1; y++ {
		if y < 0 || y >= len(g.Cells) {
			continue
		}

		row := g.Cells[y]
		for x := c - 1; x <= c+1; x++ {
			if x < 0 || x >= len(row) {
				continue
			}
			if x == c && y == r {
				// Skip the target
				continue
			}

			if row[x] == val {
				count++
			}
		}
	}
	return count
}

func (g *Grid) Print() {
	for _, row := range g.Cells {
		fmt.Println(string(row))
	}
}

func (g *Grid) Flip() (int, bool) {
	defer func() {
		tmp := g.Cells
		g.Cells = g.Next
		g.Next = tmp
	}()

	flux := false
	occupied := 0
	for i, row := range g.Cells {
		for j, cell := range row {
			next := g.Next[i][j]
			if next != cell {
				flux = true
			}
			if next == '#' {
				occupied++
			}
		}
	}

	return occupied, flux
}

func run() error {
	grid := Grid{
		Cells: make([][]rune, 0),
		Next:  make([][]rune, 0),
	}

	if err := doLines(os.Args[1], func(line string) {
		arr := make([]rune, len(line))
		next := make([]rune, len(line))
		for i, c := range line {
			arr[i] = c
			next[i] = c
		}
		grid.Cells = append(grid.Cells, arr)
		grid.Next = append(grid.Next, next)
	}); err != nil {
		return err
	}

	flux := true
	occupied := 0
	for flux {
		for y := 0; y < len(grid.Cells); y++ {
			for x := 0; x < len(grid.Cells[0]); x++ {
				if grid.Cells[y][x] == '.' {
					continue
				} else {
					count := grid.CountAround(x, y, '#')
					if grid.Cells[y][x] == 'L' && count == 0 {
						grid.Next[y][x] = '#'
					} else if grid.Cells[y][x] == '#' && count >= 4 {
						grid.Next[y][x] = 'L'
					} else {
						grid.Next[y][x] = grid.Cells[y][x]
					}
				}
			}
		}
		occupied, flux = grid.Flip()
	}

	fmt.Println(occupied)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
