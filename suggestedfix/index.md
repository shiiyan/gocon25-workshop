author: kamata
summary: Introduction to building a Go Suggested Fix tool
id: suggestedfix
categories: codelab,markdown
environments: Web
status: Published

# Go Suggested Fix ツールの作成 Codelab

このコードラボでは、`interface{}`を`any`に置き換える修正提案を提供するGo解析ツールを段階的に作成します。実践的な例を通じて、[Go AST](https://pkg.go.dev/go/ast)の仕組みと静的解析ツールの開発方法を学習します。
各ステップで `skeleton/stepX` のコードに手を加えたあと `solution/stepX` と比較し、なぜその実装になるのかを確認できるよう構成しています。

## 学習目標

- [Go AST](https://pkg.go.dev/go/ast) と [inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector) を使用した効率的な型検査方法を理解する
- 段階的なアナライザー開発のアプローチを習得する
- SuggestedFix による自動修正機能の実装方法を学ぶ

## 📋 進め方

1. 各ステップの「学習内容」と「ゴール」を読み、追加したい診断や修正提案の振る舞いを整理します
2. まず `go test ./skeleton/stepX/...` を実行し、現状のテストの失敗状況を確認して課題を把握します
3. `skeleton/stepX/analyzer.go` を編集し、解析ロジックを少しずつ追加しながら再実行して振る舞いの変化を確かめます
4. 期待どおりの診断や SuggestedFix が出るまで `go test ./skeleton/stepX/...` を繰り返し、挙動を検証します
5. 完成形の考え方を整理するために、対応する `solution/stepX/analyzer.go` を読み、差分の理由を言語化します
6. 行き詰まったら「ヒント」や参考資料で理解を補強し、再び skeleton に戻って実装をブラッシュアップします


### 実装タスク

`skeleton/stepX/analyzer.go` をテストが通るように修正してください：

```bash
cd suggestedfix/skeleton/stepX

# テストを実行
go test -v -count=1 .
```

### 前提条件
- Go 1.25 以降がインストールされていること
- Go プログラミングの基本的な理解

---

## 環境のセットアップ

### Goインストールの確認

```bash
go version
```

### 作業ディレクトリの構成

このワークショップでは、以下の構成で作業します：

```
suggestedfix/
├── skeleton/      # 各ステップのスケルトンコード
│   ├── step1/
│   │   ├── analyzer.go       # 実装する解析器
│   │   ├── analyzer_test.go  # テストコード
│   │   └── testdata/         # テスト用データ
│   ├── step2/
│   └── step3/
└── solution/      # 各ステップの完成コード
    ├── step1/
    ├── step2/
    └── step3/
```

### Goモジュールの初期化

```bash
cd suggestedfix
go mod init suggestedfix
```

### 必要な依存関係のインストール

```bash
go get golang.org/x/tools/go/analysis
go get golang.org/x/tools/go/analysis/passes/inspect
go get golang.org/x/tools/go/ast/inspector
```

---

## Step 1: [Inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector) を使用したインターフェース型の検出

### ゴール
✅ [`inspector`](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector) を使用して AST から `interface{}` 型リテラルを効率的に見つける

### 学習内容

静的解析ツールを作成する第一歩として、コードから特定の型（今回は interface 型）を検出する方法を学びます。

### Inspector の仕組みと利点

#### Inspector とは

[`go/ast/inspector`](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector) パッケージは、AST を効率的に走査するための最適化されたツールです。[analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) パッケージでは、[`inspect.Analyzer`](https://pkg.go.dev/golang.org/x/tools/go/analysis/passes/inspect#Analyzer) を通じて提供されます。

#### 内部動作の仕組み

Inspector は AST を事前にインデックス化することで高速な走査を実現します：

1. **事前インデックス化**: AST 全体を一度走査し、ノードタイプごとにインデックスを作成
2. **型フィルタリング**: 指定された型のノードのみを効率的に訪問
3. **メモリ共有**: 複数のアナライザー間でインデックスを共有

#### Preorder と Postorder

```go
// Preorder: 親ノードを子ノードより先に訪問
inspect.Preorder(nodeFilter, func(n ast.Node) {
    // 親→子の順序で処理
})

// Postorder: 子ノードを親ノードより先に訪問
inspect.Postorder(nodeFilter, func(n ast.Node) {
    // 子→親の順序で処理
})
```

#### nodeFilter の仕組み

```go
// 型のゼロ値ポインタを使って、対象とする型を指定
nodeFilter := []ast.Node{
    (*ast.InterfaceType)(nil),  // InterfaceType のみ
    (*ast.FuncDecl)(nil),       // 複数指定も可能
}
```

この配列に指定された型のノードのみがコールバック関数に渡されます。

#### パフォーマンスの違い

大規模なコードベース（1000ファイル）での比較：
- **ast.Inspect**: 全ノード（約100万個）を訪問 → 約500ms
- **inspector.Preorder**: 対象ノード（約1000個）のみ訪問 → 約10ms

つまり、特定の型のノードのみを処理したい場合、[Inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector) は50倍以上高速です。

### AST の基礎知識

#### AST とは何か

AST（Abstract Syntax Tree）は、ソースコードの構文構造を木構造で表現したものです。コンパイラやツールがコードを理解し操作するための中間表現として使用されます。

#### Go における AST ノード

Go の AST は [`go/ast`](https://pkg.go.dev/go/ast) パッケージで定義されており、すべてのノードは [`ast.Node`](https://pkg.go.dev/go/ast#Node) インターフェースを実装しています：

```go
type Node interface {
    Pos() token.Pos // ノードの開始位置
    End() token.Pos // ノードの終了位置
}
```

主要なノードタイプ：
- **式（Expression）ノード**: [`ast.Expr`](https://pkg.go.dev/go/ast#Expr) を実装
  - [`*ast.Ident`](https://pkg.go.dev/go/ast#Ident): 識別子（変数名、型名など）
  - [`*ast.BasicLit`](https://pkg.go.dev/go/ast#BasicLit): リテラル（数値、文字列など）
  - [`*ast.CallExpr`](https://pkg.go.dev/go/ast#CallExpr): 関数呼び出し
  - [`*ast.InterfaceType`](https://pkg.go.dev/go/ast#InterfaceType): インターフェース型

- **文（Statement）ノード**: [`ast.Stmt`](https://pkg.go.dev/go/ast#Stmt) を実装
  - [`*ast.AssignStmt`](https://pkg.go.dev/go/ast#AssignStmt): 代入文
  - [`*ast.IfStmt`](https://pkg.go.dev/go/ast#IfStmt): if 文
  - [`*ast.ForStmt`](https://pkg.go.dev/go/ast#ForStmt): for 文

- **宣言（Declaration）ノード**: [`ast.Decl`](https://pkg.go.dev/go/ast#Decl) を実装
  - [`*ast.GenDecl`](https://pkg.go.dev/go/ast#GenDecl): 汎用宣言（var, const, type, import）
  - [`*ast.FuncDecl`](https://pkg.go.dev/go/ast#FuncDecl): 関数宣言

#### AST の走査方法

AST を走査する主な方法は3つあります：

1. **[ast.Inspect](https://pkg.go.dev/go/ast#Inspect)**: 再帰的にすべてのノードを訪問
```go
ast.Inspect(node, func(n ast.Node) bool {
    // すべてのノードに対して実行
    return true // false を返すと子ノードをスキップ
})
```

2. **[ast.Walk](https://pkg.go.dev/go/ast#Walk)**: Visitor パターンを使用
```go
type visitor struct{}
func (v *visitor) Visit(n ast.Node) ast.Visitor {
    // ノード処理
    return v
}
ast.Walk(&visitor{}, node)
```

3. **[inspector.Preorder](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector.Preorder)/[inspector.Postorder](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector.Postorder)**: 効率的な型フィルタリング（後述）

#### InterfaceType ノードの詳細

```go
// ソースコード例
var x interface{}           // 空のインターフェース
type Reader interface {     // メソッドを持つインターフェース
    Read([]byte) (int, error)
}

// AST表現
*ast.GenDecl                // 宣言ノード
  └── *ast.ValueSpec        // 値の仕様
      └── Type: *ast.InterfaceType  // インターフェース型
          └── Methods: *ast.FieldList // メソッドリスト
              └── List: []*ast.Field   // 各メソッド
```

### 実装で必要な知識

#### [analysis.Pass](https://pkg.go.dev/golang.org/x/tools/go/analysis#Pass) の役割

```go
type Pass struct {
    Analyzer   *Analyzer        // 実行中のアナライザー
    Fset       *token.FileSet   // ファイル位置情報
    Files      []*ast.File      // 解析対象ファイル
    ResultOf   map[*Analyzer]interface{} // 依存アナライザーの結果
    Report     func(Diagnostic) // 診断を報告
    // ... 他のフィールド
}
```

`pass.ResultOf[inspect.Analyzer]` から [Inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector) インスタンスを取得できます。

#### [token.Pos](https://pkg.go.dev/go/token#Pos) の概念

AST ノードの位置は [token.Pos](https://pkg.go.dev/go/token#Pos) で表現されます：
- ファイル内のバイトオフセットを表す整数値
- `pass.Fset` を使って実際のファイル位置に変換可能
- `Pos()` はノードの開始位置、`End()` は終了位置

#### 必要な修正箇所

1. nodeFilter の設定
2. Preorder コールバック内での診断報告

### ポイント
- [`inspector.Preorder`](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector.Preorder) で特定の型のノードのみを効率的に走査
- AST ノードの型（[`*ast.InterfaceType`](https://pkg.go.dev/go/ast#InterfaceType)）を理解することが重要
- `pass.Report` で検出結果を報告（Diagnostic 構造体を使用）

---

## Step 2: 空のインターフェースの判定

### ゴール
✅ 検出したインターフェース型から、[`interface{}`](https://go.dev/ref/spec#Interface_types) のみを特定して報告する

### 学習内容

すべてのインターフェース型を検出するだけでなく、メソッドを持たない空のインターフェース（`interface{}`）のみを特定する方法を学びます。

### `InterfaceType` の詳細構造

#### `InterfaceType` struct の定義

[`InterfaceType`](https://pkg.go.dev/go/ast#InterfaceType) の構造：
```go
type InterfaceType struct {
    Interface  token.Pos  // "interface" キーワードの位置
    Methods    *FieldList // メソッドのリスト
    Incomplete bool       // 構文エラーがある場合 true
}

type FieldList struct {
    Opening token.Pos // 開き括弧 '{' の位置
    List    []*Field  // フィールド（メソッド）のリスト
    Closing token.Pos // 閉じ括弧 '}' の位置
}

type Field struct {
    Doc     *CommentGroup // ドキュメントコメント
    Names   []*Ident      // フィールド/メソッド名
    Type    Expr          // 型表現
    Tag     *BasicLit     // フィールドタグ（インターフェースでは nil）
    Comment *CommentGroup // 行コメント
}
```

#### 様々なインターフェース型の AST 表現

```go
// 1. 空のインターフェース（Go 1.18以前のスタイル）
var x interface{}
// → Methods: nil または Methods.List: []

// 2. any 型（Go 1.18以降）
var y any
// → これは *ast.Ident であり、*ast.InterfaceType ではない！

// 3. メソッドを持つインターフェース
type Writer interface {
    Write([]byte) (int, error)
}
// → Methods.List: []*Field{...} （1つの要素）

// 4. 埋め込みインターフェース
type ReadWriter interface {
    io.Reader
    io.Writer
}
// → Methods.List: []*Field{...} （2つの要素、Names は nil）

// 5. 型制約インターフェース（ジェネリクス）
type Number interface {
    ~int | ~float64
}
// → Methods.List に型要素が含まれる
```

#### 空のインターフェース判定の詳細

空のインターフェースかどうかを判定する際の考慮点：

1. **Methods が nil の場合**: 明示的に `interface{}` と書かれた
2. **Methods.List が空配列の場合**: `interface { }` のようにスペース付きで書かれた
3. **Methods.List に要素がある場合**: メソッドまたは埋め込み型がある

#### Field の解釈

インターフェースの Field は以下のパターンがあります：

```go
// メソッド定義
Read([]byte) (int, error)
// → Names: []*Ident{"Read"}
// → Type: *ast.FuncType

// 埋め込みインターフェース
io.Reader
// → Names: nil
// → Type: *ast.SelectorExpr または *ast.Ident

// 型要素（ジェネリクス）
~int
// → 特殊な型表現
```

### 実装タスク

`skeleton/step2/analyzer.go` を修正：

Step1では全てのインターフェース型を検出しましたが、Step2では空のインターフェース（`interface{}`）のみを検出するように修正します。

1. **空のインターフェースの判定を追加**
   - InterfaceType の Methods フィールドをチェック
   - 空の場合（nil または List が空）のみレポート

2. **適切なメッセージで報告**
   - より詳細な診断情報を提供

3. **solution を読み合わせ**
   - `solution/step2/analyzer.go` を開き、条件判定の実装方法を確認

### 実装で必要な知識

#### 型アサーションの安全性

[Inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector) を使用する場合、nodeFilter で指定した型のノードのみが渡されるため、型アサーションは常に成功します。しかし、防御的プログラミングとして `ok` パターンを使うことも可能です：

```go
iface, ok := n.(*ast.InterfaceType)
if !ok {
    return // これは実際には到達しない
}
```

#### nil チェックの重要性

Go の AST では、オプショナルな要素は nil になることがあります：
- `Methods` フィールドが nil の場合
- `Methods.List` が nil の場合（通常は空スライス `[]` ですが）

両方のケースを考慮した条件分岐が必要です。

#### 診断メッセージの一貫性

アナライザー全体で一貫したメッセージを使用するため、定数として定義することが推奨されます：

```go
const message = "interface{} can be replaced with any"
```

### なぜ空のインターフェースだけを対象にするのか？

```go
// 置換対象（空のインターフェース）
var x interface{}  // → var x any

// 置換対象外（メソッドを持つインターフェース）
type Writer interface {
    Write([]byte) (int, error)  // これは置換しない
}
```

### ポイント
- AST ノードの構造（`Methods` フィールド）を理解する
- 条件分岐により特定のパターンのみを検出
- より具体的な診断メッセージを提供

### ✅ チェックリスト
- [ ] Methods フィールドのチェックを実装した
- [ ] 空のインターフェースのみを検出できることを確認した
- [ ] 通常のインターフェースは検出されないことを確認した
- [ ] 診断メッセージが適切に表示されることを確認した

---

## Step 3: SuggestedFix による自動修正の追加

### ゴール
✅ 診断に自動修正機能を追加し、エディタやツールで簡単に適用できるようにする

### 学習内容

単なる警告だけでなく、具体的な修正方法を提供する [SuggestedFix](https://pkg.go.dev/golang.org/x/tools/go/analysis#SuggestedFix) の実装方法を学びます。これにより、開発者は手動で修正する手間を省けます。

### SuggestedFix の詳細な仕組み

#### データ構造の階層

```go
type Diagnostic struct {
    Pos            token.Pos      // 診断の開始位置
    End            token.Pos      // 診断の終了位置
    Category       string         // カテゴリ（オプション）
    Message        string         // エラーメッセージ
    SuggestedFixes []SuggestedFix // 修正提案（複数可）
    Related        []RelatedInfo  // 関連情報
}

type SuggestedFix struct {
    Message   string      // 修正の説明（エディタに表示）
    TextEdits []TextEdit  // 実際の編集内容
}

type TextEdit struct {
    Pos     token.Pos  // 編集開始位置
    End     token.Pos  // 編集終了位置
    NewText []byte     // 置換後のテキスト
}
```

#### TextEdit の動作原理

[TextEdit](https://pkg.go.dev/golang.org/x/tools/go/analysis#TextEdit) は、ソースコードの特定範囲を新しいテキストで置換します：

```go
// 元のコード: "interface{}"
// Pos: 'i' の位置
// End: '}' の次の位置
// NewText: []byte("any")
// 結果: "any"
```

重要な特性：
- **位置の精度**: Pos と End は正確にトークンの境界を指定する必要がある
- **バイト配列**: NewText は []byte 型（UTF-8 エンコーディング）
- **複数編集**: 1つの SuggestedFix に複数の TextEdit を含められる

#### 複数の修正提案

1つの診断に複数の修正方法を提案できます：

```go
SuggestedFixes: []analysis.SuggestedFix{
    {
        Message: "Replace with any",
        TextEdits: []analysis.TextEdit{{
            Pos: pos, End: end,
            NewText: []byte("any"),
        }},
    },
    {
        Message: "Replace with generic type",
        TextEdits: []analysis.TextEdit{{
            Pos: pos, End: end,
            NewText: []byte("T"),
        }},
    },
}
```

#### 位置計算の注意点

```go
// 正しい位置の取得
iface := n.(*ast.InterfaceType)
pos := iface.Pos()  // "interface" キーワードの開始位置
end := iface.End()  // "}" の次の位置

// 間違った例
pos := iface.Interface  // これも同じだが、Pos() メソッドを使うべき
```

### なぜ SuggestedFix が重要か？

#### 手動修正 vs 自動修正

**手動修正の問題点：**
- 時間がかかる
- タイプミスのリスク
- 大規模なコードベースでは現実的でない

**自動修正のメリット：**
- 一貫性のある修正
- 時間の節約
- エラーの削減
- CI/CD パイプラインでの自動適用も可能

### 実装タスク

`skeleton/step3/analyzer.go` を修正：

Step2では空のインターフェースを検出できるようになりました。Step3では、これに自動修正機能を追加します。

1. **[SuggestedFix](https://pkg.go.dev/golang.org/x/tools/go/analysis#SuggestedFix) を作成**
   - [TextEdit](https://pkg.go.dev/golang.org/x/tools/go/analysis#TextEdit) で置換内容を定義（[`interface{}`](https://go.dev/ref/spec#Interface_types) → [`any`](https://pkg.go.dev/builtin#any)）
   - 適切なメッセージを設定

2. **[Diagnostic](https://pkg.go.dev/golang.org/x/tools/go/analysis#Diagnostic) に追加**
   - SuggestedFixes フィールドに設定
   - 既存の診断メッセージはそのまま

3. **完成版で適用例を確認**
   - `solution/step3/analyzer.go` を読解し、SuggestedFix の実装方法を理解

### 実装で必要な知識

#### [Diagnostic](https://pkg.go.dev/golang.org/x/tools/go/analysis#Diagnostic) と [SuggestedFix](https://pkg.go.dev/golang.org/x/tools/go/analysis#SuggestedFix) の関係

Diagnostic は問題を報告し、SuggestedFix はその解決方法を提供します：

```go
pass.Report(analysis.Diagnostic{
    Pos:     pos,              // 診断範囲の開始
    End:     end,              // 診断範囲の終了
    Message: "問題の説明",      // ユーザーに表示される診断メッセージ
    SuggestedFixes: []analysis.SuggestedFix{
        // 0個以上の修正提案
    },
})
```

#### 修正提案の構成要素

1. **Message**: エディタの Quick Fix メニューに表示される説明
2. **[TextEdits](https://pkg.go.dev/golang.org/x/tools/go/analysis#TextEdit)**: 実際に適用される編集操作のリスト

#### 実装時の考慮事項

- 現在の実装では診断のみを報告している
- SuggestedFixes フィールドを追加する必要がある
- Pos と End は既存の診断と同じ値を使用できる

### 実際の適用例

```go
// 修正前
var data interface{}
func Process(v interface{}) {}

// 自動修正後
var data any
func Process(v any) {}
```

### analysistest の仕組み

テストでは [`analysistest`](https://pkg.go.dev/golang.org/x/tools/go/analysis/analysistest) パッケージを使用しており、`testdata/src/a/a.go` のコメントで期待される診断を指定しています：

```go
func example1(x interface{}) interface{} { // want "interface{} can be replaced with any"
    return x
}
```

`// want` コメントがある行で、指定されたメッセージの診断が報告されることを検証します。

更に analysistest.RunWithSuggestedFixes を使用することで、SuggestedFixes の適用もテストできます。
golden ファイル（`a.go.golden`）に修正後のコードを保存し、テストで自動的に比較します。

### ポイント
- [TextEdit](https://pkg.go.dev/golang.org/x/tools/go/analysis#TextEdit) で具体的な編集内容を定義
- [SuggestedFix](https://pkg.go.dev/golang.org/x/tools/go/analysis#SuggestedFix) を [Diagnostic](https://pkg.go.dev/golang.org/x/tools/go/analysis#Diagnostic) に追加することが重要
- エディタとの連携により開発効率が大幅に向上

### ✅ チェックリスト
- [ ] [SuggestedFix](https://pkg.go.dev/golang.org/x/tools/go/analysis#SuggestedFix) の作成を実装した
- [ ] [TextEdit](https://pkg.go.dev/golang.org/x/tools/go/analysis#TextEdit) で正しい位置と内容を指定した
- [ ] [Diagnostic](https://pkg.go.dev/golang.org/x/tools/go/analysis#Diagnostic) に SuggestedFixes を追加した
- [ ] 修正提案が機能することを理解した
- [ ] 実行して動作を確認した

---

## まとめ

このコードラボで学んだこと：

### 📚 Step 1: [Inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector) の活用
- AST の基本構造と走査方法
- [Inspector](https://pkg.go.dev/golang.org/x/tools/go/ast/inspector#Inspector) による効率的なノード検索
- 特定の型のノードのみを対象とした処理

### 📚 Step 2: 条件による絞り込み
- AST ノードの詳細な構造の理解
- 空のインターフェースの判定ロジック
- より精度の高い診断の実装

### 📚 Step 3: 自動修正の提供
- [SuggestedFix](https://pkg.go.dev/golang.org/x/tools/go/analysis#SuggestedFix) による修正提案の仕組み
- [TextEdit](https://pkg.go.dev/golang.org/x/tools/go/analysis#TextEdit) を使った具体的な編集内容の定義
- 開発効率を向上させる自動化ツールの作成

## 実践的な応用

学んだ技術は以下のような場面で活用できます：

1. **コードマイグレーション**
   - 古いAPIから新しいAPIへの移行
   - 非推奨機能の置き換え

2. **コード品質の向上**
   - コーディング規約の自動適用
   - アンチパターンの検出と修正

3. **リファクタリング支援**
   - 大規模な構造変更の自動化
   - 一貫性のある変更の適用

## 次のステップ

1. **理解を深める**
   - `solution/*` のコードを詳しく読み、実装の詳細を確認
   - なぜ `ok` パターンで型アサーションをチェックするのか理解する

2. **応用練習**
   - 他の修正提案（例：error チェック、命名規則）を実装してみる
   - 複数の SuggestedFix を提供するアナライザーを作成

3. **さらなる学習**
   - [golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) のドキュメントを読む
   - 実際のプロジェクトで独自のアナライザーを作成
   - [`singlechecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/singlechecker) と [`multichecker`](https://pkg.go.dev/golang.org/x/tools/go/analysis/multichecker) の違いを理解する

## 参考資料

- [Goで作る静的解析ツール開発入門](https://zenn.dev/hsaki/books/golang-static-analysis)
- [go/analysis package documentation](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [Writing a Go Analyzer](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)
- [staticcheck source code](https://github.com/dominikh/go-tools) - 実践的な例として
