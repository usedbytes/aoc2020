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

	Orientation int
	Flip        bool
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

func Rotate90(image [][]byte) [][]byte {
	size := len(image)
	newImage := make([][]byte, size)
	for y := 0; y < size; y++ {
		newRow := make([]byte, size)
		for x := 0; x < size; x++ {
			newRow[x] = image[x][size-y-1]
		}
		newImage[y] = newRow
	}
	return newImage
}

func (t *Tile) Rotate90() {
	t.Borders[0], t.Borders[1], t.Borders[2], t.Borders[3] =
		t.Borders[1], t.Borders[2], t.Borders[3], t.Borders[0]

	t.Neighbours[0], t.Neighbours[1], t.Neighbours[2], t.Neighbours[3] =
		t.Neighbours[1], t.Neighbours[2], t.Neighbours[3], t.Neighbours[0]

	// It would be much more efficient to just store a transform and
	// then have a Content accessor that reads out with the right transform
	// But this is easier to think about
	t.Content = Rotate90(t.Content)
}

func HFlip(image [][]byte) [][]byte {
	size := len(image)
	newImage := make([][]byte, size)
	for y := 0; y < size; y++ {
		newRow := make([]byte, size)
		for x := 0; x < size; x++ {
			newRow[x] = image[y][size-x-1]
		}
		newImage[y] = newRow
	}
	return newImage
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

	// It would be much more efficient to just store a transform and
	// then have a Content accessor that reads out with the right transform
	// But this is easier to think about
	t.Content = HFlip(t.Content)
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

func SearchFor(monster [][]byte, image [][]byte, match, replace byte) int {
	count := 0
	for y := 0; y < len(image)-len(monster); y++ {
		for x := 0; x < len(image[0])-len(monster[0]); x++ {
			found := true
			for my := 0; my < len(monster); my++ {
				imgRow := image[y+my]
				monsterRow := monster[my]
				for mx := 0; mx < len(monsterRow); mx++ {
					if monsterRow[mx] == match && imgRow[x+mx] != match {
						found = false
						break
					}
				}
				if !found {
					break
				}
			}
			if found {
				// Found one, mark it
				count++
				for my := 0; my < len(monster); my++ {
					imgRow := image[y+my]
					monsterRow := monster[my]
					for mx := 0; mx < len(monsterRow); mx++ {
						if monsterRow[mx] == match {
							imgRow[x+mx] = replace
						}
					}
				}
			}
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
		if len(line) == 0 {
			tile.Borders = make([]*Border, 4)
			for i, b := range borders {
				border := &Border{}
				border.Parse(b)
				tile.Borders[i] = border
				borders[i] = ""
			}
			tiles[tile.ID] = tile
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
		xform(t)
	}

	tileImgSize := int(math.Sqrt(float64(len(tiles))))
	tileImage := make([][]*Tile, tileImgSize)
	for y := 0; y < tileImgSize; y++ {
		tileImage[y] = make([]*Tile, tileImgSize)
	}
	tileImage[0][0] = t

	// Then work through, transforming each neighbour until it fits
	for y := 0; y < len(tileImage); y++ {
		for x := 0; x < len(tileImage[0]); x++ {
			if tileImage[y][x] == nil {
				panic(fmt.Sprintln("not assigned yet", x, y))
			}
			t = tileImage[y][x]

			for i, nid := range t.Neighbours {
				if nid == 0 {
					continue
				}

				n := tiles[nid]
				nx, ny := x+dx[i], y+dy[i]
				opp := opposites[i]

				if tileImage[ny][nx] != nil {
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
					tileImage[ny][nx] = n
				}
			}
		}
	}

	// Now we have an image of tiles, generate a straightforward raster image
	image := make([][]byte, tileImgSize*(tileSize-2))
	for y := 0; y < len(tileImage); y++ {
		for ty := 0; ty < tileSize-2; ty++ {
			image[y*(tileSize-2)+ty] = make([]byte, tileImgSize*(tileSize-2))
			for x := 0; x < len(tileImage[0]); x++ {
				t := tileImage[y][x]
				for tx := 0; tx < tileSize-2; tx++ {
					image[y*(tileSize-2)+ty][x*(tileSize-2)+tx] = t.Content[ty][tx]
				}
			}
		}
	}

	imgTransforms := []func([][]byte) [][]byte{
		func(img [][]byte) [][]byte { return Rotate90(img) },
		func(img [][]byte) [][]byte { return Rotate90(img) },
		func(img [][]byte) [][]byte { return Rotate90(img) },
		func(img [][]byte) [][]byte { return HFlip(img) },
		func(img [][]byte) [][]byte { return Rotate90(img) },
		func(img [][]byte) [][]byte { return Rotate90(img) },
		func(img [][]byte) [][]byte { return Rotate90(img) },
		// Hax: Just to make sure we run the search on all orientations
		func(img [][]byte) [][]byte { return img },
	}
	monster := [][]byte{
		[]byte("                  # "),
		[]byte("#    ##    ##    ###"),
		[]byte(" #  #  #  #  #  #   "),
	}
	for _, xform := range imgTransforms {
		count := SearchFor(monster, image, '#', 'O')
		if count > 0 {
			fmt.Println("Found", count, "monsters")
			break
		}
		image = xform(image)
	}

	count := 0
	for y := 0; y < len(image); y++ {
		for x := 0; x < len(image[y]); x++ {
			if image[y][x] == '#' {
				count++
			}
		}
		fmt.Println(string(image[y]))
	}
	fmt.Println("Roughness:", count)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
