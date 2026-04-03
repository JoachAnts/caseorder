package caseorder_test

import (
	"testing"

	"github.com/JoachAnts/caseorder/internal/caseorder"
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

func TestOverlap(t *testing.T) {
	testdata := analysistest.TestData()
	// Expression switches (no tag) have conditions that can overlap — reordering
	// would change behaviour. The analyzer must not emit any diagnostics.
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "expressions")
}

func TestNegativeNumbers(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "negative")
}

func TestComments(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.Analyzer, "comments")
}

func TestDescendingStrings(t *testing.T) {
	cfg := caseorder.DefaultConfig()
	cfg.Order = "desc"
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.NewWithConfig(&cfg), "desc_strings")
}

func TestDescendingNumbers(t *testing.T) {
	cfg := caseorder.DefaultConfig()
	cfg.Order = "desc"
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.NewWithConfig(&cfg), "desc_numbers")
}

func TestCaseSensitive(t *testing.T) {
	cfg := caseorder.DefaultConfig()
	for i := range cfg.Comparators {
		if cfg.Comparators[i].Type == "alphabetical" {
			cfg.Comparators[i].IgnoreCase = false
		}
	}
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.NewWithConfig(&cfg), "case_sensitive")
}

func TestNoAutofix(t *testing.T) {
	cfg := caseorder.DefaultConfig()
	cfg.Autofix.Enabled = false
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, caseorder.NewWithConfig(&cfg), "no_autofix")
}
