package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

func Intersection(a, b []string) []string {
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

func Remove(needles, haystack []string) []string {
	maxLen := len(haystack)
	result := make([]string, 0, maxLen)
	for _, inA := range haystack {
		remove := false
		for _, inB := range needles {
			if inA == inB {
				remove = true
			}
		}
		if !remove {
			result = append(result, inA)
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
				if a, ok := allergens[cleaned]; ok {
					allergens[cleaned] = Intersection(a, food)
				} else {
					allergens[cleaned] = food
				}
			}
		}

		foods = append(foods, food)

		return nil
	}); err != nil {
		return err
	}

	couldBeAllergen := make(map[string]bool)
	allergenList := []string{}
	for allergen, candidates := range allergens {
		fmt.Println(allergen, len(candidates), candidates)
		allergenList = append(allergenList, allergen)
		for _, i := range candidates {
			couldBeAllergen[i] = true
		}
	}

	// Part 1
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

	// Part 2
	dangerous := map[string]string{}
	for len(allergens) > 0 {
		for allergen, candidates := range allergens {
			if len(candidates) == 1 {
				// Delete allergens as we identify their ingredient
				dangerous[allergen] = candidates[0]
				delete(allergens, allergen)
				for a2, c2 := range allergens {
					dj := Remove(candidates, c2)
					allergens[a2] = dj
				}

				// Break to get a new iteration of the map.
				break
			}
		}
	}

	// Messy and feels redundant, but we need them ordered by allergen name
	sort.Strings(allergenList)
	orderedDangerous := []string{}
	for _, allergen := range allergenList {
		orderedDangerous = append(orderedDangerous, dangerous[allergen])
	}

	fmt.Println(strings.Join(orderedDangerous, ","))

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
