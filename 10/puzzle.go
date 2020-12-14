package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
)

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: %s INPUT", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	adapters := make([]int, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		n, err := strconv.Atoi(line)
		if err != nil {
			return err
		}

		adapters = append(adapters, n)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println("I have", len(adapters), "adapters")

	// Add outlet at 0 Jolts
	adapters = append(adapters, 0)
	sort.Ints(adapters)

	numOneJolt := 0
	numThreeJolt := 0
	currentJolts := 0
	deviceJolts := adapters[len(adapters)-1] + 3
	adapters = append(adapters, deviceJolts)

	for _, a := range adapters[1:] {
		switch a - currentJolts {
		case 1:
			numOneJolt++
		case 3:
			numThreeJolt++
		default:
			return fmt.Errorf("unexpected difference: %d - %d = %d", a, currentJolts, a-currentJolts)
		}
		currentJolts = a
	}

	fmt.Println("One Jolts:", numOneJolt, ", Three Jolts:", numThreeJolt, ", Product:", numOneJolt*numThreeJolt)

	// The general solution would be to recurisvely remove adapters and see
	// if the result is still valid, and explore the solution space that
	// way, but like the puzzle page says "there must be more than a
	// trillion valid ways", so I assume that's not a good idea (and that
	// sounds like a hint not to go for a recursive solution).
	//
	// So instead, let's look for adapters that can be removed in isolation
	// and still leave a valid configuration. Then we can treat groups of
	// those "removable" adapters as places where the chain can be broken.
	// Depending on the number of adapters in the group, there are a
	// different number of ways to remove some or all of those adapters
	// while maintaining a valid chain:
	//
	// Group of 1
	// ----------
	// For example with the sequence 26, 27, 28, 31 we can only remove 27.
	// If we remove 28, then that leaves 26, 27, __, 31 - which is a gap of
	// 4 which is invalid. That's a "group" of 1: We can only remove one
	// adapter. This represents TWO ways to make a valid configuration:
	//  1. Adapter 27 present
	//  2. Adapter 27 absent
	//
	// Group of 2
	// ----------
	// If we have 0, 1, 2, 3, 6, then we can remove 1, 2, or both and still
	// be left with valid configurations - a "group" of 2, with FOUR
	// possible configurations:
	//  1. 0, 1, 2, 3
	//  2. 0, _, 2, 3
	//  3. 0, 1, _, 3
	//  4. 0, _, _, 3
	//
	// Group of 3
	// ----------
	// 33, 36, 37, 38, 39, 40, 43 is a group of three. We can remove any of
	// 37, 38, 39, but not all three at once. This is a group of 3 and
	// represents SEVEN valid configurations (note it's not 8, because the
	// "all three absent" configuration is not valid.
	//
	// So to find the total number of valid configurations, we find all the
	// "groups", and multiply together all the combinations represented by
	// that "group" (1, 4, or 7)
	// That requires a single traversal of the list, so it's only O(N)

	numCombinations := 1
	numInGroup := 0
	for i := 1; i < len(adapters)-1; i++ {
		prev := adapters[i-1]
		current := adapters[i]
		next := adapters[i+1]

		if (next - prev) >= 1 && (next - prev) <= 3 {
			// Still valid, keep adding to the group
			numInGroup++
			fmt.Printf("--> %3d [%3d] %3d valid\n", prev, current, next)
		} else {
			// Invalid, so evaluate the group
			switch numInGroup {
			case 0:
				break
			case 1:
				numCombinations *= 2
			case 2:
				numCombinations *= 4
			case 3:
				numCombinations *= 7
			default:
				return fmt.Errorf("group size too large: %d", numInGroup)
			}

			fmt.Printf("--X %3d [%3d] %3d invalid. Block of %d\n", prev, current, next, numInGroup)

			// Reset the group
			numInGroup = 0
		}
	}

	fmt.Println("Number of combinations:", numCombinations)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
