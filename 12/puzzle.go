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
		err := do(line)
		if err != nil {
			return err
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

type Heading int

const (
	HeadingEast  Heading = 0
	HeadingNorth         = 1
	HeadingWest          = 2
	HeadingSouth         = 3
)

type Ship struct {
	X, Y    int
	Heading Heading
}

func (s *Ship) Execute(c Command) {
	switch c.Opcode {
	case 'N':
		s.Y += c.Arg
	case 'E':
		s.X += c.Arg
	case 'S':
		s.Y -= c.Arg
	case 'W':
		s.X -= c.Arg
	case 'L':
		quarters := c.Arg / 90
		h := (int(s.Heading) + quarters) % 4
		s.Heading = Heading(h)
	case 'R':
		quarters := c.Arg / 90
		h := (int(s.Heading) - quarters) % 4
		if h < 0 {
			h += 4
		}
		s.Heading = Heading(h)
	case 'F':
		switch s.Heading {
		case HeadingNorth:
			s.Execute(Command{Opcode: 'N', Arg: c.Arg})
		case HeadingEast:
			s.Execute(Command{Opcode: 'E', Arg: c.Arg})
		case HeadingSouth:
			s.Execute(Command{Opcode: 'S', Arg: c.Arg})
		case HeadingWest:
			s.Execute(Command{Opcode: 'W', Arg: c.Arg})
		default:
			panic(fmt.Sprintln("Unknown heading", s.Heading))
		}
	default:
		panic(fmt.Sprintln("Unknown opcode", c.Opcode))
	}
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func (s *Ship) Manhattan() int {
	return abs(s.X) + abs(s.Y)
}

type Command struct {
	Opcode byte
	Arg    int
}

func run() error {
	ship := &Ship{}
	if err := doLines(os.Args[1], func(line string) error {
		arg, err := strconv.Atoi(line[1:])
		if err != nil {
			return err
		}
		cmd := Command{Opcode: line[0], Arg: arg}
		ship.Execute(cmd)

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(ship.Manhattan())

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
