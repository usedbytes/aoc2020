package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
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

type Operation struct {
	Symbol string
	Order  int // Lower order == evaluated sooner
	Func   func(a, b int) int
}

var Operations map[rune]*Operation = map[rune]*Operation{}

type Node struct {
	Order        int
	Op           *Operation
	Value        int
	L, R, Parent *Node
}

func AddLeaf(newNode, current *Node) *Node {
	if current != nil {
		if current.L == nil {
			current.L = newNode
		} else if current.R == nil {
			current.R = newNode
		} else {
			panic("both children already assigned")
		}
	}
	newNode.Parent = current
	return newNode
}

func ReplaceNode(newNode, replace *Node) *Node {
	newNode.L = replace
	newNode.Parent = replace.Parent
	if replace.Parent != nil {
		if replace.Parent.R == replace {
			replace.Parent.R = newNode
		} else {
			replace.Parent.L = newNode
		}
	}
	replace.Parent = newNode

	return newNode
}

func Parse(s string) (*Node, int) {
	var (
		current   *Node
		i, skipTo int
		r         rune
	)
	for i, r = range s {
		if skipTo > i {
			// Can't figure out if there's a better way to advance than this
			continue
		}

		if unicode.IsDigit(r) {
			toks := strings.SplitN(s[i:], " ", 2)
			n, err := strconv.Atoi(strings.TrimRight(toks[0], ")"))
			if err != nil {
				panic(err)
			}

			newNode := &Node{
				Value: n,
			}

			// A value node always starts as a leaf (an Operation might re-order it later)
			current = AddLeaf(newNode, current)
		} else if op, ok := Operations[r]; ok {
			newNode := &Node{
				Op: op,
			}

			// Walk up the tree until we're at the root, or a later-order operation
			for ; current.Parent != nil && newNode.Op.Order >= current.Parent.Op.Order; current = current.Parent {
				continue
			}
			// Replace the node we landed at
			current = ReplaceNode(newNode, current)
		} else if strings.ContainsRune("(", r) {
			newNode, di := Parse(s[i+1:])

			// Parenthesised operations are always leaves (parentheses force early evaluation)
			current = AddLeaf(newNode, current)

			skipTo = i + di + 1
		} else if strings.ContainsRune(")", r) {
			i++
			break
		}
	}

	// Walk back to the root
	for ; current.Parent != nil; current = current.Parent {
		continue
	}

	return current, i
}

func (n *Node) String() string {
	if n.Op == nil {
		return fmt.Sprintf("%d", n.Value)
	}

	return fmt.Sprintf("(%s %s %s)", n.L.String(), n.Op.Symbol, n.R.String())
}

func (n *Node) Eval() int {
	if n.Op == nil {
		return n.Value
	}

	return n.Op.Func(n.L.Eval(), n.R.Eval())
}

func run() error {
	// For Part 1, we just evaluate left-to-right, so multiplication and
	// addition have the same precedence (0)
	// For Part 2, addition comes first, so we increase multiplication
	// Order to 1
	mulOrder := 0
	if len(os.Args) > 2 {
		mulOrder = 1
	}

	Operations['+'] = &Operation{
		Symbol: "+",
		Order:  0,
		Func: func(a, b int) int {
			return a + b
		},
	}
	Operations['*'] = &Operation{
		Symbol: "*",
		Order:  mulOrder,
		Func: func(a, b int) int {
			return a * b
		},
	}

	result := 0
	if err := doLines(os.Args[1], func(line string) error {
		root, _ := Parse(line)
		n := root.Eval()
		fmt.Println("Eval", line, "->\n", root, "=", n)
		result += n

		return nil
	}); err != nil {
		return err
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
