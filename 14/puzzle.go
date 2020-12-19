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

type Machine struct {
	MaskSet   uint64
	MaskClear uint64
	Mem       map[uint64]uint64
	MaskX     []int
}

func (m *Machine) Write(addr, value uint64) {
	m.Mem[addr] = value
}

func run() error {
	// Obviously storing the whole address space would be impractical,
	// but the file is only a few hundred lines long, so worst case
	// the map will have a few hundred entries.
	//
	// Edit: lol @ Part2. Now we're writing a whole bunch more locations.
	// It still looks practical to track it as a map... but perhaps there's
	// a better way
	m := &Machine{
		Mem:   make(map[uint64]uint64),
		MaskX: make([]int, 0, 36),
	}

	part2 := len(os.Args) > 2
	if err := doLines(os.Args[1], func(line string) error {
		if strings.HasPrefix(line, "mask = ") {
			mask := line[len("mask = "):]
			m.MaskSet = 0
			m.MaskClear = 0
			m.MaskX = m.MaskX[:0]
			bit := uint64(1 << 35)
			for i, b := range mask {
				switch b {
				case '1':
					m.MaskSet |= bit
				case '0':
					m.MaskClear |= bit
				case 'X':
					m.MaskX = append(m.MaskX, 35-i)
					break
				default:
					panic(fmt.Sprintln("invalid mask bit", b))
				}
				bit >>= 1
			}
		} else {
			var addr, value uint64
			n, err := fmt.Sscanf(line, "mem[%d] = %d", &addr, &value)
			if n != 2 {
				return fmt.Errorf("couldn't parse write instruction %s", line)
			} else if err != nil {
				return err
			}

			if !part2 {
				m.Mem[addr] = (value & ^m.MaskClear) | m.MaskSet
			} else {
				// Mask = 0 means unchanged in part 2, only set bits.
				addr = addr | m.MaskSet
				for cnt := 0; cnt < (1 << len(m.MaskX)); cnt++ {
					maskAddr := addr
					for bit := 0; bit < len(m.MaskX); bit++ {
						// MaskX is actually in MSB-to-LSB order, but
						// it shouldn't matter because we're doing all
						// combinations
						if cnt&(1<<bit) != 0 {
							maskAddr |= (1 << m.MaskX[bit])
						} else {
							maskAddr &= ^(1 << m.MaskX[bit])
						}
					}
					m.Write(maskAddr, value)
				}
			}
		}
		return nil
	}); err != nil {
		return err
	}

	fmt.Println("Memory locations:", len(m.Mem))
	sum := uint64(0)
	for _, v := range m.Mem {
		sum += v
	}
	fmt.Println("sum", sum)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
