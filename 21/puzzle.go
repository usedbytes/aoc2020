package main

import (
	"bufio"
	"fmt"
	"os"
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

func Intersect(a, b []string) []string {
	if len(a) == 0 {
		return b
	} else if len(b) == 0 {
		return a
	}

	maxLen := len(a)
	if len(b) > len(a) {
		maxLen = len(b)
	}
	result := make([]string, 0, maxLen)
	for _, inA := range a {
		for _, inB := range b {
			if inA == inB {
				result = append(result, inA)
			}
		}
	}
	return result
}

func run() error {
	foods := [][]string{}
	ingredients := map[string]bool{}
	allergens := map[string][]string{}

	if err := doLines(os.Args[1], func(line string) error {
		// Split at spaces, then just clean each token afterwards
		toks := strings.Split(line, " ")
		isIngredient := true
		food := []string{}
		for _, t := range toks {
			cleaned := strings.Trim(t, " ,()")
			if cleaned == "contains" {
				isIngredient = false
				continue
			}

			if isIngredient {
				ingredients[cleaned] = true
				food = append(food, cleaned)
			} else {
				// The allergens list never has false positives, so we can intersect
				// all the corresponding ingredients to get the possible contributors
				allergens[cleaned] = Intersect(allergens[cleaned], food)
			}
		}

		foods = append(foods, food)

		return nil
	}); err != nil {
		return err
	}

	couldBeAllergen := make(map[string]bool)
	for _, a := range allergens {
		for _, i := range a {
			couldBeAllergen[i] = true
		}
	}

	count := 0
	for i, _ := range ingredients {
		if _, ok := couldBeAllergen[i]; !ok {
			// 'i' is not an allergen, count how often it appears
			for _, f := range foods {
				for _, j := range f {
					if j == i {
						count++
					}
				}
			}
		}
	}
	fmt.Println(count)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
