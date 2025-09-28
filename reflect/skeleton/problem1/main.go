package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	p1 := Person{Name: "Alice", Age: 20}
	p2 := Person{}

	copyPersonByReflect(&p1, &p2)

	p1.Age = 30

	fmt.Println(p2) // {Alice 20}
}

// func copyPerson(src, dst *Person) {
// 	dst.Name = src.Name
// 	dst.Age = src.Age
// }

func copyPersonByReflect(src, dst *Person) {
	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()
	dstValue.Set(srcValue)
}
