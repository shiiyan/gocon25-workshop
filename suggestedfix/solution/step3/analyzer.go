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

const message = "interface{} can be replaced with any"

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
				Message: message,
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
