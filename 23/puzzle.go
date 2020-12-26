package main

import (
	"fmt"
	"os"
	"strconv"
)

type Node struct {
	Next, Prev *Node
	Value      int
}

func min(list []int) int {
	min := int(^uint(0) >> 1)
	for _, v := range list {
		if v < min {
			min = v
		}
	}
	return min
}

func max(list []int) int {
	max := 0
	for _, v := range list {
		if v > max {
			max = v
		}
	}
	return max
}

func Remove(after *Node, num int) *Node {
	removed := after.Next
	cursor := removed
	for i := 0; i < num; i++ {
		cursor = cursor.Next
	}
	(cursor.Prev).Next = nil
	after.Next = cursor
	cursor.Prev = after

	return removed
}

func Insert(after *Node, nodes *Node) {
	next := after.Next

	after.Next = nodes
	nodes.Prev = after

	cursor := nodes
	for cursor.Next != nil {
		cursor = cursor.Next
	}
	next.Prev = cursor
	cursor.Next = next
}

func Print(list, current *Node) {
	cursor := list
	for {
		if cursor == current {
			fmt.Printf(" (%d) ", cursor.Value)
		} else {
			fmt.Printf("  %d  ", cursor.Value)
		}
		cursor = cursor.Next

		if cursor == nil || cursor == list {
			break
		}
	}
	fmt.Println("")
}

func PrintN(list *Node, n int) {
	cursor := list
	for i := 0; i < n; i++ {
		fmt.Printf("  %d  ", cursor.Value)
		cursor = cursor.Next

		if cursor == nil || cursor == list {
			break
		}
	}
	fmt.Println("")
}

func Find(list *Node, value int) *Node {
	cursor := list
	for {
		if cursor.Value == value {
			break
		}

		cursor = cursor.Next

		if cursor == nil || cursor == list {
			cursor = nil
			break
		}
	}

	return cursor
}

func run() error {
	input := os.Args[1]
	part2 := len(os.Args) > 2

	// Assume min is 1 and max is len(input) (or 1000000 for Part 2)
	minVal := 1
	maxVal := len(input)
	numMoves := 100
	if part2 {
		maxVal = 1000000
		numMoves = 10000000
	}

	var head, tail *Node

	cups := make([]int, len(input))
	// lut maps from Value (index) to a *Node, for O(1) lookups
	lut := make([]*Node, maxVal+1)
	for i := range input {
		n, err := strconv.Atoi(input[i : i+1])
		if err != nil {
			return err
		}
		cups[i] = n
		node := &Node{
			Prev:  tail,
			Next:  nil,
			Value: n,
		}
		if head == nil {
			head = node
		} else {
			tail.Next = node
		}
		lut[n] = node
		tail = node
	}

	// Pad out the list if we need to
	for i := len(input) + 1; i <= maxVal; i++ {
		node := &Node{
			Prev:  tail,
			Next:  nil,
			Value: i,
		}
		lut[node.Value] = node
		tail.Next = node
		tail = node
	}

	// Loop the list to itself, so now we have a ring
	tail.Next = head
	head.Prev = tail

	// Start with the first cup
	current := head

	for move := 0; move < numMoves; move++ {
		var destination *Node
		removed := Remove(current, 3)
		destinationVal := current.Value - 1
		if destinationVal < minVal {
			destinationVal = maxVal
		}

		// Check if the value we're looking for is in the _short_ list
		inRemoved := removed
		for inRemoved != nil {
			inRemoved = Find(removed, destinationVal)
			if inRemoved != nil {
				destinationVal--
				if destinationVal < minVal {
					destinationVal = maxVal
				}
			}
		}

		// Look up the destination cup
		destination = lut[destinationVal]

		// Insert the removed cups
		Insert(destination, removed)

		// Move on
		current = current.Next
	}

	one := lut[1]
	if !part2 {
		Print(one, current)
	} else {
		fmt.Println(one.Next.Value * one.Next.Next.Value)
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
