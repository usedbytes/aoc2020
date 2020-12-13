package main

import (
	"testing"
)

func doForm(answers string) *Form {
	var f Form
	for _, a := range answers {
		f.Answer(a)
	}

	return &f
}

func doFormTest(t *testing.T, answers string, expected int) {
	f := doForm(answers)
	sum := f.Sum()
	if sum != expected {
		t.Errorf("%v: expected %d got %d", answers, expected, sum)
	}
}

func TestForm(t *testing.T) {
	doFormTest(t, "a", 1)
	doFormTest(t, "abc", 3)
	doFormTest(t, "az", 2)
}

func doCombineTest(t *testing.T, answers []string, op CombineFunc, expected int) {
	var group Form
	for _, a := range answers[0] {
		group.Answer(a)
	}

	for _, ind := range answers[1:] {
		var individual Form
		for _, a := range ind {
			individual.Answer(a)
		}
		group.Combine(individual, op)
	}

	sum := group.Sum()
	if sum != expected {
		t.Errorf("%v: expected %d got %d", answers, expected, sum)
	}
}

func TestCobine(t *testing.T) {
	doCombineTest(t, []string{"abc"}, Or, 3)
	doCombineTest(t, []string{"a", "b", "c"}, Or, 3)
	doCombineTest(t, []string{"ab", "ac"}, Or, 3)
	doCombineTest(t, []string{"a", "a", "a", "a"}, Or, 1)
	doCombineTest(t, []string{"b"}, Or, 1)

	doCombineTest(t, []string{"abc"}, And, 3)
	doCombineTest(t, []string{"a", "b", "c"}, And, 0)
	doCombineTest(t, []string{"ab", "ac"}, And, 1)
	doCombineTest(t, []string{"a", "a", "a", "a"}, And, 1)
	doCombineTest(t, []string{"b"}, And, 1)
}
