# Go Suggested Fix ツールの作成

## 概要

このコードラボでは、`interface{}`を`any`に置き換える修正提案を提供するGo解析ツールを段階的に作成します。各ステップでコードを徐々に改善し、最終的に自動修正可能なツールを完成させます。

**学習内容:**
- Go ASTとinspectorを使用した型の検査方法
- 段階的なアナライザー開発のアプローチ
- 修正提案（SuggestedFix）の実装方法

**前提条件:**
- Go 1.25以降がインストールされていること
- Goプログラミングの基本的な理解
- Go対応のテキストエディタ（VS Code、GoLand、gopls付きVimなど）

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

## ステップ1: Inspectorを使用したインターフェース型の検出

### 目標
`inspector`を使用してASTから`interface{}`型リテラルを見つける基本的なアナライザーを作成します。

### 考え方
1. **なぜinspectorを使うのか？**
   - `inspector`は事前にASTをインデックス化し、特定の型のノードを効率的に走査できます
   - 単純な`ast.Walk`より高速で、特定のノードタイプに絞った処理が可能です

2. **InterfaceTypeノードの特定**
   - Go ASTでは、`interface{}`は`*ast.InterfaceType`として表現されます
   - `inspector.Preorder`で特定の型のノードだけを訪問できます

### スケルトンコード（skeleton/step1/analyzer.go）

```go
package suggestedfix

import (
    "go/ast"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "interfacetoany",
	Doc:      "check for interface{} and suggest replacing with any",
    Requires: []*analysis.Analyzer{inspect.Analyzer},
    Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		pass.Report(analysis.Diagnostic{
			Pos:     n.Pos(),
			Message: "interface found",
		})
	})

    return nil, nil
}
```

### 解決へのアプローチ

**ステップ1-1: nodeFilterの作成**
```go
nodeFilter := []ast.Node{
    (*ast.InterfaceType)(nil),  // InterfaceTypeノードのみを対象
}
```

**ステップ1-2: Preorderでの走査**
```go
inspect.Preorder(nodeFilter, func(n ast.Node) {
    iface := n.(*ast.InterfaceType)
    // この時点で全てのinterface型を検出
    pass.Reportf(iface.Pos(), "interface型を検出しました")
})
```

### 完成コード（solution/step1/analyzer.go）

```go
package suggestedfix

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "interfacetoany",
	Doc:      "check for interface{} and suggest replacing with any",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.InterfaceType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		iface, ok := n.(*ast.InterfaceType)
		if !ok {
			return
		}

		pass.Report(analysis.Diagnostic{
			Pos:     iface.Pos(),
			End:     iface.End(),
			Message: "interface{} can be replaced with any",
		})
	})

	return nil, nil
}
```

これでまず`interface{}`型を検出する基本的なアナライザーが完成しました。次のステップでは、これらのインターフェースが空であるかどうかを判定し、`any`への置換を提案します。

---

## ステップ2: 空のインターフェースの判定

### 目標
検出したインターフェース型から、メソッドを持たない空のインターフェース（`interface{}`）だけを抽出します。

### 考え方
1. **InterfaceTypeの構造理解**
   - `ast.InterfaceType`は`Methods`フィールドを持ちます
   - `Methods`が`nil`または空のリストの場合、それは`interface{}`です

2. **なぜメソッドのチェックが必要か？**
   - `interface{ Read() }`のような通常のインターフェースは変換対象外
   - `interface{}`のみが`any`への置換対象

### スケルトンコード（skeleton/step2/analyzer.go）

```go
package suggestedfix

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "interfacetoany",
	Doc:      "check for interface{} and suggest replacing with any",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.InterfaceType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		iface := n.(*ast.InterfaceType)
		pass.Report(analysis.Diagnostic{
			Pos:     iface.Pos(),
			End:     iface.End(),
			Message: "interface{} can be replaced with any",
		})
	})

	return nil, nil
}
```

### 解決へのアプローチ

**ステップ2-1: 空のインターフェースの判定ロジック**
```go
if iface.Methods == nil || len(iface.Methods.List) == 0 {
    // これは空のインターフェース
}
```

**ステップ2-2: 適切なメッセージでレポート**
```go
pass.Report(analysis.Diagnostic{
    Pos:     iface.Pos(),
    End:     iface.End(),
    Message: "interface{}の代わりにanyを使用してください",
})
```

### 完成コード（solution/step2/analyzer.go）

```go
package suggestedfix

import (
    "go/ast"
    "golang.org/x/tools/go/analysis"
    "golang.org/x/tools/go/analysis/passes/inspect"
    "golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
    Name:     "emptyinterface",
    Doc:      "interface{}をanyに置換することを提案",
    Requires: []*analysis.Analyzer{inspect.Analyzer},
    Run:      run,
}

func run(pass *analysis.Pass) (interface{}, error) {
    inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)
    
    nodeFilter := []ast.Node{
        (*ast.InterfaceType)(nil),
    }
    
    inspect.Preorder(nodeFilter, func(n ast.Node) {
        iface, ok := n.(*ast.InterfaceType)
        if !ok {
            return
        }

        // 空のインターフェースかチェック
        if iface.Methods == nil || len(iface.Methods.List) == 0 {
            pass.Report(analysis.Diagnostic{
                Pos:     iface.Pos(),
                End:     iface.End(),
                Message: "interface{}の代わりにanyを使用してください",
            })
        }
    })
    
    return nil, nil
}
```

