package main

import (
	"testing"
)

func TestXMAS(t *testing.T) {
	x := NewXMAS(5)

	seq := []int{35, 20, 15, 25, 47, 40, 62, 55, 65, 95, 102, 117, 150, 182, 127, 219, 299, 277, 309, 576}

	for i, v := range seq {
		if i >= 5 {
			valid := x.Valid(v)
			if !valid && v != 127 {
				t.Errorf("%d reported as invalid, should be valid", v)
			} else if valid && v == 127 {
				t.Errorf("%d reported as valid, should be invalid", v)
			}
		}

		x.Receive(v)
	}
}
