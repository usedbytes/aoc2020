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

type OpFunc func(a, b int) int

func OpAdd(a, b int) int {
	return a + b
}

func OpMul(a, b int) int {
	return a * b
}

func Eval(s string) (int, int) {
	var (
		op     OpFunc
		args   []int
		skipTo int
	)

	for i, r := range s {
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

			args = append(args, n)
		} else if strings.ContainsRune("+*", r) {
			if r == '+' {
				op = OpAdd
			} else if r == '*' {
				op = OpMul
			} else {
				panic(r)
			}
		} else if strings.ContainsRune("(", r) {
			n, di := Eval(s[i+1:])
			skipTo = i + di + 1
			args = append(args, n)
		} else if strings.ContainsRune(")", r) {
			return args[0], i + 1
		}

		if len(args) == 2 {
			args[0] = op(args[0], args[1])
			args = args[:1]
		}
	}
	return args[0], len(s)
}

func run() error {

	result := 0
	if err := doLines(os.Args[1], func(line string) error {
		n, _ := Eval(line)
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
