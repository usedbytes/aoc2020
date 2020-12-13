package main

import (
	"testing"
)

func TestColoredBagRE(t *testing.T) {
	rule := "shiny gold bags"
	matches := coloredBagsRE.FindStringSubmatch(rule)

	if len(matches) != 2 || matches[1] != "shiny gold" {
		t.Errorf("%#v", matches)
	}
}

func TestParseBagColor(t *testing.T) {
	rule := "shiny gold bags"

	color, err := parseBagColor(rule)
	if err != nil {
		t.Error(err)
	}

	if color != "shiny gold" {
		t.Errorf("%#v != shiny gold", color)
	}
}

func TestNumberedBagRE(t *testing.T) {
	rule := "3 shiny gold bags"
	matches := numberedBagRE.FindStringSubmatch(rule)

	if len(matches) != 3 || matches[1] != "3" || matches[2] != "shiny gold" {
		t.Errorf("%#v", matches)
	}

	rule = "1 shiny gold bag"
	matches = numberedBagRE.FindStringSubmatch(rule)

	if len(matches) != 3 || matches[1] != "1" || matches[2] != "shiny gold" {
		t.Errorf("%#v", matches)
	}
}

func TestNoOtherBagRE(t *testing.T) {
	rule := "no other bags"
	matches := noOtherBagRE.FindStringSubmatch(rule)

	if len(matches) != 2 || matches[1] != "no other bags" {
		t.Errorf("%#v", matches)
	}
}

func TestNumberedOrNoneRE(t *testing.T) {
	rule := "1 shiny gold bag"
	allMatches := numberedOrNoneRE.FindAllStringSubmatch(rule, -1)
	if len(allMatches) != 1 || allMatches[0][3] != "shiny gold" {
		t.Errorf("%#v", allMatches)
	}

	rule = "no other bags"
	allMatches = numberedOrNoneRE.FindAllStringSubmatch(rule, -1)
	if len(allMatches) != 1 || allMatches[0][4] != "no other bags" {
		t.Errorf("%#v", allMatches)
	}

	bag1 := "1 shiny gold bag"
	bag2 := "2 posh purple bags"
	rule = bag1 + ", " + bag2
	allMatches = numberedOrNoneRE.FindAllStringSubmatch(rule, -1)
	if len(allMatches) != 2 || allMatches[0][3] != "shiny gold" || allMatches[1][3] != "posh purple" {
		t.Errorf("%#v", allMatches)
	}

	bag1 = "1 shiny gold bag"
	bag2 = "2 posh purple bags"
	bag3 := "4 brilliant white bags"
	rule = bag1 + ", " + bag2 + ", " + bag3 + "."
	allMatches = numberedOrNoneRE.FindAllStringSubmatch(rule, -1)
	if len(allMatches) != 3 || allMatches[0][3] != "shiny gold" || allMatches[1][3] != "posh purple" || allMatches[2][3] != "brilliant white" {
		t.Errorf("%#v", allMatches)
	}
}

func TestParseBagContents(t *testing.T) {

	rule := "1 shiny gold bag"
	contents, err := parseBagContents(rule)
	if err != nil {
		t.Error(err)
	} else if len(contents) != 1 {
		t.Errorf("%s: %#v", rule, contents)
	} else if contents[0].Color != "shiny gold" || contents[0].Count != 1 {
		t.Errorf("%s: %#v", rule, contents)
	}

	rule = "1 shiny gold bag, 2 posh purple bags"
	contents, err = parseBagContents(rule)
	if err != nil {
		t.Error(err)
	} else if len(contents) != 2 {
		t.Errorf("%s: %#v", rule, contents)
	} else if contents[0].Color != "shiny gold" || contents[0].Count != 1 {
		t.Errorf("%s: %#v", rule, contents)
	} else if contents[1].Color != "posh purple" || contents[1].Count != 2 {
		t.Errorf("%s: %#v", rule, contents)
	}

	rule = "1 shiny gold bag, no other bags"
	_, err = parseBagContents(rule)
	if err == nil {
		t.Errorf("%v should have failed", rule)
	}

	rule = "no other bags, 1 shiny gold bag"
	_, err = parseBagContents(rule)
	if err == nil {
		t.Errorf("%v should have failed", rule)
	}

	rule = "no other bags"
	contents, err = parseBagContents(rule)
	if err != nil {
		t.Error(err)
	} else if len(contents) != 0 {
		t.Errorf("%s: %#v", rule, contents)
	}
}

func TestParseBag(t *testing.T) {
	rule := "light red bags contain 1 bright white bag, 2 muted yellow bags."
	bag, err := NewBag(rule)
	if err != nil {
		t.Error(err)
	} else if bag.Color != "light red" {
		t.Errorf("%v: %#v", rule, *bag)
	} else if bag.Contents["bright white"] != 1 {
		t.Errorf("%v: %#v", rule, *bag)
	} else if bag.Contents["muted yellow"] != 2 {
		t.Errorf("%v: %#v", rule, *bag)
	}

	rule = "faded blue bags contain no other bags."
	bag, err = NewBag(rule)
	if err != nil {
		t.Error(err)
	} else if bag.Color != "faded blue" {
		t.Errorf("%v: %#v", rule, *bag)
	} else if len(bag.Contents) != 0 {
		t.Errorf("%v: %#v", rule, *bag)
	}
}

func TestContains(t *testing.T) {
	bag := &Bag{
		Color: "bright white",
		Contents: map[string]int {
			"muted yellow": 1,
		},
	}

	color := "light red"
	res := bag.Contains(color)
	if res {
		t.Errorf("%#v contains %s: %v", bag, color, res)
	}

	bag.Contents["light red"] = 7
	res = bag.Contains(color)
	if !res {
		t.Errorf("%#v contains %s: %v", bag, color, res)
	}
}
