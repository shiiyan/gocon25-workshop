package main

import (
	"fmt"
	"strconv"
)

type Slice []any

func (s Slice) Filter(f func(any) bool) Slice {
	var result Slice
	for _, v := range s {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

func Map(s Slice, f func(any) any) Slice {
	result := make(Slice, len(s))
	for i, v := range s {
		result[i] = f(v)
	}
	return result
}

func main() {
	ints := Slice{1, 2, 3, 4, 5}
	evens := ints.Filter(func(i any) bool { return (i.(int))%2 == 0 })
	fmt.Printf("%#v\n", evens) // Output: main.Slice{2, 4}

	v := Map(ints, func(v any) any { return strconv.Itoa(v.(int)) })
	fmt.Printf("%#v\n", v) // Output: main.Slice{"1", "2", "3", "4", "5"}
}
