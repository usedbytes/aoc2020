package main

import (
	"testing"
)

func segmentTest(t *testing.T, segments string, setSize, expected int) {
	r, err := BinarySegment(segments, setSize)
	if err != nil {
		t.Fatal(err)
	} else if r != expected {
		t.Errorf("expected %d got %d", expected, r)
	}
}

func TestBinarySegment(t *testing.T) {
	segmentTest(t, "F", 2, 0)
	segmentTest(t, "B", 2, 1)
	segmentTest(t, "FB", 4, 1)
	segmentTest(t, "BFFFBBF", 128, 70)
	segmentTest(t, "RRR", 8, 7)
	segmentTest(t, "RLL", 8, 4)
}
