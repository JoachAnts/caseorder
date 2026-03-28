package switchorder_test

import (
	"testing"

	"github.com/JoachAnts/switch-order"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSwitchOrderAlphabetical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "alphabetical")
}

func TestSwitchOrderWithFallthrough(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.RunWithSuggestedFixes(t, testdata, switchorder.Analyzer, "fallthru")
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
