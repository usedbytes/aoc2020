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
}

func run() error {
	// Obviously storing the whole address space would be impractical,
	// but the file is only a few hundred lines long, so worst case
	// the map will have a few hundred entries.
	m := &Machine{
		Mem: make(map[uint64]uint64),
	}

	if err := doLines(os.Args[1], func(line string) error {
		if strings.HasPrefix(line, "mask = ") {
			mask := line[len("mask = "):]
			m.MaskSet = 0
			m.MaskClear = 0
			bit := uint64(1 << 35)
			for _, b := range mask {
				switch b {
				case '1':
					m.MaskSet |= bit
				case '0':
					m.MaskClear |= bit
				case 'X':
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

			m.Mem[addr] = (value & ^m.MaskClear) | m.MaskSet
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
