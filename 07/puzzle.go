package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var colorStr string = "([a-z]+ [a-z]+)"

var coloredBagsStr string = colorStr + " bags"
var coloredBagsRE *regexp.Regexp = regexp.MustCompile(coloredBagsStr)

var numberedBagStr string = "([0-9]+) " + colorStr + " bag(?:s)?"
var numberedBagRE *regexp.Regexp = regexp.MustCompile(numberedBagStr)

var noOtherBagStr string = "(no other bags)"
var noOtherBagRE *regexp.Regexp = regexp.MustCompile(noOtherBagStr)

var numberedOrNoneStr string = "((?:" + numberedBagStr + ")|(?:" + noOtherBagStr + "))"
var numberedOrNoneRE *regexp.Regexp = regexp.MustCompile(numberedOrNoneStr)

type NumberedBag struct {
	Color string
	Count int
}

type Bag struct {
	Color    string
	Contents map[string]int
}

func (b *Bag) Contains(color string) bool {
	_, ok := b.Contents[color]
	return ok
}

func parseBagColor(color string) (string, error) {
	matches := coloredBagsRE.FindStringSubmatch(color)

	if len(matches) != 2 {
		return "", fmt.Errorf("couldn't parse as colored bag: %s", color)
	}

	return matches[1], nil
}

func parseBagContents(contents string) ([]NumberedBag, error) {
	allMatches := numberedOrNoneRE.FindAllStringSubmatch(contents, -1)

	if len(allMatches) == 0 {
		return nil, fmt.Errorf("couldn't parse as bag contents: %s", contents)
	}

	numberedBags := make([]NumberedBag, 0)

	for i, match := range allMatches {
		if match[0] == "no other bags" {
			if i != 0 || i != len(allMatches)-1 || match[4] != "no other bags" {
				return nil, fmt.Errorf("\"no other bags\" must be the only contents: %s", contents)
			}

			return nil, nil
		}

		num, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}

		color := match[3]

		numberedBags = append(numberedBags, NumberedBag{Color: color, Count: num})
	}

	return numberedBags, nil
}

func NewBag(rule string) (*Bag, error) {
	bag := Bag{
		Contents: make(map[string]int),
	}

	tokens := strings.Split(rule, " contain ")

	color, err := parseBagColor(tokens[0])
	if err != nil {
		return nil, err
	}

	bag.Color = color

	contents, err := parseBagContents(tokens[1])
	if err != nil {
		return nil, err
	}

	for _, b := range contents {
		bag.Contents[b.Color] = b.Count
	}

	return &bag, nil
}

func run() error {
	if len(os.Args) != 2 {
		return fmt.Errorf("Usage: %s INPUT", os.Args[0])
	}

	f, err := os.Open(os.Args[1])
	if err != nil {
		return err
	}
	defer f.Close()

	bags := make(map[string]*Bag, 0)

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		bag, err := NewBag(line)
		if err != nil {
			return err
		}

		bags[bag.Color] = bag
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	fmt.Println("Parsed", len(bags))

	containers := map[string]bool{
		"shiny gold": true,
	}
	l := len(containers)
	for {
		for outerColor, outer := range bags {
			for c, _ := range containers {
				if outer.Contains(c) {
					containers[outerColor] = true
				}
			}
		}
		if len(containers) == l {
			break
		}
		l = len(containers)
	}

	fmt.Println("Number that can contain shiny gold:", len(containers)-1)

	return nil
}

func main() {
	err := run()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
}
