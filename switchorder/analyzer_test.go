package switchorder_test

import (
	"testing"

	"github.com/JoachAnts/switch-order/switchorder"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestSwitchOrderAlphabetical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, switchorder.Analyzer, "alphabetical")
}

func TestSwitchOrderNumerical(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, switchorder.Analyzer, "numbers")
}
