package main

import (
	"bufio"
	"fmt"
	"math"
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

type Border struct {
	Size           int
	Pattern        uint32
	ReversePattern uint32
}

func (b *Border) Flip() {
	b.Pattern, b.ReversePattern = b.ReversePattern, b.Pattern
}

func (b *Border) String() string {
	s := ""
	for i := 0; i < b.Size; i++ {
		if b.Pattern&(1<<i) != 0 {
			s = "#" + s
		} else {
			s = "." + s
		}
	}
	return s
}

func (b *Border) Parse(s string) {
	b.Size = len(s)
	b.Pattern = 0
	b.ReversePattern = 0
	for i, c := range s {
		if c == '#' {
			b.Pattern |= 1
			b.ReversePattern |= (1 << (b.Size - 1))
		}
		if i < len(s)-1 {
			b.Pattern <<= 1
			b.ReversePattern >>= 1
		}
	}
}

// Tile:
//     MSB.........LSB
// LSB       /|\ N(0)  MSB
//   . W(3)   |        .
//   .   <----|---->   .
//   .        |  E(1)  .
// MSB  S(2) \|/       LSB
//     LSB.........MSB
type Tile struct {
	ID      int
	Size    int
	Borders []*Border

	Neighbours [4]int
	Content    [][]byte
}

func (t *Tile) String() string {
	bstr := [4]string{
		t.Borders[0].String(),
		t.Borders[1].String(),
		t.Borders[2].String(),
		t.Borders[3].String(),
	}

	idStr := fmt.Sprintf("%d", t.ID)
	s := ""
	for y := 0; y < t.Size; y++ {
		row := make([]byte, t.Size)
		for x := 0; x < len(row); x++ {
			if y == 0 {
				row[x] = bstr[0][x]
			} else if y == t.Size-1 {
				row[x] = bstr[2][t.Size-x-1]
			} else {
				if x == 0 {
					row[x] = bstr[3][t.Size-y-1]
				} else if x == t.Size-1 {
					row[x] = bstr[1][y]
				} else {
					if y == t.Size/2 && (x-3) >= 0 && (x-3) < len(idStr) {
						row[x] = idStr[x-3]
					} else {
						row[x] = t.Content[y-1][x-1]
					}
				}
			}
		}
		s += string(row) + "\n"
	}

	/*
		s += fmt.Sprintf("Tile: %d\n", t.ID)
		s += fmt.Sprintf("Borders: %s, %s, %s, %s\n",
			t.Borders[0].String(), t.Borders[1].String(), t.Borders[2].String(), t.Borders[3].String())
		s += fmt.Sprintf("Neighbours: N: %d, E: %d, S: %d, W: %d\n",
				t.Neighbours[0], t.Neighbours[1], t.Neighbours[2], t.Neighbours[3])
	*/
	return s
}

func (t *Tile) Rotate90() {
	t.Borders[0], t.Borders[1], t.Borders[2], t.Borders[3] =
		t.Borders[1], t.Borders[2], t.Borders[3], t.Borders[0]

	t.Neighbours[0], t.Neighbours[1], t.Neighbours[2], t.Neighbours[3] =
		t.Neighbours[1], t.Neighbours[2], t.Neighbours[3], t.Neighbours[0]
}

func (t *Tile) HFlip() {
	t.Borders[1], t.Borders[3] =
		t.Borders[3], t.Borders[1]

	t.Borders[0].Flip()
	t.Borders[1].Flip()
	t.Borders[2].Flip()
	t.Borders[3].Flip()

	t.Neighbours[1], t.Neighbours[3] =
		t.Neighbours[3], t.Neighbours[1]
}

func (t *Tile) NumNeighbours() int {
	count := 0
	for _, nid := range t.Neighbours {
		if nid != 0 {
			count++
		}
	}
	return count
}

func Reverse(s string) string {
	ret := ""
	for _, c := range s {
		ret = string(c) + ret
	}
	return ret
}

func run() error {

	tiles := make(map[int]*Tile)
	var tile *Tile
	var tileLine int
	var tileSize int
	var borders [4]string
	if err := doLines(os.Args[1], func(line string) error {
		fmt.Println(line)
		if len(line) == 0 {
			tile.Borders = make([]*Border, 4)
			for i, b := range borders {
				border := &Border{}
				border.Parse(b)
				tile.Borders[i] = border
				borders[i] = ""
			}
			tiles[tile.ID] = tile
			fmt.Println(tile)
			tile = nil
			tileLine = 0

			return nil
		}

		if tile == nil {
			var id int
			n, err := fmt.Sscanf(line, "Tile %d:", &id)
			if n != 1 {
				return fmt.Errorf("couldn't scan tile ID")
			} else if err != nil {
				return err
			}

			if id == 0 {
				panic("can't handle id == 0")
			}

			tile = &Tile{
				ID: id,
			}
		} else {
			if tileLine == 1 {
				tileSize = len(line)
				tile.Size = tileSize
				borders[0] = line
				tile.Content = make([][]byte, tileSize-2)
			}
			borders[1] = borders[1] + string(line[len(line)-1])
			borders[3] = string(line[0]) + borders[3]
			if tileLine == tileSize {
				borders[2] = Reverse(line)
			}
			if tileLine > 1 && tileLine < tileSize {
				row := make([]byte, tileSize-2)
				for i, c := range line[1 : tileSize-1] {
					row[i] = byte(c)
				}
				tile.Content[tileLine-2] = row
			}
		}

		tileLine++

		return nil
	}); err != nil {
		return err
	}

	fmt.Printf("Read %d tiles\n", len(tiles))
	corners := make([]*Tile, 0, 4)

	for _, t := range tiles {
		for i, b := range t.Borders {
			for _, t2 := range tiles {
				if t == t2 {
					continue
				}
				for _, b2 := range t2.Borders {
					if b.Pattern == b2.Pattern || b.Pattern == b2.ReversePattern {
						t.Neighbours[i] = t2.ID
					}
				}
			}
		}
		if t.NumNeighbours() == 2 {
			corners = append(corners, t)
		}
	}

	if len(corners) != 4 {
		return fmt.Errorf("Couldn't find 4 corners")
	}

	// Part 1
	product := 1
	for _, c := range corners {
		product *= c.ID
	}
	fmt.Println(product)

	// Part 2
	transforms := []func(t *Tile){
		func(t *Tile) { t.Rotate90() },
		func(t *Tile) { t.Rotate90() },
		func(t *Tile) { t.Rotate90() },
		func(t *Tile) { t.HFlip() },
		func(t *Tile) { t.Rotate90() },
		func(t *Tile) { t.Rotate90() },
		func(t *Tile) { t.Rotate90() },
	}
	opposites := []int{2, 3, 0, 1}
	dx := []int{0, 1, 0, -1}
	dy := []int{-1, 0, 1, 0}

	// Let's start with putting corners[0] in the top-left
	t := corners[0]

	// We need to orient it so that the connecting sides are East and South
	for _, xform := range transforms {
		east := t.Neighbours[1]
		south := t.Neighbours[2]
		if east != 0 && south != 0 {
			break
		}
		fmt.Println("Transform")
		xform(t)
	}

	imgSize := int(math.Sqrt(float64(len(tiles))))
	image := make([][]*Tile, imgSize)
	for y := 0; y < imgSize; y++ {
		image[y] = make([]*Tile, imgSize)
	}
	image[0][0] = t

	// Then work through, transforming each neighbour until it fits
	for y := 0; y < len(image); y++ {
		for x := 0; x < len(image[0]); x++ {
			if image[y][x] == nil {
				panic(fmt.Sprintln("not assigned yet", x, y))
			}
			t = image[y][x]

			for i, nid := range t.Neighbours {
				if nid == 0 {
					continue
				}

				n := tiles[nid]
				nx, ny := x+dx[i], y+dy[i]
				opp := opposites[i]

				if image[ny][nx] != nil {
					// Neighbour already assigned, just check it is OK
					if t.Borders[i].Pattern != n.Borders[opp].ReversePattern {
						panic(fmt.Sprintf("neighbour not matching, %d side %d -> %d side %d", t.ID, i, n.ID, opp))
					}
				} else {
					for _, xform := range transforms {
						// Needs to match the reverse pattern
						if t.Borders[i].Pattern == n.Borders[opp].ReversePattern {
							break
						}
						xform(n)
					}
					if t.Borders[i].Pattern != n.Borders[opp].ReversePattern {
						panic("not matching after all possible transforms")
					}
					image[ny][nx] = n
				}
			}
		}
	}

	for y := 0; y < len(image); y++ {
		for ty := 0; ty < tileSize-2; ty++ {
			for x := 0; x < len(image[0]); x++ {
				t := image[y][x]
				for tx := 0; tx < tileSize-2; tx++ {
					fmt.Printf("%s", string(t.Content[ty][tx]))
				}
				fmt.Printf(" ")
			}
			fmt.Println("")
		}
		fmt.Println("")
	}

	/*
		fmt.Println(t)
		for i, nid := range t.Neighbours {
			n := tiles[nid]
			opp := opposites[i]
			for _, xform := range transforms {
				if t.Borders[i].Pattern == n.Borders[opp].Pattern {
					break
				}
				xform(n)
			}
			image[y+dy[i]][x+dx[i]] = n
			fmt.Println("Neighbour", i)
			fmt.Println(n)
		}
	*/

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
