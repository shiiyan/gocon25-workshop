package suggestedfix

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := analysistest.TestData()
	results := analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, "a")

	if len(results) != 1 {
		t.Fatalf("unexpected package count: got %d want 1", len(results))
	}

	fset := results[0].Action.Package.Fset
	var checked bool

	for _, diag := range results[0].Action.Diagnostics {
		if diag.Message != message {
			continue
		}
		pos := fset.Position(diag.Pos)
		checked = true
		if len(diag.SuggestedFixes) == 0 {
			t.Fatalf("diagnostic at %s is missing suggested fixes", pos)
		}
	}

	if !checked {
		t.Fatalf("no diagnostics reported with expected message")
	}
}
