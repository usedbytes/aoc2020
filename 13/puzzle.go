package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

// https://en.wikipedia.org/wiki/Extended_Euclidean_algorithm#Computing_multiplicative_inverses_in_modular_structures
func ModuloInverse(a, n int) int {
	t, newT := 0, 1
	r, newR := n, a

	for newR != 0 {
		q := r / newR
		t, newT = newT, t-q*newT
		r, newR = newR, r-q*newR
	}

	if r > 1 {
		panic(fmt.Sprintf("%d not invertible", a))
	}
	if t < 0 {
		t += n
	}

	return t
}

// http://homepages.math.uic.edu/~leon/mcs425-s08/handouts/chinese_remainder.pdf
func CRT(divisors, remainders []int) int {
	ms := divisors
	as := remainders
	zs := make([]int, len(ms))
	ys := make([]int, len(ms))
	ws := make([]int, len(ms))
	m := 1
	for i := range ms {
		z := 1
		m *= ms[i]
		for j, d := range ms {
			if i == j {
				continue
			}
			z *= d
		}
		zs[i] = z
	}
	x := 0
	for i := range ms {
		ys[i] = ModuloInverse(zs[i], ms[i])
		ws[i] = (ys[i] * zs[i]) % m

		x += as[i] * ws[i]
	}

	x %= m

	return x
}

func run() error {
	start := -1
	buses := make([]int, 0)
	minutesAfter := make([]int, 0)
	if err := doLines(os.Args[1], func(line string) error {
		if start == -1 {
			t, err := strconv.Atoi(line)
			if err != nil {
				return err
			}
			start = t
		} else {
			sbuses := strings.Split(line, ",")
			for i, s := range sbuses {
				if s == "x" {
					continue
				}
				b, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				buses = append(buses, b)
				minutesAfter = append(minutesAfter, i)
			}
		}

		return nil
	}); err != nil {
		return err
	}

	{
		// Part 1
		min := start
		minB := 0
		for _, b := range buses {
			diff := b - (start % b)
			if diff < min {
				min = diff
				minB = b
			}
		}

		fmt.Println("Bus", minB, "leaves in", min, "minutes -", minB*min)
	}
	{

		// Part 2
		// A dumb iterative solution takes too slow, even when trying
		// to optimise the step size.
		// Some googling around solving simultaneous modulo equations
		// turned up Chinese Remainder Theorem
		// Feels like a bit of a hollow victory. Apparently this is
		// a "very common", "very important" algorithm that comes up
		// as lot in crypto - but I've never heard of it

		remainders := make([]int, len(minutesAfter))
		for i := range minutesAfter {
			remainders[i] = buses[i] - minutesAfter[i]
		}

		t := CRT(buses, remainders)
		fmt.Println("By Chinese Remainder Theorem:", t)

		for i, b := range buses {
			mod := (t + minutesAfter[i]) % b

			fmt.Printf("(%d + %d) %% %d = %d\n", t, minutesAfter[i], b, mod)
		}

		// Alternative approach which breaks the problem up, which I
		// saw hints towards online.
		candidate := 0
		step := buses[0]
		for i := 1; i < len(buses); i++ {
			for x := candidate; ; x += step {
				if (x+minutesAfter[i])%buses[i] == 0 {
					// If this is the smallest step that satisfies
					// all buses up to buses[i], then to keep satisfying
					// those buses and this one, the step will be
					// multiplied by this bus' period
					// I guess  this only works if the buses are coprime
					// otherwise this wouldn't be the smallest step and
					// we would jump over some cases where they re-align
					candidate = x
					step *= buses[i]
					break
				}
			}
		}
		fmt.Println(candidate)
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
