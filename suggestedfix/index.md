# Go Suggested Fix ツールの作成 Codelab

このコードラボでは、`interface{}`を`any`に置き換える修正提案を提供するGo解析ツールを段階的に作成します。実践的な例を通じて、Go ASTの仕組みと静的解析ツールの開発方法を学習します。
各ステップで `skeleton/stepX` のコードに手を加えたあと `solution/stepX` と比較し、なぜその実装になるのかを確認できるよう構成しています。

## 学習目標

- Go AST と inspector を使用した効率的な型検査方法を理解する
- 段階的なアナライザー開発のアプローチを習得する
- SuggestedFix による自動修正機能の実装方法を学ぶ

## 📋 進め方

1. 各ステップの「学習内容」と「ゴール」を読み、追加したい診断や修正提案の振る舞いを整理します
2. まず `go run skeleton/stepX/main.go testdata/sample.go` を実行し、現状の出力を確認して課題を把握します
3. `skeleton/stepX/main.go` を編集し、解析ロジックを少しずつ追加しながら再実行して振る舞いの変化を確かめます
4. 期待どおりの診断や SuggestedFix が出るまで `go run skeleton/stepX/main.go testdata/sample.go`（必要なら `-fix` オプションも）を繰り返し、挙動を検証します
5. 完成形の考え方を整理するために、対応する `solution/stepX/main.go` を読み、差分の理由を言語化します
6. 行き詰まったら「ヒント」や参考資料で理解を補強し、再び skeleton に戻って実装をブラッシュアップします

### 前提条件
- Go 1.25 以降がインストールされていること
- Go プログラミングの基本的な理解
- Go 対応のテキストエディタ（VS Code、GoLand、gopls付きVimなど）

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
│   ├── step2/
│   └── step3/
├── solution/      # 各ステップの完成コード
│   ├── step1/
│   ├── step2/
│   └── step3/
└── testdata/      # テスト用データ
    └── sample.go
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

## Step 1: Inspector を使用したインターフェース型の検出

### ゴール
✅ `inspector` を使用して AST から `interface{}` 型リテラルを効率的に見つける

### 学習内容

静的解析ツールを作成する第一歩として、コードから特定の型（今回は interface 型）を検出する方法を学びます。

### なぜ Inspector を使うのか？

#### 概念対比：通常の ast.Walk vs Inspector

**通常の ast.Walk の場合：**
```go
ast.Inspect(node, func(n ast.Node) bool {
    // すべてのノードを訪問（非効率）
    switch n := n.(type) {
    case *ast.InterfaceType:
        // 処理
    }
    return true
})
```

**Inspector を使った場合：**
```go
inspect.Preorder([]ast.Node{(*ast.InterfaceType)(nil)}, func(n ast.Node) {
    // InterfaceType ノードのみを効率的に訪問
    iface := n.(*ast.InterfaceType)
    // 処理
})
```

**メリット：**
1. **パフォーマンス**: 事前インデックス化により高速な走査が可能
2. **簡潔性**: 特定の型のノードのみに絞った処理が簡単
3. **メモリ効率**: 不要なノードをスキップできる

### AST の基礎知識

Go のソースコードは以下のような AST（Abstract Syntax Tree）に変換されます：

```go
// ソースコード
var x interface{}

// AST表現（概念図）
*ast.GenDecl
  └── *ast.ValueSpec
      └── Type: *ast.InterfaceType
          └── Methods: nil (空のインターフェース)
```

### 実装タスク

`skeleton/step1/main.go` を以下のように修正してください：

1. **nodeFilter を作成**
   - InterfaceType ノードのみを対象とするフィルタを定義

2. **Preorder で走査**
   - inspector を使って効率的にノードを訪問
   - 検出したインターフェース型の位置を報告

3. **完成コードと照合して理解を深める**
   - `solution/step1/main.go` の同じ箇所を読み、抽象化の仕方やログ出力の違いを確認

### 実装のヒント

```go
// 1. nodeFilter の定義
nodeFilter := []ast.Node{
    (*ast.InterfaceType)(nil),  // InterfaceType ノードのみを対象
}

// 2. inspector による走査
inspect.Preorder(nodeFilter, func(n ast.Node) {
    iface := n.(*ast.InterfaceType)
    // インターフェース型を検出したことを報告
    pass.Reportf(iface.Pos(), "interface型を検出しました")
})
```

