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

	fmt.Println(p2) // {Alice ALICE_DNA}
	fmt.Println(f2) // {Icecream Sweet}
}

// `any` を使うと、任意の型を引数で受け取れるようになります
func copyStruct(src, dst any) {
	// `reflect.ValueOf()` と `Value.Elem()` を使うと、ポインター型の値情報から値そのものを取り出せます
	// 引数の struct のフィールドを反復して処理します
	// ループ中の処理(1): `Value.NumField()` を使うと、フィールドの数を取得できます
	// ループ中の処理(2): `Value.Filed(index)` を使うと、その `index` のフィールドを取得できます
	// ループ中の処理(3): `Value.CanSet()` を使うと、フィールドが代入可能であるかをチェックできます
	// ループ中の処理(4): `Value.Set(Value)` を使うと、値情報を別の値情報にマッピングできます

	srcValue := reflect.ValueOf(src).Elem()
	dstValue := reflect.ValueOf(dst).Elem()
	for i := 0; i < srcValue.NumField(); i++ {
		if dstValue.Field(i).CanSet() && srcValue.Type().Field(i).Tag.Get("copyable") != "false" {
			dstValue.Field(i).Set(srcValue.Field(i))
		}
	}
}
