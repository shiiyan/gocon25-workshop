package main

import (
	"fmt"
)

// Numeric is a type constraint that includes all numeric types
// ~int means any type with underlying type int (e.g., type MyInt int)
type Numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Custom types with underlying numeric types
type UserID int
type Score float64
type Counter uint

// Sum calculates the sum of a slice of numeric values
// The ~ in the constraint allows custom types with numeric underlying types
func Sum[T Numeric](values []T) T {
	var result T
	for _, v := range values {
		result += v
	}
	return result
}

// Max returns the maximum value from a slice
// TODO: Workshop Exercise 1 - Complete this function
// Hint: Use the Numeric constraint and compare values
func Max[T Numeric](values []T) T {
	if len(values) == 0 {
		var zero T
		return zero
	}
	
	// TODO: Complete the implementation
	// Start with max := values[0]
	// Then iterate and compare
	
	max := values[0]
	for _, v := range values[1:] {
		if v > max {
			max = v
		}
	}
	return max
}

// Scale multiplies all elements by a factor
func Scale[T Numeric](values []T, factor T) []T {
	result := make([]T, len(values))
	for i, v := range values {
		result[i] = v * factor
	}
	return result
}

// TODO: Workshop Exercise 2 - Create a generic MinMax function
// that returns both min and max values from a slice
// func MinMax[T Numeric](values []T) (T, T) {
//     // Your implementation here
// }

func main() {
	// Example 1: Using Sum with different numeric types
	fmt.Println("=== Example 1: Sum with different types ===")
	
	// Regular int slice
	ints := []int{1, 2, 3, 4, 5}
	fmt.Printf("Sum of ints %v: %d\n", ints, Sum(ints))
	
	// Custom type with underlying int
	userIDs := []UserID{101, 102, 103}
	fmt.Printf("Sum of UserIDs %v: %d\n", userIDs, Sum(userIDs))
	
	// Custom type with underlying float64
	scores := []Score{85.5, 90.0, 78.5}
	fmt.Printf("Sum of Scores %v: %.1f\n", scores, Sum(scores))
	
	// Example 2: Using Max
	fmt.Println("\n=== Example 2: Max values ===")
	fmt.Printf("Max of ints: %d\n", Max(ints))
	fmt.Printf("Max of scores: %.1f\n", Max(scores))
	
	// Example 3: Using Scale
	fmt.Println("\n=== Example 3: Scaling values ===")
	counters := []Counter{10, 20, 30}
	scaled := Scale(counters, 2)
	fmt.Printf("Original counters: %v\n", counters)
	fmt.Printf("Scaled by 2: %v\n", scaled)
	
	// TODO: Workshop Exercise 3
	// Test your MinMax function here
	// fmt.Println("\n=== Exercise: MinMax ===")
	// min, max := MinMax(ints)
	// fmt.Printf("Min: %d, Max: %d\n", min, max)
}