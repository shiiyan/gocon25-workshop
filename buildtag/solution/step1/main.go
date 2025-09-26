package main

import (
	"fmt"
	"strconv"
)

type Slice[T any] []T

func (s Slice[T]) Filter(f func(T) bool) Slice[T] {
	var result Slice[T]
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map[a, b any](s Slice[a], f func(a) b) Slice[b] {
	result := make(Slice[b], len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func main() {
	ints := Slice[int]{1, 2, 3, 4, 5}
	evens := ints.Filter(func(i int) bool { return i%2 == 0 })
	fmt.Printf("%#v\n", evens) // Output: main.Slice[int]{2, 4}

	v := Map(ints, strconv.Itoa)
	fmt.Printf("%#v\n", v) // Output: main.Slice[string]{"1", "2", "3", "4", "5"}
}
