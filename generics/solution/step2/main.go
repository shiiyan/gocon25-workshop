package main

import "fmt"

// Generic container that works with any type implementing fmt.Stringer
type Container[T fmt.Stringer] struct {
	items []T
}

// Add item to container
func (c *Container[T]) Add(item T) {
	c.items = append(c.items, item)
}

// Print all items
func (c *Container[T]) PrintAll() {
	for _, item := range c.items {
		fmt.Println(item.String())
	}
}

// Example types implementing Stringer
type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

type Product struct {
	Name  string
	Price float64
}

func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

func main() {
	// Container for Person type
	people := Container[Person]{}
	people.Add(Person{"Alice", 30})
	people.Add(Person{"Bob", 25})
	fmt.Println("People:")
	people.PrintAll()
	// Output:
	// People:
	// Alice (30 years)
	// Bob (25 years)
	//

	// Container for Product type
	products := Container[Product]{}
	products.Add(Product{"Laptop", 999.99})
	products.Add(Product{"Mouse", 25.50})
	fmt.Println("\nProducts:")
	products.PrintAll()
	// Output:
	// Products:
	// Laptop: $999.99
	// Mouse: $25.50
}

