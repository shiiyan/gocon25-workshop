package main

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

// これは引っ掛け
// これがあると analysis 側で定義している `-fix` フラグとコンフリクトする
var fix string
func init() {
	Analyzer.Flags.StringVar(&fix, "fix", "", "fix with suggested fixes")
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

		// Check if it's an empty interface (interface{})
		if iface.Methods == nil || len(iface.Methods.List) == 0 {
			pos := iface.Pos()
			end := iface.End()

			pass.Report(analysis.Diagnostic{
				Pos:     pos,
				End:     end,
				Message: "interface{} can be replaced with any",
				SuggestedFixes: []analysis.SuggestedFix{ // ここを空にする予定
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
