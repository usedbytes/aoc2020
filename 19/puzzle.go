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

type Rule interface {
	Validate([]byte, map[int]Rule) (bool, int)
}

type DerivedRule struct {
	SubRules [][]int
}

func (dr *DerivedRule) Validate(data []byte, rules map[int]Rule) (bool, int) {
	var (
		ok bool
		consumed int
	)
	for _, option := range dr.SubRules {
		ok = true
		consumed = 0
		for _, element := range option {
			var n int
			ok, n = rules[element].Validate(data[consumed:], rules)
			if !ok {
				break
			}
			consumed += n
		}
		if ok {
			break
		}
	}

	return ok, consumed
}

type BaseRule struct {
	Char byte
}

func (br *BaseRule) Validate(data []byte, rules map[int]Rule) (bool, int) {
	if len(data) == 0 {
		return false, 0
	}
	return data[0] == br.Char, 1
}

func Validate(data []byte, against int, rules map[int]Rule) bool {
	ok, consumed := rules[against].Validate(data, rules)
	if ok && consumed == len(data) {
		return true
	}

	return false
}

func run() error {

	rules := make(map[int]Rule)

	parsing := true
	res := 0

	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			parsing = false
			return nil
		}

		if parsing {
			ruleNumberThenRules := strings.Split(line, ": ")

			ruleNo, err := strconv.Atoi(ruleNumberThenRules[0])
			if err != nil {
				return err
			}

			var rule Rule
			options := strings.Split(ruleNumberThenRules[1], " | ")
			if len(options) == 1 && options[0][0] == '"' {
				rule = &BaseRule{
					Char: options[0][1],
				}
			} else {
				dr := &DerivedRule{
					SubRules: make([][]int, len(options)),
				}
				rule = dr
				for o, option := range options {
					toks := strings.Split(option, " ")
					subRules := make([]int, len(toks))

					for i := range subRules {
						n, err := strconv.Atoi(toks[i])
						if err != nil {
							return err
						}
						subRules[i] = n
					}
					dr.SubRules[o] = subRules
				}
			}
			rules[ruleNo] = rule
		} else {
			ok := Validate([]byte(line), 0, rules)
			if ok {
				res++
			}
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(res)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
