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

type Border struct {
	Size           int
	Pattern        uint32
	ReversePattern uint32
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
	s2 := ""
	for i := 0; i < b.Size; i++ {
		if b.ReversePattern&(1<<i) != 0 {
			s2 = "#" + s2
		} else {
			s2 = "." + s2
		}
	}
	s += " (" + s2 + ")"
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
	Borders []*Border

	Neighbours map[int]int
}

func (t *Tile) String() string {
	s := ""
	s += fmt.Sprintf("Tile: %d\n", t.ID)
	s += fmt.Sprintf("Borders: %s, %s, %s, %s\n",
		t.Borders[0].String(), t.Borders[1].String(), t.Borders[2].String(), t.Borders[3].String())
	return s
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
			tile = &Tile{
				ID:         id,
				Neighbours: make(map[int]int),
			}
		} else {
			if tileLine == 1 {
				tileSize = len(line)
				borders[0] = line
			}
			borders[1] = borders[1] + string(line[len(line)-1])
			borders[3] = string(line[0]) + borders[3]
			if tileLine == tileSize {
				borders[2] = Reverse(line)
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
		if len(t.Neighbours) == 2 {
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

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
