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

func Pop(list []int) (int, []int) {
	return list[0], list[1:]
}

func PushBack(list []int, val int) []int {
	return append(list, val)
}

func Score(list []int) int {
	score := 0
	for i := 1; i <= len(list); i++ {
		score += list[len(list)-i] * i
	}

	return score
}

func run() error {
	hands := [][]int{}
	var hand []int

	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			return nil
		}

		if strings.HasPrefix(line, "Player") {
			if hand != nil {
				hands = append(hands, hand)
				hand = []int{}
			}

			return nil
		}

		n, err := strconv.Atoi(line)
		if err != nil {
			return err
		}

		hand = append(hand, n)

		return nil
	}); err != nil {
		return err
	}

	if hand != nil {
		hands = append(hands, hand)
		hand = []int{}
	}

	rounds := 0
	for len(hands[0]) > 0 && len(hands[1]) > 0 {
		var a, b int
		a, hands[0] = Pop(hands[0])
		b, hands[1] = Pop(hands[1])

		if a > b {
			hands[0] = PushBack(hands[0], a)
			hands[0] = PushBack(hands[0], b)
		} else {
			hands[1] = PushBack(hands[1], b)
			hands[1] = PushBack(hands[1], a)
		}
		rounds++
	}

	winner := hands[0]
	if len(winner) == 0 {
		winner = hands[1]
	}

	fmt.Println(Score(winner))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
