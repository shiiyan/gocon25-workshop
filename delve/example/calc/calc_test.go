package main

import "testing"

func TestCalc(t *testing.T) {
	t.Parallel()

	want := 10
	got := calc(2, 3)

	if got != want {
		t.Errorf("got %d, want %d", got, want)
	}
}
