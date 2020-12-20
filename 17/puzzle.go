package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	Min, Max     int
	Values       []byte
	defaultValue byte
}

type GridND struct {
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

	if x < g.Min || x > g.Max || len(g.Values) == 0 {
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
	}

	g.Values[x-g.Min] = val
}

func (g *GridND) Set(coords []int, val byte) {
	if len(coords) != g.Dims {
		panic(fmt.Sprintf("wrong number of coords for Grid, need %d got %d", g.Dims, len(coords)))
	}

	first := coords[0]

	if first < g.Min || first > g.Max || len(g.Values) == 0 {
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

	first := coords[0]
	if first < g.Min || first > g.Max || len(g.Values) == 0 {
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

func apply(grid, next Gridder, ranges [][2]int, coords []int) {
	dim := len(coords)
	coords = append(coords, 0)
	for x := ranges[dim][0] - 1; x <= ranges[dim][1]+1; x++ {
		coords[dim] = x

		if dim == grid.Dimensions()-1 {
			current := grid.Get(coords, '.')
			count := grid.CountAround(coords, '#', true)
			if current == '#' {
				if count == 2 || count == 3 {
					next.Set(coords, '#')
				} else {
					next.Set(coords, '.')
				}
			} else if current == '.' {
				if count == 3 {
					next.Set(coords, '#')
				}
			}
		} else {
			apply(grid, next, ranges, coords)
		}
	}
}

func run() error {
	dims := 3
	if len(os.Args) > 2 {
		var err error
		dims, err = strconv.Atoi(os.Args[2])
		if err != nil {
			return err
		}
	}

	coords := make([]int, dims)
	grid := NewGrid(dims, 0, 0, '.')

	coords[len(coords)-2] = 0

	if err := doLines(os.Args[1], func(line string) error {
		for x, c := range []byte(line) {
			coords[len(coords)-1] = x
			grid.Set(coords, c)
		}
		coords[len(coords)-2]++
		return nil
	}); err != nil {
		return err
	}

	fmt.Println("Starting configuration:")
	fmt.Println(grid.String())
	for cycle := 0; cycle < 6; cycle++ {
		ranges := grid.Range()
		next := grid.Dup()
		apply(grid, next, ranges, nil)
		grid = next
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
