author: kamata
summary: Introduction to generics
id: generics
categories: codelab,markdown
environments: Web
status: Published

# Go  generics  Codelab

このコードラボでは、generics について、実践的な例を通じて学習します。

## 学習目標

- 型パラメータを使った型安全なコードの実装方法を理解する
- Interface を型制約として活用する方法を学ぶ
- 複雑な型制約のパターンを習得する

## 進め方

1. 各ステップの「学習内容」を読み、概念を理解します
2. `skeleton/stepX/main.go` を編集して実装します
3. 仕上げに `solution/stepX/main.go` と比較して理解を深めます

### 前提条件

- Go 1.18 以降（Go 1.25以降を推奨）
- 基本的な Go の文法の理解

### 実行して確認

各ステップのコードを修正したら、以下のコマンドで動作を確認してください。

```bash
go run skeleton/stepX/main.go
```

---

## Step 1:  generics の基礎 - 型パラメータによる抽象化

### このステップで学ぶこと

**型パラメータ** という新しい概念を理解し、同一のロジックを複数の型で再利用する方法を学びます。

### なぜ generics が必要か？

Go 1.18以前は、複数の型に対して同じ処理を書く場合、以下のような選択肢しかありませんでした。

1. **型ごとに関数を複製する**（コードの重複）
2. **`interface{}`（現在の`any`）を使う**（型安全性の喪失）

現在の `skeleton/step1/main.go` は `any` を使った実装例です。この方法には以下の問題があります。

```go
// 問題1: 型情報が失われる
func (s Slice) Filter(f func(any) bool) Slice  // any型として扱う

// 問題2: 使う側で型アサーションが必要
evens := ints.Filter(func(i any) bool {
    return (i.(int))%2 == 0  // .(int) で型アサーション
})

// 問題3: 実行時にパニックのリスク
// もし間違った型でアクセスしたら実行時エラー
```

### 型パラメータの概念

型パラメータは「型を後から決める」仕組みです。

```go
// [T any] が型パラメータ
// T は「何かの型」を表すプレースホルダー
type Slice[T any] []T

// 使うときに具体的な型を指定
var ints Slice[int]       // T = int
var strings Slice[string] // T = string
```

**重要な概念**

- `T` は型変数（Type Variable）と呼ばれる
- `any` は型制約（Type Constraint）
    - この場合「どんな型でもOK」
- 型パラメータは関数にも使える：`func Map[A, B any](...)`

### 実装タスク

このステップの内容を踏まえて、`skeleton/step1/main.go` を generics を用いて修正してください。

### 理解度チェック

- 型パラメータが「型のプレースホルダー」であることを理解した
- 型アサーションが不要になる理由を説明できる
- コンパイル時型チェックの利点を理解した
- 複数の型パラメータ（Map関数のA, B）の使い方を理解した

---

## Step 2: Interface を型制約として使う

### このステップで学ぶこと

Interface を型制約として使用し、型パラメータに「条件」を付ける方法を学びます。

### 型制約の必要性

Step 1 では `any` を使いましたが、これは「どんな型でもOK」という意味です。しかし、実際のコードでは「特定のメソッドを持つ型」という条件を付けたい場合があります。

```go
// 現在の問題：T は any なので String() メソッドが保証されない
type Container[T any] struct {
    items []T
}

func (c *Container[T]) PrintAll() {
    for _, item := range c.items {
        fmt.Println(item.String())  // コンパイルエラー！
    }
}
```

### Interface による型制約

Interface を型制約として使うことで、「この条件を満たす型のみ」を指定できます。

```go
// T は fmt.Stringer を実装した型のみ
type Container[T fmt.Stringer] struct {
    items []T
}
```

### 通常の Interface を使う場合との違い

**重要な違い：型の統一性**

```go
// 従来の方法：異なる型を混在できる
func PrintAll(items []fmt.Stringer) {
    // items = []fmt.Stringer{Person{}, Product{}} // 異なる型OK
}

//  generics の型制約：同一の型で統一
type Container[T fmt.Stringer] struct {
    items []T  // すべて同じ具体的な型T
}
// Container[Person] と Container[Product] は別の型
```

**型制約のメリット**

1. **型の一貫性** — コンテナ内の全要素が同じ型
2. **パフォーマンス** — interface のボックス化が不要な場合がある
3. **型情報の保持** — 元の型の情報が失われない

### なぜ型制約が有用か？

