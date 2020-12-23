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
	Validate([]byte, map[int]Rule) (bool, []int)
}

type DerivedRule struct {
	SubRules [][]int

	Memo map[string][]int
}

func in(haystack []int, needle int) bool {
	for _, n := range haystack {
		if n == needle {
			return true
		}
	}
	return false
}

func (dr *DerivedRule) Validate(data []byte, rules map[int]Rule) (bool, []int) {

	// Check if we've been here before
	if res, memod := dr.Memo[string(data)]; memod {
		return len(res) > 0, res
	}

	var endPositions []int
	for _, option := range dr.SubRules {
		var (
			ruleIdx int
			rule    int

			// Each option starts at the beginning
			startPositions []int = []int{0}

			optionIsBad bool = false
		)

		for ruleIdx, rule = range option {
			var (
				thisEndPositions []int
			)

			for _, start := range startPositions {
				ok, tmp := rules[rule].Validate(data[start:], rules)
				if !ok {
					// No successful matches for this starting position
					continue
				} else if len(tmp) == 0 {
					panic("OK but no results")
				}

				for _, consumed := range tmp {
					if consumed == 0 {
						panic("zero consumed")
					}

					end := start + consumed

					// This would use up the whole data, but we have more rules to match
					// in this option, so that's not OK
					if end >= len(data) && ruleIdx != len(option)-1 {
						continue
					}

					// Don't add duplicates
					if in(thisEndPositions, end) {
						continue
					}

					thisEndPositions = append(thisEndPositions, end)
				}
			}

			// thisEndPositions now contains all the ways to get "past" 'rule'
			if len(thisEndPositions) == 0 {
				// None of the possible starting points for this rule produced
				// anything viable so this option can't go anywhere
				optionIsBad = true
				break
			}

			// The next rule starts at the end of this one
			startPositions = thisEndPositions
		}
		if !optionIsBad {
			// We managed to traverse the whole option, so we can
			// reach everything in startPositions
			endPositions = append(endPositions, startPositions...)
		}
	}

	if _, ok := dr.Memo[string(data)]; !ok {
		dr.Memo[string(data)] = endPositions
	}

	return len(endPositions) > 0, endPositions
}

type BaseRule struct {
	Char byte
}

func (br *BaseRule) Validate(data []byte, rules map[int]Rule) (bool, []int) {
	if len(data) == 0 {
		return false, nil
	} else if data[0] != br.Char {
		return false, nil
	}

	return true, []int{1}
}

func Validate(data []byte, against int, rules map[int]Rule) bool {
	ok, consumed := rules[against].Validate(data, rules)
	if !ok {
		return false
	}

	for _, c := range consumed {
		if c == len(data) {
			return true
		}
	}

	return false
}

func run() error {
	part2 := false
	if len(os.Args) > 2 {
		part2 = true
	}

	rules := make(map[int]Rule)

	parsing := true
	res := 0

	if err := doLines(os.Args[1], func(line string) error {
		if len(line) == 0 {
			// Done parsting rules
			parsing = false

			if part2 {
				fmt.Println("Overriding rules (Part 2)")
				// Override the two rules
				// 8: 42 | 42 8
				// 11: 42 31 | 42 11 31
				rules[8] = &DerivedRule{
					SubRules: [][]int{
						{42},
						{42, 8},
					},
					Memo: make(map[string][]int),
				}
				rules[11] = &DerivedRule{
					SubRules: [][]int{
						{42, 31},
						{42, 11, 31},
					},
					Memo: make(map[string][]int),
				}
			}
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
					Memo:     make(map[string][]int),
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
			// Validate message
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
