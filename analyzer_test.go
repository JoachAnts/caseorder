package caseorder_test

import (
	"testing"

	"github.com/JoachAnts/caseorder"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSwitchOrderAlphabetical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "alphabetical")
}

func TestSwitchOrderWithFallthrough(t *testing.T) {
	testdata := analysistest.TestData()
	// Default config has allow-fallthrough: false, so diagnostics are reported but no fix is suggested.
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "fallthru")
}

func TestSwitchOrderWithFallthroughAutofix(t *testing.T) {
	cfg := caseorder.DefaultConfig()
	cfg.Autofix.AllowFallthrough = true
	a := caseorder.NewWithConfig(&cfg)
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, a, "fallthru_autofix")
}

func TestSwitchOrderNumerical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "numbers")
}

func TestSwitchOrderMulti(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "multi")
}

func TestSwitchOrderLarge(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "large")
}

func TestSwitchOrderDefault(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "default")
}

func TestSwitchOrderEdgeCases(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "edgecases")
}