```go
// Interface を型制約として使うと...
container := Container[Person]{}
container.Add(Person{"Alice", 30})  // OK
container.Add(Product{"Book", 10})  // コンパイルエラー（型が違う）

// 通常の interface 引数
items := []fmt.Stringer{}
items = append(items, Person{"Alice", 30})  // OK
items = append(items, Product{"Book", 10})  // OK（混在可能）
```

### 実装タスク

このステップの内容を踏まえて、`skeleton/step2/main.go` を修正してください。

### 理解度チェック

- 型制約により特定のメソッドの存在を保証できることを理解した
- 型制約と通常の interface 引数の違いを説明できる
- 型の一貫性がなぜ重要かを理解した

---

## Step 3: 複雑な型制約 - Interface の合成
### このステップで学ぶこと
複数の制約を組み合わせた高度な型制約パターンを理解し、ポインタ専用の制約を正しく表現できるようになります。

### Pointer 制約を正しく書く難しさ

JSON のアンマーシャル処理を汎用化したい場合、値そのもの (`T`) を返しつつ `*T` にだけ定義されたメソッドを呼び出さなければなりません。`json.Unmarshal` はポインタを受け取るため、型制約で「`*T` かつ `json.Unmarshaler`」を厳密に指定する必要があります。

### 解決したい問題

```go
var user User
json.Unmarshal(data, &user)  // ポインタが必須

user, err := Unmarshal[User](data) //  generics で値型として受け取りたい
```

### 複合型制約（Type Set Intersection）

以下のように複数の型制約を同時に満たす型を指定します。

```go
type Unmershaller[T any] interface {
    *T
    json.Unmarshaler
}
```

これは「`*T` **かつ** `json.Unmarshaler` を実装した型」を意味します。

### 2つの型パラメータの連携

```go
func Unmarshal[T any, PT Unmershaller[T]](data []byte) (T, error) {
    var v T
    err := PT(&v).UnmarshalJSON(data)
    return v, err
}
```

- `T` : 呼び出し側へ返したい値型（例：`User`）
- `PT`: `*T` と互換性があり、`json.Unmarshaler` を実装した型

### 型推論の進化

Go 1.20 以降では `PT` を省略でき、`Unmarshal[User](data)` のようにシンプルに呼び出せます。コンパイラが `Unmershaller[User]` を満たす型として `*User` を推論します。

### 現在の skeleton の問題

`skeleton/step3/main.go` の `Unmershaller` は `json.Unmarshaler` しか制約として指定していません。そのため、コンパイラは `PT` と `*T` の関係を理解できず、以下のエラーが発生します：

```
cannot convert &v (value of type *T) to type PT
```

### 実装タスク

このステップの内容を踏まえて、`skeleton/step3/main.go` を修正してください。

### 理解度チェック

- 複合型制約（`*T` かつ `json.Unmarshaler`）の意味を説明できる
- なぜ 2 つの型パラメータ（`T` と `PT`）が必要か理解している
- `PT(&v)` のキャストが必要な理由を説明できる
- 値型だけを許すとコンパイル エラーになることを把握している

---

## まとめ

### 学んだ概念

**Step 1: 型パラメータの基礎**

- 型パラメータは「型のプレースホルダー」
- コンパイル時の型安全性を保ちながら汎用的なコードが書ける
- 型推論により、多くの場合型指定を省略できる

**Step 2: Interface による型制約**

- Interface を型制約として使い、型パラメータに条件を付けられる
- 通常の interface 引数とは異なり、型の一貫性が保証される
- パフォーマンスと型安全性の両立が可能

**Step 3: 複雑な型制約**

- 複数の制約を組み合わせた複合型制約が定義できる
- 複数の型パラメータを連携させて高度な抽象化が可能
- 型推論の改善により、使いやすいAPIを提供できる

### 次のステップ

1. 標準ライブラリの generics 活用例を調べる（[`iter`]( https://pkg.go.dev/iter ) パッケージなど）
2. 自分のプロジェクトで generics が活用できる箇所を探す
3. より高度なパターン（型制約の再帰的定義など）を学習する

## 参考資料

- [YouTube: GopherCon 2024: Advanced Generics Patterns - Axel Wagner](https://www.youtube.com/watch?v=dab3I-HcTVk)
- [PDF: GopherCon 2024: Advanced Generics Patterns - Axel Wagner](https://github.com/gophercon/2024-talks/tree/main/AxelWagner-AdvancedGenericsPatterns)

