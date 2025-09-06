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
