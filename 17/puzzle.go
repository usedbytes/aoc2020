package main

import (
	"bufio"
	"fmt"
	"os"
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

func set(a []byte, b byte) {
	for i := range a {
		a[i] = b
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Grid1D struct {
	// Range is [Min, Max)
	Min, Max     int
	Values       []byte
	defaultValue byte
}

type GridND struct {
	// Range is [Min, Max)
	Dims         int
	Min, Max     int
	Values       []Gridder
	defaultValue byte
}

func NewGrid(dimensions, min, max int, defaultValue byte) Gridder {
	var grid Gridder
	if dimensions == 1 {
		grid = NewGrid1D(min, max, defaultValue)
	} else {
		gnd := &GridND{
			Dims:         dimensions,
			Min:          min,
			Max:          max,
			Values:       make([]Gridder, max-min+1),
			defaultValue: defaultValue,
		}
		for i := range gnd.Values {
			gnd.Values[i] = NewGrid(dimensions-1, min, max, defaultValue)
		}
		grid = gnd
	}

	return grid
}

func NewGrid1D(min, max int, defaultValue byte) Gridder {
	g := &Grid1D{
		Min:          min,
		Max:          max,
		Values:       make([]byte, max-min+1),
		defaultValue: defaultValue,
	}

	set(g.Values, g.defaultValue)

	return g
}

func (g *Grid1D) Dup() Gridder {
	dup := &Grid1D{
		Min:          g.Min,
		Max:          g.Max,
		Values:       make([]byte, len(g.Values)),
		defaultValue: g.defaultValue,
	}
	copy(dup.Values, g.Values)
	return dup
}

func (g *GridND) Dup() Gridder {
	dup := &GridND{
		Dims:         g.Dims,
		Min:          g.Min,
		Max:          g.Max,
		Values:       make([]Gridder, len(g.Values)),
		defaultValue: g.defaultValue,
	}
	for i := range g.Values {
		dup.Values[i] = g.Values[i].Dup()
	}

	return dup
}

func (g *Grid1D) Set(coords []int, val byte) {
	if len(coords) != 1 {
		panic("wrong number of coords for Grid1D")
	}

	x := coords[0]

	//fmt.Printf("Grid%dD, Min: %d, Max: %d, Set(%v, %s)\n", 1, g.Min, g.Max, coords, string(val))

	if x < g.Min || x > g.Max || len(g.Values) == 0 {
		//fmt.Println("Realloc")
		newMin := min(g.Min, x)
		newMax := max(g.Max, x)
		newSize := newMax - newMin + 1

		newValues := make([]byte, newSize)
		set(newValues, g.defaultValue)
		offset := g.Min - newMin
		copy(newValues[offset:], g.Values)

		g.Min = newMin
		g.Max = newMax
		g.Values = newValues

		//fmt.Printf("Realloc Grid%dD, Min: %d, Max: %d, Size: %d, Set(%v, %s)\n", 1, g.Min, g.Max, newSize, coords, string(val))
	}

	g.Values[x-g.Min] = val
}

func (g *GridND) Set(coords []int, val byte) {
	if len(coords) != g.Dims {
		panic(fmt.Sprintf("wrong number of coords for Grid, need %d got %d", g.Dims, len(coords)))
	}

	//fmt.Printf("Grid%dD, Min: %d, Max: %d, Set(%v, %s)\n", g.Dims, g.Min, g.Max, coords, string(val))
	first := coords[0]

	if first < g.Min || first > g.Max || len(g.Values) == 0 {
		//fmt.Println("Realloc")
		newMin := min(g.Min, first)
		newMax := max(g.Max, first)
		newSize := newMax - newMin + 1

		newValues := make([]Gridder, newSize)
		offset := g.Min - newMin
		copy(newValues[offset:], g.Values)

		for i := 0; i < offset; i++ {
			newValues[i] = NewGrid(g.Dimensions()-1, 0, 0, g.defaultValue)
		}
		for i := offset + len(g.Values); i < newSize; i++ {
			newValues[i] = NewGrid(g.Dimensions()-1, 0, 0, g.defaultValue)
		}

		g.Min = newMin
		g.Max = newMax
		g.Values = newValues
		//fmt.Printf("Realloc Grid%dD, Min: %d, Max: %d, Size: %d, Set(%v, %s)\n", g.Dims, g.Min, g.Max, newSize, coords, string(val))
	}

	g.Values[first-g.Min].Set(coords[1:], val)
}

func (g *Grid1D) Get(coords []int, defaultVal byte) byte {
	if len(coords) != 1 {
		panic("wrong number of coords for Grid1D")
	}

	x := coords[0]
	if x < g.Min || x > g.Max || len(g.Values) == 0 {
		return defaultVal
	}
	return g.Values[x-g.Min]
}

func (g *GridND) Get(coords []int, defaultVal byte) byte {
	if len(coords) != g.Dims {
		panic(fmt.Sprintf("wrong number of coords for Grid, need %d got %d", g.Dims, len(coords)))
	}

	//fmt.Printf("Grid%dD, Min: %d, Max: %d, Get(%v, %s)\n", g.Dims, g.Min, g.Max, coords, string(defaultVal))
	first := coords[0]
	if first < g.Min || first > g.Max || len(g.Values) == 0 {
		//fmt.Println("OOB")
		return defaultVal
	}
	return g.Values[first-g.Min].Get(coords[1:], defaultVal)
}

func (g *Grid1D) CountAround(coords []int, val byte, ignoreCentre bool) int {
	ddim := []int{-1, 0, 1}
	count := 0
	x := coords[0] - g.Min
	for _, dx := range ddim {
		if x+dx < 0 || x+dx >= len(g.Values) {
			continue
		}
		if ignoreCentre && dx == 0 {
			continue
		}

		if g.Values[x+dx] == val {
			count++
		}
	}
	return count
}

func (g *GridND) CountAround(coords []int, val byte, ignoreCentre bool) int {
	if len(coords) != g.Dims {
		panic(fmt.Sprintf("wrong number of coords for Grid, need %d got %d", g.Dims, len(coords)))
	}

	ddim := []int{-1, 0, 1}
	count := 0
	first := coords[0] - g.Min
	for _, dx := range ddim {
		if first+dx < 0 || first+dx >= len(g.Values) {
			continue
		}
		count += g.Values[first+dx].CountAround(coords[1:], val, (dx == 0) && ignoreCentre)
	}
	return count
}

func (g *Grid1D) Count(val byte) int {
	count := 0
	for _, cell := range g.Values {
		if cell == val {
			count++
		}
	}
	return count
}

func (g *GridND) Count(val byte) int {
	count := 0
	for _, cell := range g.Values {
		count += cell.Count(val)
	}
	return count
}

func (g *Grid1D) String() string {
	return fmt.Sprintf("%2d: %s", g.Min, string(g.Values))
}

func (g *GridND) String() string {
	s := fmt.Sprintf("%d dim: %d\n", g.Dims, g.Min)
	for i, val := range g.Values {
		s += val.String()
		if i < len(g.Values)-1 {
			s += "\n"
		}
	}
	return s
}

func (g *Grid1D) Dimensions() int {
	return 1
}

func (g *GridND) Dimensions() int {
	return g.Dims
}

func (g *Grid1D) Range() [][2]int {
	return [][2]int{
		{g.Min, g.Max},
	}
}

func (g *GridND) Range() [][2]int {
	dimRanges := make([][2]int, g.Dims)

	var childRanges [][][2]int
	for _, child := range g.Values {
		childRanges = append(childRanges, child.Range())
	}

	dimRanges[0][0] = g.Min
	dimRanges[0][1] = g.Max
	for dim := 0; dim < g.Dims-1; dim++ {
		dimRanges[dim+1][0] = 0
		dimRanges[dim+1][1] = 0
		for _, child := range childRanges {
			dimRanges[dim+1][0] = min(child[dim][0], dimRanges[dim+1][0])
			dimRanges[dim+1][1] = max(child[dim][1], dimRanges[dim+1][1])
		}
	}

	return dimRanges
}

type Gridder interface {
	Set(coords []int, val byte)
	Get(coords []int, defaultVal byte) byte
	CountAround(coords []int, val byte, ignoreCentre bool) int
	Count(val byte) int
	Dup() Gridder
	String() string
	Dimensions() int
	Range() [][2]int
}

func run() error {

	grid := NewGrid(3, 0, 0, '.')

	y := 0
	if err := doLines(os.Args[1], func(line string) error {
		for x, c := range []byte(line) {
			grid.Set([]int{0, y, x}, c)
		}
		y++
		return nil
	}); err != nil {
		return err
	}

	fmt.Println("Starting configuration:")
	fmt.Println(grid.String())
	for cycle := 0; cycle < 6; cycle++ {
		ranges := grid.Range()
		next := grid.Dup()
		for z := ranges[0][0] - 1; z <= ranges[0][1]+1; z++ {
			for y := ranges[1][0] - 1; y <= ranges[1][1]+1; y++ {
				for x := ranges[2][0] - 1; x <= ranges[2][1]+1; x++ {
					current := grid.Get([]int{z, y, x}, '.')
					count := grid.CountAround([]int{z, y, x}, '#', true)
					//action := "none"
					if current == '#' {
						if count == 2 || count == 3 {
							next.Set([]int{z, y, x}, '#')
							//action = "keep"
						} else {
							next.Set([]int{z, y, x}, '.')
							//action = "kill"
						}
					} else if current == '.' {
						if count == 3 {
							next.Set([]int{z, y, x}, '#')
							//action = "spawn"
						}
					}
					//fmt.Printf("(%d, %d, %d) count: %d, action: %s\n", z, y, x, count, action)
				}
			}

		}
		grid = next
		//fmt.Println("Cycle", cycle)
		//fmt.Println(next.String())
		//fmt.Println("-----")
	}

	fmt.Println(grid.Count('#'))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
