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

	var head, tail *Node

	cups := make([]int, len(input))
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
		tail = node
	}
	numCups := len(cups)
	minVal := min(cups)
	maxVal := max(cups)

	maxNum := 1000000
	numMoves := 10000000

	for i := maxVal+1; i <= maxNum; i++ {
		node := &Node{
			Prev:  tail,
			Next:  nil,
			Value: i,
		}
		tail.Next = node
		tail = node
		numCups++
	}
	maxVal = numCups

	tail.Next = head
	head.Prev = tail

	current := head

	//fmt.Println("current:", current.Value, "prev:", current.Prev.Value)

	//lastSearchLen := 0

	window := 100

	//PrintN(head, 20)
	//PrintN(tail, 20)

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

		//fmt.Println("Looking for", destinationVal)

		found := false
		lower, upper := current, current
		for !found {
			cursor := lower
			for i := 0; i < window; i++ {
				//fmt.Println(cursor.Value)
				if cursor.Value == destinationVal {
					destination = cursor
					found = true
					//fmt.Println("Searched", i)
					break
				}
				cursor = cursor.Prev
				lower = cursor
			}
			cursor = upper
			for i := 0; !found && i < window; i++ {
				//fmt.Println(cursor.Value)
				if cursor.Value == destinationVal {
					destination = cursor
					found = true
					//fmt.Println("Searched", i)
					break
				}
				cursor = cursor.Next
				upper = cursor
			}
			if !found {
				//fmt.Println("Not found. Expanding window", window, window * 2)
				window *= 2
				if window > numCups/2 {
					window = numCups/2
				}
			}
		}
		Insert(destination, removed)
		current = current.Next

		if move % 1000 == 0 {
			fmt.Println(move)
			PrintN(current, 20)
		}
	}

	//Print(current, current)
	one := Find(current, 1)
	PrintN(one, 10)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
