package main

import (
	"fmt"
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

// `any` を使うと、任意の型を引数で受け取れるようになります
func copyStruct() {
	// `reflect.ValueOf()` と `Value.Elem()` を使うと、ポインター型の値情報から値そのものを取り出せます
	// 引数の struct のフィールドを反復して処理します
	// ループ中の処理(1): `Value.NumField()` を使うと、フィールドの数を取得できます
	// ループ中の処理(2): `Value.Filed(index)` を使うと、その `index` のフィールドを取得できます
	// ループ中の処理(3): `Value.CanSet()` を使うと、フィールドが代入可能であるかをチェックできます
	// ループ中の処理(4): `Value.Set(Value)` を使うと、値情報を別の値情報にマッピングできます
}
