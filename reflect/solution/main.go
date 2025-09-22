package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	DNA  string
	Soul string `copyable:"false"`
}

type Food struct {
	Name             string
	Kind             string
	secretIngredient string
}

func main() {
	p1 := Person{Name: "Alice", DNA: "ALICE_DNA", Soul: "ALICE_SOUL"}
	p2 := Person{}
	f1 := Food{Name: "Icecream", Kind: "Sweet", secretIngredient: "Salt"}
	f2 := Food{}

	copyStruct(&p1, &p2)
	copyStruct(&f1, &f2)

	fmt.Println(p2) // {Alice 20}
	fmt.Println(f2) // {Icecream Sweet}
}

func copyStruct(dst, src interface{}) {
	dv := reflect.ValueOf(dst).Elem()
	sv := reflect.ValueOf(src).Elem()

	for i := 0; i < dv.NumField(); i++ {
		if dv.Field(i).CanSet() {
			dv.Field(i).Set(sv.Field(i))
		}
	}
}
