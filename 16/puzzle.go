package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Section int

const (
	RulesSection Section = iota
	MyTicketSection
	NearbyTicketsSection
	MaxSection
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

type Range struct {
	Min, Max int
}

func (r *Range) Valid(n int) bool {
	return n >= r.Min && n <= r.Max
}

func (r *Range) Parse(s string) error {
	minThenMax := strings.Split(s, "-")
	if len(minThenMax) != 2 {
		return fmt.Errorf("couldn't split range: %s", s)
	}

	tmp, err := strconv.Atoi(minThenMax[0])
	if err != nil {
		return err
	}
	r.Min = tmp

	tmp, err = strconv.Atoi(minThenMax[1])
	if err != nil {
		return err
	}
	r.Max = tmp

	return nil
}

func (r *Range) String() string {
	return fmt.Sprintf("%d-%d", r.Min, r.Max)
}

type Field struct {
	Name string
	A, B Range
}

func (f *Field) Valid(n int) bool {
	return f.A.Valid(n) || f.B.Valid(n)
}

func (f *Field) Parse(s string) error {
	fieldThenRanges := strings.Split(s, ": ")
	if len(fieldThenRanges) != 2 {
		return fmt.Errorf("couldn't split field: %s", s)
	}
	f.Name = fieldThenRanges[0]

	aThenB := strings.Split(fieldThenRanges[1], " or ")
	if len(aThenB) != 2 {
		return fmt.Errorf("couldn't split ranges: %s", fieldThenRanges[1])
	}

	err := f.A.Parse(aThenB[0])
	if err != nil {
		return err
	}

	err = f.B.Parse(aThenB[1])
	if err != nil {
		return err
	}

	return nil
}

func (f *Field) String() string {
	return fmt.Sprintf("%s: %s or %s", f.Name, f.A.String(), f.B.String())
}

type Ticket struct {
	Values []int
}

func (t *Ticket) Parse(s string) error {
	vals := strings.Split(s, ",")
	for _, v := range vals {
		n, err := strconv.Atoi(v)
		if err != nil {
			return err
		}
		t.Values = append(t.Values, n)
	}
	return nil
}

func CountBitsSet(val uint32) (int, int) {
	lastBit := 0
	n := 0
	for bit := 0; bit < 32; bit++ {
		if val&(1<<bit) != 0 {
			n++
			lastBit = bit
		}
	}
	return lastBit, n
}

func run() error {

	section := RulesSection
	fields := make([]*Field, 0)
	var myTicket Ticket
	tickets := make([]*Ticket, 0)

	if err := doLines(os.Args[1], func(line string) error {
		//fmt.Println(line)
		if len(line) == 0 {
			section++
		}

		switch section {
		case RulesSection:
			// Parse rule
			f := &Field{}
			err := f.Parse(line)
			if err != nil {
				return err
			}
			fields = append(fields, f)
		case MyTicketSection:
			if line == "your ticket:" {
				return nil
			}
			// Parse ticket
			myTicket.Parse(line)
		case NearbyTicketsSection:
			if line == "nearby tickets:" {
				return nil
			}
			ticket := &Ticket{}
			ticket.Parse(line)
			tickets = append(tickets, ticket)
		default:
			return fmt.Errorf("too many sections")
			// Parse ticket into nearby
		}

		return nil
	}); err != nil {
		return err
	}

	// Part 1
	scanningErrorRate := 0
	notInvalidTickets := make([]*Ticket, 0, len(tickets))
	for _, t := range tickets {
		ticketCouldBeValid := true
		for _, v := range t.Values {
			fieldCouldBeValid := false
			for _, f := range fields {
				if f.Valid(v) {
					fieldCouldBeValid = true
				}
			}

			// XXX: Can we quit on the first invalid? It's not clear.
			if !fieldCouldBeValid {
				scanningErrorRate += v
				ticketCouldBeValid = false
			}
		}
		if ticketCouldBeValid {
			notInvalidTickets = append(notInvalidTickets, t)
		}
	}
	fmt.Println(scanningErrorRate)
	fmt.Println("tickets", len(tickets), "valid", len(notInvalidTickets))

	// Part 2
	// We have to assume that all remaining tickets are actually valid.
	numValues := len(myTicket.Values)
	if numValues > 32 {
		panic("not enough bits in bitmask")
	}
	if numValues != len(fields) {
		panic("different numbers of fields and values")
	}
	possibleIndices := make([]uint32, numValues)

	// possibleIndices[i] has a bit set for each index in Ticket.Values
	// which _could_ represent field[i].
	// We start with all bits set, and then clear them if we ever encounter
	// an invalid Value for field[i] at that position
	for i := range possibleIndices {
		possibleIndices[i] = uint32(uint64(1<<numValues) - 1)
	}
	for _, t := range notInvalidTickets {
		for vi, v := range t.Values {
			for fi, f := range fields {
				if !f.Valid(v) {
					possibleIndices[fi] &= ^uint32(1 << vi)
				}
			}
		}
	}

	// Attempt to assign the fields to their indices, one by one, reducing
	// ambiguity as we go
	indices := make([]int, numValues)
	for assigned := 0; assigned < numValues; assigned++ {
		found := false
		for i, mask := range possibleIndices {
			if bit, n := CountBitsSet(mask); n == 1 {
				// If only one bit is set, then field[i] must be at
				// index 'bit'
				indices[i] = bit
				// No other field can be at this index, so clear this
				// bit in all of the masks.
				for j := range possibleIndices {
					possibleIndices[j] &= ^mask
				}
				found = true
				break
			}
		}
		if !found {
			panic("no unambigious bit found")
		}
	}

	// Now we know that field[i] is at index indices[i] in Ticket.Values

	result := 1
	for fi, f := range fields {
		if strings.HasPrefix(f.Name, "departure") {
			idx := indices[fi]
			result *= myTicket.Values[idx]
		}
	}
	fmt.Println(result)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
