package main

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestAnalyzer(t *testing.T) {
	testdata := "/Users/kenshi.kamata/gocon25-workshop/suggestedfix/testdata"
	analysistest.RunWithSuggestedFixes(t, testdata, Analyzer, ".")
}