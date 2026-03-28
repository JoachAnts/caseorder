package switchorder_test

import (
	"testing"

	switchorder "github.com/JoachAnts/switch-order"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSwitchOrderAlphabetical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "alphabetical")
}

func TestSwitchOrderWithFallthrough(t *testing.T) {
	testdata := analysistest.TestData()
	// Default config has allow-fallthrough: false, so diagnostics are reported but no fix is suggested.
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "fallthru")
}

func TestSwitchOrderWithFallthroughAutofix(t *testing.T) {
	cfg := switchorder.DefaultConfig()
	cfg.Autofix.AllowFallthrough = true
	a := switchorder.NewWithConfig(cfg)
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, a, "fallthru_autofix")
}

func TestSwitchOrderNumerical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "numbers")
}

func TestSwitchOrderMulti(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "multi")
}

func TestSwitchOrderLarge(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "large")
}

func TestSwitchOrderEdgeCases(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "edgecases")
}
