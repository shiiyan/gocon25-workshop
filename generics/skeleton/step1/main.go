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

func Map[T any, U any](s Slice[T], f func(T) U) Slice[U] {
	result := make(Slice[U], len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func main() {
	ints := Slice[int]{1, 2, 3, 4, 5}
	evens := ints.Filter(func(i int) bool { return i%2 == 0 })
	fmt.Printf("%#v\n", evens) // Output: main.Slice{2, 4}

	v := Map(ints, func(v int) string { return strconv.Itoa(v) })
	fmt.Printf("%#v\n", v) // Output: main.Slice{"1", "2", "3", "4", "5"}
}
