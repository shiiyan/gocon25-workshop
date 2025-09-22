package main

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	t.Parallel()

	want := time.Date(2025, 9, 27, 0, 0, 0, 0, time.UTC)
	got := now()

	if !got.Equal(want) {
		t.Errorf("now() = %v; want %v", got, want)
	}
}
