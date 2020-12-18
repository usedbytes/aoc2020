package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

func (g *Grid) CastRay(fromX, fromY, dx, dy, maxSteps int) rune {
	x, y := fromX, fromY
	for steps := 0; maxSteps < 0 || steps < maxSteps; steps++ {
		x, y = x+dx, y+dy
		if y < 0 || y >= len(g.Cells) {
			return '.'
		}
		if x < 0 || x >= len(g.Cells[0]) {
			return '.'
		}
		cell := g.Cells[y][x]
		if cell != '.' {
			return cell
		}
	}

	return '.'
}

func (g *Grid) CountAround(c, r, distance int, val rune) int {
	rays := [8][2]int{
		{0, -1},
		{1, -1},
		{1, 0},
		{1, 1},
		{0, 1},
		{-1, 1},
		{-1, 0},
		{-1, -1},
	}
	count := 0
	for _, ray := range rays {
		hit := g.CastRay(c, r, ray[0], ray[1], distance)
		if hit == val {
			count++
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

	distance := 1
	threshold := 4
	if len(os.Args) > 2 {
		d, err := strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}
		distance = d
	}
	if len(os.Args) > 3 {
		t, err := strconv.Atoi(os.Args[3])
		if err != nil {
			return err
		}
		threshold = t
	}
	for flux {
		for y := 0; y < len(grid.Cells); y++ {
			for x := 0; x < len(grid.Cells[0]); x++ {
				if grid.Cells[y][x] == '.' {
					continue
				} else {
					count := grid.CountAround(x, y, distance, '#')
					if grid.Cells[y][x] == 'L' && count == 0 {
						grid.Next[y][x] = '#'
					} else if grid.Cells[y][x] == '#' && count >= threshold {
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
