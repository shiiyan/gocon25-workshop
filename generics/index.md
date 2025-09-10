# Go ジェネリクス Codelab

このコードラボでは、Go 1.18で導入されたジェネリクスについて、実践的な例を通じて学習します。

## 学習目標

- Type constraint（型制約）を使った汎用的な実装方法を理解する
- Interface を型制約として活用する方法を学ぶ
- 複雑な型制約の実装パターンを習得する

## 📋 進め方

1. 各ステップの「学習内容」と「ゴール」を読みます
2. `skeleton/stepX/main.go` を編集して実装します
3. 修正前後の動作を確認します：
   - `go run skeleton/stepX/main.go` （修正前：エラーまたは型アサーション版）
   - `go run skeleton/stepX/main.go` （修正後：型安全版）
4. 必要なら「ヒント」を参照します
5. 自力で解けたら `solution/stepX/main.go` と見比べて理解を深めます

### 前提条件
- Go 1.25 を推奨
- 実行例: `go run skeleton/step1/main.go`
- 参考解答: `go run solution/step1/main.go`

---

## Step 1: ジェネリクスの基礎

### ゴール
✅ 1つの実装を複数の型で再利用できることを体験する

### 学習内容

ジェネリクスを使うことで、型安全性を保ちながら複数の型に対して動作する汎用的なコードを書くことができます。

**現在のコードの問題点：**
1. **型安全性の欠如**: 型アサーション `v.(int)` が必要
2. **実行時エラーのリスク**: 間違った型でアクセスするとパニックが発生
3. **コンパイル時のチェック不可**: エラーは実行時まで発見できない

### 実装タスク

`skeleton/step1/main.go` を以下のように修正してください：

1. **`Slice` 型をジェネリクス型に変更**
   ```go
   type Slice[T any] []T
   ```

2. **`Filter` メソッドを型安全に変更**
   - 型アサーションを削除
   - 関数の引数を `func(T) bool` に変更

3. **`Map` 関数を2つの型パラメータで実装**
   ```go
   func Map[A, B any](s Slice[A], f func(A) B) Slice[B]
   ```

### 実行して確認

```bash
# 修正前を実行（型アサーションが必要でパニックのリスクあり）
go run skeleton/step1/main.go

# 修正後を実行（型安全）
go run skeleton/step1/main.go

# 期待される出力
main.Slice[int]{2, 4}
main.Slice[string]{"1", "2", "3", "4", "5"}

# 参考解答の確認
go run solution/step1/main.go
```

### ポイント
- `[T any]` は型パラメータの宣言（T は任意の型）
- 型パラメータを使うことで、コンパイル時に型チェックが可能
- 型アサーションが不要になり、実行時エラーのリスクが減少
- 型パラメータ名は慣例的に大文字（`T`, `A`, `B`）を使用

### ✅ チェックリスト
- [ ] Slice型に型パラメータ `[T any]` を追加した
- [ ] Filterメソッドから型アサーションを削除した
- [ ] Map関数に2つの型パラメータを使用した
- [ ] 実行してエラーなく動作することを確認した

---

## Step 2: Interface を使った型制約

### ゴール
✅ Interface を関数引数ではなく型制約として使う理由を理解する

### 学習内容

型制約として interface を使うことで、特定のメソッドを持つ型のみを受け付ける汎用的な実装が可能になります。

### なぜ型制約として interface を使うのか？

#### 概念対比：通常の interface 引数 vs 型制約

**通常の interface 引数の場合：**
```go
// 異なる型を混在させられるが、fmt.Stringerとしてしか扱えない
func PrintItems(items []fmt.Stringer) {
    // items = []fmt.Stringer{Person{}, Product{}} // 異なる型の混在OK
    for _, item := range items {
        fmt.Println(item.String()) // fmt.Stringerとしてのみアクセス
    }
}
```

**ジェネリクスの型制約の場合：**
```go
type Container[T fmt.Stringer] struct {
    items []T  // T は fmt.Stringer を満たす「特定の型」で統一
}

// Container[Person]とContainer[Product]は別の型として扱われる
func (c *Container[T]) Add(v T) {
    c.items = append(c.items, v) // 具体的な型Tとして扱える
}
```

**メリット：**
1. **型の一貫性**: Container 内のすべての要素が同じ具体的な型で統一
2. **パフォーマンス**: インターフェースのボックス化（boxing）が不要
   - 注：コンパイラが具体的な型ごとに特化したコードを生成できる可能性がある
3. **型情報の保持**: 元の型の情報が保持され、必要に応じてアクセス可能

### 実装タスク

`skeleton/step2/main.go` の Container 型を修正：

**現在の問題：**
```go
type Container[T any] struct {  // any では String() メソッドが保証されない
    items []T
}
```
`PrintAll()` メソッド内で `item.String()` を呼び出すとコンパイルエラー

**修正方法：**
```go
type Container[T fmt.Stringer] struct {
    items []T
}
```

### 実行して確認

```bash
# 修正前を実行（コンパイルエラー）
go run skeleton/step2/main.go

# 修正後を実行
go run skeleton/step2/main.go

# 期待される出力
People:
Alice (30 years)
Bob (25 years)

Products:
Laptop: $999.99
Mouse: $25.50

# 参考解答の確認
go run solution/step2/main.go
```

