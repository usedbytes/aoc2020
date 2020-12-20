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

type Grid struct {
	MinX, MinY int
	Cells      [][]byte
	Next       [][]byte
}

func (g *Grid) CountAround(x, y int, val byte, ignoreCentre bool) int {
	dxs := []int{-1, 0, 1}
	dys := []int{-1, 0, 1}
	count := 0
	y -= g.MinY
	x -= g.MinX
	for _, dy := range dys {
		if y+dy < 0 || y+dy >= len(g.Cells) {
			continue
		}

		row := g.Cells[y+dy]
		for _, dx := range dxs {
			if x+dx < 0 || x+dx >= len(row) {
				continue
			}
			if ignoreCentre && dy == 0 && dx == 0 {
				continue
			}

			if row[x+dx] == val {
				count++
			}
		}
	}

	return count
}

func (g *Grid) Print() {
	fmt.Println("X:", g.MinX, "Y:", g.MinY)
	for _, row := range g.Cells {
		fmt.Println(string(row))
	}
}

func (g *Grid) Flip() {
	for i, row := range g.Next {
		copy(g.Cells[i], row)
	}
}

func (g *Grid) Count(val byte) int {
	count := 0
	for _, row := range g.Cells {
		for _, cell := range row {
			if cell == val {
				count++
			}
		}
	}

	return count
}

func set(a []byte, b byte) {
	for i := range a {
		a[i] = b
	}
}

func (g *Grid) Expand(newMin, newMax int) {
	newSize := (newMax - newMin) + 1
	// X and Y size/offset always the same
	offset := g.MinY - newMin
	newCells := make([][]byte, newSize)
	newNext := make([][]byte, newSize)

	//fmt.Println("Grid Expand ->", newMin, newMax)
	//fmt.Println("Offset", offset)
	//fmt.Println("NewSize", newSize)

	// Fill in the new bits below the current contents
	if offset > 0 {
		for i := 0; i < offset; i++ {
			newCells[i] = make([]byte, newSize)
			set(newCells[i], '.')

			newNext[i] = make([]byte, newSize)
			set(newNext[i], '.')
		}
	}

	// Copy over the current contents
	for i := 0; i < len(g.Cells); i++ {
		newCells[i+offset] = make([]byte, newSize)
		set(newCells[i+offset], '.')
		copy(newCells[i+offset][offset:], g.Cells[i])

		newNext[i+offset] = make([]byte, newSize)
		set(newNext[i+offset], '.')
		copy(newNext[i+offset][offset:], g.Next[i])

	}

	// Fill in the new bits after the current contents
	if newMax >= g.MinY+len(g.Cells) {
		for i := offset + len(g.Cells); i < newSize; i++ {
			newCells[i] = make([]byte, newSize)
			set(newCells[i], '.')

			newNext[i] = make([]byte, newSize)
			set(newNext[i], '.')
		}
	}

	g.MinX = newMin
	g.MinY = newMin
	g.Cells = newCells
	g.Next = newNext
	//fmt.Println("Grid Expand. Now min:", g.MinX, g.MinY, "Size:", len(g.Cells))
}

func (g *Grid) Set(x, y int, v byte) {
	//fmt.Println("Grid Set(", x, y, string(v), ")")
	localX := x - g.MinX
	localY := y - g.MinY
	// These conditions depend on X and Y size/offset always being equal
	if localY < 0 {
		//fmt.Println("localY small", localY)
		g.Expand(y, g.MinY+len(g.Next))
	} else if localY >= len(g.Next) {
		//fmt.Println("localY big", localY)
		g.Expand(g.MinY, y)
	} else if localX < 0 {
		//fmt.Println("localX small", localX)
		g.Expand(x, g.MinX+len(g.Next[localY]))
	} else if localX >= len(g.Next[localY]) {
		//fmt.Println("localX big", localX)
		g.Expand(g.MinX, x)
	} else {
		g.Next[localY][localX] = v
		// Done, don't try again
		return
	}

	// Keep trying until it's big enough
	g.Set(x, y, v)
}

func (g *Grid) Get(x, y int, def byte) byte {
	localX := x - g.MinX
	localY := y - g.MinY
	if localY < 0 || localY >= len(g.Next) || localX < 0 || localX >= len(g.Next[localY]) {
		return def
	} else {
		return g.Cells[localY][localX]
	}

	return def
}

type Grid3D struct {
	MinZ   int
	MinX, MaxX, MinY, MaxY int
	Planes []*Grid
}