これで、空のインターフェースを検出し、適切なメッセージで報告するアナライザーが完成しました。次のステップでは、これらの診断に自動修正提案を追加します。

---

## ステップ3: SuggestedFixによる自動修正の追加

### 目標
診断にSuggestedFixを追加し、エディタやツールで自動修正できるようにします。

### 考え方
1. **SuggestedFixの役割**
   - 単なる警告ではなく、具体的な修正方法を提供
   - エディタは「Quick Fix」として表示し、ワンクリックで適用可能

2. **TextEditの構造**
   - `Pos`: 置換開始位置
   - `End`: 置換終了位置
   - `NewText`: 新しいテキスト

3. **なぜこれが重要か？**
   - 手動修正は時間がかかりエラーが起きやすい
   - 大規模なコードベースでは自動修正が必須

### スケルトンコード（skeleton/step3/analyzer.go）

```go
package suggestedfix

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "interfacetoany",
	Doc:      "check for interface{} and suggest replacing with any",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.InterfaceType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		iface, ok := n.(*ast.InterfaceType)
		if !ok {
			return
		}

		if iface.Methods == nil || len(iface.Methods.List) == 0 {
			pass.Report(analysis.Diagnostic{
				Pos:     iface.Pos(),
				End:     iface.End(),
				Message: "interface{} can be replaced with any",
			})
		}
	})

	return nil, nil
}
```

### 解決へのアプローチ

**ステップ3-1: SuggestedFixの作成**
```go
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
```

**ステップ3-2: DiagnosticにSuggestedFixesを追加**
```go
pass.Report(analysis.Diagnostic{
    Pos:            iface.Pos(),
    End:            iface.End(),
    Message:        "interface{}の代わりにanyを使用してください",
    SuggestedFixes: []analysis.SuggestedFix{fix}, // ← 重要！
})
```

### 完成コード（solution/step3/analyzer.go）

```go
package suggestedfix

import (
	"go/ast"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
	"golang.org/x/tools/go/ast/inspector"
)

var Analyzer = &analysis.Analyzer{
	Name:     "interfacetoany",
	Doc:      "check for interface{} and suggest replacing with any",
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Run:      run,
}

func run(pass *analysis.Pass) (any, error) {
	inspect := pass.ResultOf[inspect.Analyzer].(*inspector.Inspector)

	nodeFilter := []ast.Node{
		(*ast.InterfaceType)(nil),
	}

	inspect.Preorder(nodeFilter, func(n ast.Node) {
		iface, ok := n.(*ast.InterfaceType)
		if !ok {
			return
		}

		if iface.Methods == nil || len(iface.Methods.List) == 0 {
			pos := iface.Pos()
			end := iface.End()

			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				End:     end,
				Message: "interface{} can be replaced with any",
				SuggestedFixes: []analysis.SuggestedFix{
					{
						Message: "Replace interface{} with any",
						TextEdits: []analysis.TextEdit{
							{
								Pos:     pos,
								End:     end,
								NewText: []byte("any"),
							},
						},
					},
				},
			})
		}
	})

	return nil, nil
}
```

---

## テストコード

各ステップで動作を確認するためのテストコード：

### testdata/src/a/a.go

```go
package a

// interface{}を使用している例
func processData(data interface{}) { // want "interface{}の代わりにanyを使用してください"
    println(data)
}

// 既にanyを使用している例（エラーなし）
func modernFunction(data any) {
    println(data)
}

// メソッドを持つインターフェース（エラーなし）
type Reader interface {
    Read([]byte) (int, error)
}

// 複数のinterface{}
func multipleInterfaces(a interface{}, b interface{}) interface{} { // want "interface{}の代わりにanyを使用してください" "interface{}の代わりにanyを使用してください" "interface{}の代わりにanyを使用してください"
    return a
}
```

### analyzer_test.go

```go
package suggestedfix

import (
    "testing"
    "golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.Run(t, testdata, Analyzer, "a")
}

func TestAnalyzerFix(t *testing.T) {
    testdata := analysistest.TestData()
    analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "a")
}
```

---

## 重要なポイントのまとめ

### 各ステップで学んだこと

1. **ステップ1**: Inspectorを使用することで、特定の型のノードだけを効率的に走査できる
2. **ステップ2**: ASTノードの構造を理解し、条件に基づいてフィルタリングする方法
3. **ステップ3**: SuggestedFixを追加することで、ツールが実用的になる

