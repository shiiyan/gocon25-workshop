package main

import "fmt"

type Container[T fmt.Stringer] struct {  // TODO: any を適切な制約に変更
	items []T
}

// Add メソッド（完成済み）
func (c *Container[T]) Add(item T) {
	c.items = append(c.items, item)
}

// PrintAll メソッド（完成済み）
func (c *Container[T]) PrintAll() {
	for _, item := range c.items {
		fmt.Println(item.String())
	}
}

// Person型
type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s (%d years)", p.Name, p.Age)
}

// Product型
type Product struct {
	Name  string
	Price float64
}

func (p Product) String() string {
	return fmt.Sprintf("%s: $%.2f", p.Name, p.Price)
}

func main() {
	// Person用のContainer
	people := Container[Person]{}
	people.Add(Person{"Alice", 30})
	people.Add(Person{"Bob", 25})
	fmt.Println("People:")
	people.PrintAll()

	// Product用のContainer
	products := Container[Product]{}
	products.Add(Product{"Laptop", 999.99})
	products.Add(Product{"Mouse", 25.50})
	fmt.Println("\nProducts:")
	products.PrintAll()
}

// 期待される出力:
// People:
// Alice (30 years)
// Bob (25 years)
//
// Products:
// Laptop: $999.99
// Mouse: $25.50
