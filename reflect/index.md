author: Nozomu Ikuta
summary: reflect パッケージによるメタプログラミング
id: reflect
categories: codelab,markdown
environments: Web
status: Published

# `reflect` パッケージによるメタプログラミング Codelab

## この Codelab について

この Codelab では、Goの標準ライブラリのひとつである `reflect` パッケージを学びます。

以下のゴールを達成して、Goでメタプログラミングができるようになりましょう！

- reflect パッケージとは何かを理解する
- reflect の基本的な使い方を理解する
- reflect を使って、任意の型の struct をコピーする

---

## reflect パッケージとは何かを理解する

[reflect](https://pkg.go.dev/reflect) パッケージは、Goの標準ライブラリのひとつです。
reflect パッケージを使うと、プログラムは、自身を鏡に反射（reflection）するように、ランタイムの型や値を見ることができるようになります。

そして、型や値を見るだけだなく、ランタイムでそれらを変更することができます。

## 基本的な使い方: 型の情報を取得する

ランタイムの型の情報を取得するには [`reflect.TypeOf()`](https://pkg.go.dev/reflect#TypeOf)を使用します。

取得できる型情報は [`Type` interface](https://pkg.go.dev/reflect#Type) を実装しています。

```go
package main

import (
	"fmt"
	"reflect"
)

type person struct {
	name string
}

func main() {
	p := person{
		name: "にゅも太郎",
	}

	t := reflect.TypeOf(p)

	fmt.Println(t.Name()) // person
}
```

## 基本的な使い方: 値の情報を取得する

ランタイムの値の情報を取得するには [`reflect.ValueOf()`](https://pkg.go.dev/reflect#ValueOf)を使用します。

取得できる値の情報は [`Value` struct](https://pkg.go.dev/reflect#Value) を実装しています。

```go
package main

import (
	"fmt"
	"reflect"
)

type person struct {
	name string
}

func main() {
	p := &person{
		name: "にゅも太郎",
	}

	v := reflect.ValueOf(p)

	fmt.Println(v.Elem()) // {にゅも太郎}
}
```

## 問題: 値のコピー（１）

ある struct A を、別の struct Bにコピーすることを考えます。
もっともシンプルな方法は、すべてのフィールドに対して代入をおこなうことです。

```go
package main

import (
	"fmt"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	p1 := Person{Name: "Alice", Age: 20}
	p2 := Person{}

	copyPerson(&p1, &p2)

	fmt.Println(p2) // {Alice 20}
}

func copyPerson(src, dst *Person) {
	dst.Name = src.Name
	dst.Age = src.Age
}
```

この方法には、以下のような問題があります。

- `Person` のフィールドが増えると `copyPerson` のコードを変更する必要がある
- 型ごとにコピー関数を用意する必要がある

reflect パッケージを使ってこの問題を解決してみましょう。

## 問題: 値のコピー（２）

任意の型の struct をコピーできる `copyStruct` 関数を実装してみましょう。
雛形のコードはこちらです。

```go
package main

import (
	"fmt"
	"uuid"
)

type Person struct {
	Name string
	DNA string
	Soul uuid `copyable:"false"`
}

type Food struct {
	Name string
	Kind string
	secretIngredient string
}

func main() {
	p1 := Person{Name: "Alice", DNA: "ALICE_DNA", Soul: uuid.New()}
	p2 := Person{}
	f1 := Food{Name: "Icecream", Kind: "Sweet", secretIngredient: "Salt"}
	f2 := Food{}

	copyStruct(&p1, &p2)
	copyStruct(&f1, &f2)

	fmt.Println(p2) // {Alice 20}
	fmt.Println(f2) // {Icecream Sweet}
}

// TODO: implement
func copyStruct() {}
```

### 考え方

- 任意の型を引数で受け取れるようにします
    - `interface{}` を使いましょう
- 引数の値情報から、値そのものを取り出します
    - `reflect.ValueOf()` を使いましょう
    - ポインター型の値情報から値を取り出すには `Value.Elem()` を使います
- 引数の struct のフィールドを反復して処理します
    - `Value.NumField()` を使うと、フィールドの数を取得できます
    - `Value.Filed(index)` を使うと、その `index` のフィールドを取得できます
- フィールドが代入可能であるかをチェックします
    - `Value.CanSet()` を使いましょう
- フィールドに値を代入します
    - `Value.Set(Value)` を使うと、値情報を別の値情報にマッピングできます

### 応用問題

フィールドが `copyable:"false"` の Struct Tag をもつ場合、代入をスキップしてみましょう。

## その他の学習

reflect パッケージを使ったメタプログラミングには、ランタイムの振る舞いを変える以外にも、以下のように使うことができます。
関連のOSSのソースコードなどを読んでみましょう。

- データベースの情報を読み込んで struct にマッピングする（ORM）
- Struct Tag と `go generate` を使って、コード生成する

## 参考資料

- [`reflect` パッケージドキュメント](https://pkg.go.dev/reflect)
