package main

import (
	"bufio"
	"fmt"
	"hash/maphash"
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

type Game struct {
	Player1, Player2 []int
	PreviousStates   map[[2]uint64]bool
	Hash             *maphash.Hash
	Recursive        bool
}

func (g *Game) Turn() bool {
	var a, b int
	a, g.Player1 = Pop(g.Player1)
	b, g.Player2 = Pop(g.Player2)

	// Assumes there's never a duplicate card. The rules don't specify
	// that case
	if a > b {
		g.Player1 = PushBack(g.Player1, a)
		g.Player1 = PushBack(g.Player1, b)
	} else {
		g.Player2 = PushBack(g.Player2, b)
		g.Player2 = PushBack(g.Player2, a)
	}

	return !((len(g.Player1) == 0) || (len(g.Player2) == 0))
}

func (g *Game) Scores() (int, int) {
	return Score(g.Player1), Score(g.Player2)
}

func NewGame(player1, player2 []int, recursive bool) *Game {
	if recursive {
		panic("recursive not implemented")
	}

	g := &Game{
		Player1:        make([]int, len(player1)),
		Player2:        make([]int, len(player2)),
		PreviousStates: make(map[[2]uint64]bool),
		Recursive:      recursive,
	}
	copy(g.Player1, player1)
	copy(g.Player2, player2)

	return g
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

	g := NewGame(hands[0], hands[1], false)

	for g.Turn() {
		// Keep going
	}

	score1, score2 := g.Scores()

	fmt.Println(score1, score2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
