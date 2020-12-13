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

	sort.Ints(adapters)
	fmt.Println("I have", len(adapters), "adapters")

	numOneJolt := 0
	numThreeJolt := 0
	currentJolts := 0
	deviceJolts := adapters[len(adapters)-1] + 3
	adapters = append(adapters, deviceJolts)

	for _, a := range adapters {
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

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