func (g *Grid3D) CountAround(x, y, z int, val byte, ignoreCentre bool) int {
	dzs := []int{-1, 0, 1}

	localZ := z - g.MinZ
	count := 0
	for _, dz := range dzs {
		if localZ+dz < 0 || localZ+dz >= len(g.Planes) {
			continue
		}
		ignoreCentre = (dz == 0)
		count += g.Planes[localZ+dz].CountAround(x, y, val, ignoreCentre)
	}

	return count
}

func (g *Grid3D) Count(val byte) int {
	count := 0
	for _, p := range g.Planes {
		pcount := p.Count(val)
		//fmt.Println("plane count", i, pcount)
		count += pcount
	}
	return count
}

func (g *Grid3D) Print() {
	for i, plane := range g.Planes {
		fmt.Println("Z:", i+g.MinZ)
		plane.Print()
		fmt.Println("--")
	}
}

func (g *Grid3D) Flip() {
	for _, plane := range g.Planes {
		plane.Flip()
	}
}

func (g *Grid3D) Expand(newMin, newMax int) {
	newSize := (newMax - newMin) + 1
	offset := g.MinZ - newMin
	newPlanes := make([]*Grid, newSize)

	// Fill in the new bits below the current contents
	if offset > 0 {
		for i := 0; i < offset; i++ {
			newPlanes[i] = &Grid{}
		}
	}

	// Copy over the current contents
	for i := 0; i < len(g.Planes); i++ {
		newPlanes[i+offset] = g.Planes[i]
	}

	// Fill in the new bits after the current contents
	if newMax >= g.MinZ+len(g.Planes) {
		for i := offset + len(g.Planes); i < newSize; i++ {
			newPlanes[i] = &Grid{}
		}
	}

	g.MinZ = newMin
	g.Planes = newPlanes

	//fmt.Println("Grid3D Expand. Now min:", g.MinZ, "Size:", len(g.Planes))
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

func (g *Grid3D) Set(x, y, z int, v byte) {
	//fmt.Println("Grid3D.Set(", x, y, z, string(v), ")")
	localZ := z - g.MinZ

	g.MaxX = max(x, g.MaxX)
	g.MinX = min(x, g.MinX)
	g.MaxY = max(y, g.MaxY)
	g.MinY = min(y, g.MinY)

	if localZ < 0 {
		g.Expand(z, g.MinZ+len(g.Planes))
	} else if localZ >= len(g.Planes) {
		g.Expand(g.MinZ, z)
	} else {
		g.Planes[localZ].Set(x, y, v)
		// Done, don't try again
		return
	}

	// Try again with the new size
	g.Set(x, y, z, v)
}

func (g *Grid3D) Get(x, y, z int, def byte) byte {
	localZ := z - g.MinZ
	if localZ < 0 {
		return def
	} else if localZ >= len(g.Planes) {
		return def
	} else {
		return g.Planes[localZ].Get(x, y, def)
	}

	return def
}

func run() error {

	grid := &Grid3D{ }

	y := 0
	if err := doLines(os.Args[1], func(line string) error {
		for x, c := range []byte(line) {
			grid.Set(x, y, 0, c)
		}
		y++
		return nil
	}); err != nil {
		return err
	}

	grid.Flip()
	grid.Print()

	for cycle := 0; cycle < 6; cycle++ {
		minZ := grid.MinZ-1
		maxZ := grid.MinZ+len(grid.Planes)
		minY := grid.MinY-1
		maxY := grid.MaxY+1
		minX := grid.MinX-1
		maxX := grid.MaxX+1
		//fmt.Printf("Cycle: %d, Z: (%d)-(%d), Y: (%d)-(%d), X: (%d)-(%d)\n", cycle, minZ, maxZ, minY, maxY, minX, maxX)
		for z := minZ; z <= maxZ; z++ {
			for y := minY; y <= maxY; y++ {
				for x := minX; x <= maxX; x++ {
					//action := "none"
					current := grid.Get(x, y, z, '.')
					count := grid.CountAround(x, y, z, '#', true)
					if current == '#' {
						if (count == 2 || count == 3) {
							grid.Set(x, y, z, '#')
							//action = "keep"
						} else {
							grid.Set(x, y, z, '.')
							//action = "kill"
						}
					} else if current == '.' {
						if (count == 3) {
							grid.Set(x, y, z, '#')
							//action = "spawn"
						}
					}
					//fmt.Printf("(%d, %d, %d) count: %d, action: %s\n", z, y, x, count, action)
				}
			}
		}
		grid.Flip()
		//fmt.Println("After Cycle", cycle)
		//grid.Print()
		//fmt.Println("-----")
	}

	//grid.Print()
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