### ポイント
- 型制約により、その型が持つべきメソッドをコンパイル時に保証
- interface を型制約として使うことで、柔軟かつ型安全な設計が可能
- 通常の interface 引数と異なり、具体的な型情報が保持される

### ✅ チェックリスト
- [ ] Container の型制約を `any` から `fmt.Stringer` に変更した
- [ ] PrintAll() メソッドでコンパイルエラーが解消された
- [ ] Person と Product の両方で Container が動作することを確認した
- [ ] 型制約と interface 引数の違いを理解した

---

## Step 3: 複雑な型制約

### ゴール
✅ ポインタ型と interface を組み合わせた高度な型制約パターンを理解する

### 学習内容

JSON Unmarshal の問題を通じて、複数の制約を組み合わせた型制約パターンを学びます。

### 背景：JSON Unmarshal の課題

通常の JSON unmarshal：
```go
var user User
err := json.Unmarshal(data, &user)  // &user (ポインタ) を渡す必要がある
```

**ジェネリクスで汎用的な Unmarshal 関数を作る際の課題：**
- 値型 T を受け取り、内部でポインタ *T を作成する必要がある
- *T が `json.Unmarshaler` interface を実装している必要がある
- これらの制約を同時に満たす必要がある

### 複雑な型制約の解決策

```go
type Unmarshaller[T any] interface {  // スペル注意：コード内ではこの綴り
    *T                    // ポインタ型の制約
    json.Unmarshaler      // interface の実装
}
```

この制約（インターフェースの合成/intersection）により：
1. `*T` がポインタ型であることを保証
2. `json.Unmarshaler` interface を実装していることを保証

### 実装の仕組み

```go
func Unmarshal[T any, PT Unmarshaller[T]](data []byte) (T, error) {
    var v T              // 値型の変数を作成
    err := PT(&v).UnmarshalJSON(data)  // PT型（*T）に変換して呼び出し
    return v, err
}
```

**型パラメータの役割：**
- `T`: デコードしたい実体の型（例: `User`）
- `PT`: `Unmarshaller[T]` を満たす型（`*T` かつ `json.Unmarshaler`）

### 型推論による簡略化

**Go バージョンによる違い：**

```go
// Go 1.18-1.19: 明示的な指定が必要
user, err := Unmarshal[User, *User](data)

// Go 1.20以降: 型推論により第2パラメータを省略可能
user, err := Unmarshal[User](data)  // PT は自動的に *User と推論される
```

**型推論の仕組み：**
- コンパイラは `T = User` から `PT = *User` を自動的に導出
- `Unmarshaller[User]` を満たす型は `*User` しかないため推論可能

### 理解すべきポイント

`skeleton/step3/main.go` はすでに完成しています。以下を理解してください：

1. **`Unmarshaller` interface の定義**
   - なぜ `*T` と `json.Unmarshaler` の両方が必要か
   - インターフェースの合成（intersection）の仕組み

2. **`Unmarshal` 関数の実装**
   - 2つの型パラメータの役割
   - `PT(&v)` で型変換している理由（ポインタレシーバーの要件を満たすため）

3. **型推論の活用**
   - `main` 関数での呼び出しで第2パラメータが省略されている
   - Go 1.20以降の改善により、APIがよりシンプルに

### 実行して確認

```bash
# 実行（すでに完成済み）
go run skeleton/step3/main.go
# Output: Alice 30

# 参考解答の確認
go run solution/step3/main.go
```

### ポイント
- 複数の制約を組み合わせることで、より厳密な型安全性を実現
- ポインタ型の制約により、メソッドレシーバーの要件を満たす
- 型推論により、使いやすいAPIを提供
- 現実的な問題（JSON デコード）への応用例

### ✅ チェックリスト
- [ ] Unmarshaller interface の複合制約を理解した
- [ ] 2つの型パラメータ（T, PT）の役割を理解した
- [ ] 型推論によるPTの省略がどう機能するか理解した
- [ ] ポインタ制約が必要な理由を理解した
- [ ] 実行して動作を確認した

---

## まとめ

このコードラボで学んだこと：

### 📚 Step 1: ジェネリクスの基礎
- 型パラメータによる型安全な汎用実装
- 型アサーションの削除とコンパイル時チェック
- `any` も「制約」の1つであること

### 📚 Step 2: Interface を型制約として使う
- 通常の interface 引数との違いと使い分け
- 型の一貫性とパフォーマンスの向上
- 具体的な型情報の保持によるメリット

### 📚 Step 3: 複雑な型制約
- ポインタ型と interface の組み合わせ（複合制約）
- 型推論による使いやすいAPI設計
- 現実的な問題（JSON デコード）への応用

## 次のステップ

1. `solution/*` のコードを詳しく読み、実装の詳細を確認
2. 自分のプロジェクトでジェネリクスを活用できる場面を探す
3. より高度な型制約パターンを学習

## 参考資料

- [Go Generics Tutorial](https://go.dev/doc/tutorial/generics)
- [Type Parameters Proposal](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md)
- [Go 1.18 Release Notes](https://go.dev/doc/go1.18#generics)
- [Go 1.20 Release Notes - Type inference improvements](https://go.dev/doc/go1.20#generics)
