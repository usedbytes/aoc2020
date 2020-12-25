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

type Game struct {
	Player1, Player2 []int
	PreviousStates   map[[2]string]bool
	Recursive        bool
}

func (g *Game) Play() int {
	for (len(g.Player1) > 0) && (len(g.Player2) > 0) {
		if g.Recursive {
			strings := [2]string{
				fmt.Sprintf("%v", g.Player1),
				fmt.Sprintf("%v", g.Player2),
			}

			if _, ok := g.PreviousStates[strings]; ok {
				// The game instantly ends with a win for player 1
				return 1
			}

			g.PreviousStates[strings] = true
		}

		var a, b int
		a, g.Player1 = Pop(g.Player1)
		b, g.Player2 = Pop(g.Player2)

		turnWinner := 0
		recurse := false
		if g.Recursive && (len(g.Player1) >= a) && (len(g.Player2) >= b) {
			recurse = true
		}

		if recurse {
			subGame := NewGame(g.Player1[:a], g.Player2[:b], true)
			turnWinner = subGame.Play()
		} else {
			// Assumes there's never a duplicate card. The rules don't specify
			// that case
			if a > b {
				turnWinner = 1
			} else {
				turnWinner = 2
			}
		}

		if turnWinner == 1 {
			g.Player1 = PushBack(g.Player1, a)
			g.Player1 = PushBack(g.Player1, b)
		} else if turnWinner == 2 {
			g.Player2 = PushBack(g.Player2, b)
			g.Player2 = PushBack(g.Player2, a)
		} else {
			panic("nobody won")
		}
	}

	if len(g.Player1) > 0 {
		return 1
	}

	return 2
}

func (g *Game) Scores() (int, int) {
	return Score(g.Player1), Score(g.Player2)
}

func NewGame(player1, player2 []int, recursive bool) *Game {
	g := &Game{
		Player1:        make([]int, len(player1)),
		Player2:        make([]int, len(player2)),
		PreviousStates: make(map[[2]string]bool),
		Recursive:      recursive,
	}
	copy(g.Player1, player1)
	copy(g.Player2, player2)

	return g
}

func run() error {
	hands := [][]int{}
	var hand []int

	recursive := len(os.Args) > 2

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

	g := NewGame(hands[0], hands[1], recursive)
	winner := g.Play()

	score1, score2 := g.Scores()

	fmt.Println(winner, score1, score2)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