### 実行して確認

```bash
# 修正前を実行（何も検出されない）
go run skeleton/step1/main.go testdata/sample.go

# 修正後を実行（interface型が検出される）
go run skeleton/step1/main.go testdata/sample.go

# 期待される出力
testdata/sample.go:3:8: interface型を検出しました
testdata/sample.go:7:14: interface型を検出しました
testdata/sample.go:11:23: interface型を検出しました

# 参考解答の確認
go run solution/step1/main.go testdata/sample.go
```

### ポイント
- `inspector.Preorder` で特定の型のノードのみを効率的に走査
- AST ノードの型（`*ast.InterfaceType`）を理解することが重要
- `pass.Reportf` で検出結果を報告

### ✅ チェックリスト
- [ ] nodeFilter に InterfaceType を指定した
- [ ] Preorder を使って走査を実装した
- [ ] interface 型を検出できることを確認した
- [ ] inspector の利点を理解した

---

## Step 2: 空のインターフェースの判定

### ゴール
✅ 検出したインターフェース型から、`interface{}` のみを特定して報告する

### 学習内容

すべてのインターフェース型を検出するだけでなく、メソッドを持たない空のインターフェース（`interface{}`）のみを特定する方法を学びます。

### InterfaceType の構造理解

```go
type InterfaceType struct {
    Interface token.Pos  // "interface" キーワードの位置
    Methods   *FieldList // メソッドのリスト
    Incomplete bool       // 構文エラーがある場合 true
}
```
- [InterfaceType ドキュメント](https://pkg.go.dev/go/ast#InterfaceType)

**判定ロジック：**
- `Methods` が `nil` → `interface{}`
- `Methods.List` が空 → `interface{}`
- `Methods.List` に要素がある → 通常のインターフェース

### 実装の考え方

```go
// interface{} の例
var empty interface{}  // Methods == nil

// 通常のインターフェースの例
type Reader interface {
    Read([]byte) (int, error)  // Methods.List に1つのメソッド
}
```

### 実装タスク

`skeleton/step2/main.go` を修正：

1. **空のインターフェースの判定を追加**
   - Methods フィールドをチェック
   - 空の場合のみレポート

2. **適切なメッセージで報告**
   - より詳細な診断情報を提供

3. **solution を読み合わせ**
   - `solution/step2/main.go` を開き、条件判定や補助関数の使い方が自分の実装とどう違うか整理

### 実装のヒント

```go
inspect.Preorder(nodeFilter, func(n ast.Node) {
    iface := n.(*ast.InterfaceType)

    // 空のインターフェースかチェック
    if iface.Methods == nil || len(iface.Methods.List) == 0 {
        pass.Report(analysis.Diagnostic{
            Pos:     iface.Pos(),
            End:     iface.End(),
            Message: "interface{}の代わりにanyを使用してください",
        })
    }
})
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

### 実行して確認

```bash
# 修正前を実行（すべてのinterface型を報告）
go run skeleton/step2/main.go testdata/sample.go

# 修正後を実行（interface{}のみを報告）
go run skeleton/step2/main.go testdata/sample.go

# 期待される出力
testdata/sample.go:3:8: interface{}の代わりにanyを使用してください
testdata/sample.go:7:14: interface{}の代わりにanyを使用してください

# 参考解答の確認
go run solution/step2/main.go testdata/sample.go
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

単なる警告だけでなく、具体的な修正方法を提供する SuggestedFix の実装方法を学びます。これにより、開発者は手動で修正する手間を省けます。

### SuggestedFix の仕組み

```go
type SuggestedFix struct {
    Message   string      // 修正の説明
    TextEdits []TextEdit  // 実際の編集内容
}

type TextEdit struct {
    Pos     token.Pos  // 編集開始位置
    End     token.Pos  // 編集終了位置
    NewText []byte     // 新しいテキスト
}
```

**動作の流れ：**
1. アナライザーが問題を検出
2. SuggestedFix で修正方法を提案
3. エディタが「Quick Fix」として表示
4. ユーザーがワンクリックで適用

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

`skeleton/step3/main.go` を修正：

1. **SuggestedFix を作成**
   - TextEdit で置換内容を定義
   - 適切なメッセージを設定

2. **Diagnostic に追加**
   - SuggestedFixes フィールドに設定

3. **完成版で適用例を確認**
   - `solution/step3/main.go` を実行・読解し、`-fix` オプション対応や補助関数の分割意図を把握

### 実装のヒント

```go
if iface.Methods == nil || len(iface.Methods.List) == 0 {
    // 修正提案を作成
    fix := analysis.SuggestedFix{
        Message: "interface{}をanyに置換",
        TextEdits: []analysis.TextEdit{
            {
                Pos:     iface.Pos(),
                End:     iface.End(),
                NewText: []byte("any"),
            },
        },
    }

    // 診断に修正提案を追加
    pass.Report(analysis.Diagnostic{
        Pos:            iface.Pos(),
        End:            iface.End(),
        Message:        "interface{}の代わりにanyを使用してください",
        SuggestedFixes: []analysis.SuggestedFix{fix},
    })
}
```

### 実際の適用例

```go
// 修正前
var data interface{}
func Process(v interface{}) {}

// 自動修正後
var data any
func Process(v any) {}
```

### 実行して確認

```bash
# 修正前を実行（修正提案なし）
go run skeleton/step3/main.go testdata/sample.go

# 修正後を実行（修正提案付き）
go run skeleton/step3/main.go testdata/sample.go

# 期待される出力
testdata/sample.go:3:8: interface{}の代わりにanyを使用してください
testdata/sample.go:7:14: interface{}の代わりにanyを使用してください
# (注: コマンドライン出力は同じですが、内部的に修正提案が追加されています)

# -fix フラグで自動修正を適用（実装次第）
go run solution/step3/main.go -fix testdata/sample.go

# 参考解答の確認
go run solution/step3/main.go testdata/sample.go
```

### エディタでの活用

VS Code や GoLand などのエディタでは：
1. 問題箇所に波線が表示される
2. カーソルを合わせると診断メッセージが表示
3. 「Quick Fix」や「💡」アイコンをクリック
4. 提案された修正を選択して適用

### ポイント
- TextEdit で具体的な編集内容を定義
- SuggestedFix を Diagnostic に追加することが重要
- エディタとの連携により開発効率が大幅に向上

### ✅ チェックリスト
- [ ] SuggestedFix の作成を実装した
- [ ] TextEdit で正しい位置と内容を指定した
- [ ] Diagnostic に SuggestedFixes を追加した
- [ ] 修正提案が機能することを理解した
- [ ] 実行して動作を確認した

---

## まとめ

このコードラボで学んだこと：

### 📚 Step 1: Inspector の活用
- AST の基本構造と走査方法
- Inspector による効率的なノード検索
- 特定の型のノードのみを対象とした処理

### 📚 Step 2: 条件による絞り込み
- AST ノードの詳細な構造の理解
- 空のインターフェースの判定ロジック
- より精度の高い診断の実装

### 📚 Step 3: 自動修正の提供
- SuggestedFix による修正提案の仕組み
- TextEdit を使った具体的な編集内容の定義
- 開発効率を向上させる自動化ツールの作成

各ステップで `skeleton/stepX` の実装を仕上げたあとに `solution/stepX` を読み返し、補助関数の切り出し方や診断メッセージの粒度を確認すると、静的解析ツールの設計意図がより鮮明になります。

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

1. `solution/*` のコードを詳しく読み、実装の詳細を確認
2. 他の修正提案（例：error チェック、命名規則）を実装してみる
3. [golang.org/x/tools/go/analysis](https://pkg.go.dev/golang.org/x/tools/go/analysis) のドキュメントを読む
4. 実際のプロジェクトで独自のアナライザーを作成

## 参考資料

- [Go AST Visualizer](https://yuroyoro.github.io/goast-viewer/)
- [go/analysis package documentation](https://pkg.go.dev/golang.org/x/tools/go/analysis)
- [Writing a Go Analyzer](https://arslan.io/2019/06/13/using-go-analysis-to-write-a-custom-linter/)
- [staticcheck source code](https://github.com/dominikh/go-tools) - 実践的な例として
- [Goで作る静的解析ツール開発入門](https://zenn.dev/hsaki/books/golang-static-analysis)
