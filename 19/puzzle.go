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
			fmt.Println("validate", string(data[consumed:]), "against", element)
			ok, n = rules[element].Validate(data[consumed:], rules)
			fmt.Println("ok", ok, "n", n)
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
	fmt.Println("char", string(br.Char))
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

	part2 := false
	if len(os.Args) > 2 {
		part2 = true
	}

	parsing := true
	res := 0

	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			parsing = false

			if part2 {
				// Override the two rules
				// 8: 42 | 42 8
				// 11: 42 31 | 42 11 31
				rules[8] = &DerivedRule{
					SubRules: [][]int{
						{ 42 },
						{ 42, 8 },
					},
				}
				rules[11] = &DerivedRule{
					SubRules: [][]int{
						{ 42, 31 },
						{ 42, 11, 31 },
					},
				}
			}

			return nil
		}

		if parsing {
			ruleNumberThenRules := strings.Split(line, ": ")
			fmt.Println(ruleNumberThenRules)
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
			/*
			ok := Validate([]byte(line), 0, rules)
			if ok {
				res++
			}
			*/
		}

		return nil
	}); err != nil {
		return err
	}

	fmt.Println(res)

	tests := []string{
		//"bbabbbbaabaabba",
		"babbbbaabbbbbabbbbbbaabaaabaaa",
		//"aaabbbbbbaaaabaababaabababbabaaabbababababaaa",
		//"bbbbbbbaaaabbbbaaabbabaaa",
		//"bbbababbbbaaaaaaaabbababaaababaabab",
		//"ababaaaaaabaaab",
		//"ababaaaaabbbaba",
		//"baabbaaaabbaaaababbaababb",
		//"abbbbabbbbaaaababbbbbbaaaababb",
		//"aaaaabbaabaaaaababaa",
		//"aaaabbaabbaaaaaaabbbabbbaaabbaabaaa",
		//"aabbbbbaabbbaaaaaabbbbbababaaaaabbaaabba",
	}

	for _, t := range tests {
		fmt.Println(t, "->", Validate([]byte(t), 0, rules))
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
