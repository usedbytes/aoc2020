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
	tail.Next = head
	head.Prev = tail

	Print(head, head)

	minVal := min(cups)
	maxVal := max(cups)

	current := head

	for move := 0; move < 100; move++ {
		var destination *Node
		removed := Remove(current, 3)
		destinationVal := current.Value - 1

		found := false
		for !found {
			cursor := current
			for i := 0; i < len(cups)-3; i++ {
				if cursor.Value == destinationVal {
					destination = cursor
					found = true
					break
				}
				cursor = cursor.Next
			}
			if !found {
				destinationVal--
				if destinationVal < minVal {
					destinationVal = maxVal
				}
			}
		}
		Insert(destination, removed)
		current = current.Next
	}

	Print(current, current)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
