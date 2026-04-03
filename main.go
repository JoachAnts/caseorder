package main

import (
	"github.com/JoachAnts/caseorder/internal/caseorder"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	cfg := caseorder.DefaultConfig()
	a := caseorder.NewWithConfig(&cfg)

	a.Flags.StringVar(&cfg.Order, "order", cfg.Order, "sort order (asc or desc)")
	a.Flags.BoolVar(&cfg.DefaultLast, "default-last", cfg.DefaultLast, "place the default case last")
	a.Flags.BoolVar(&cfg.Autofix.Enabled, "autofix", cfg.Autofix.Enabled, "emit suggested fixes")
	a.Flags.BoolVar(&cfg.Autofix.AllowFallthrough, "autofix-allow-fallthrough", cfg.Autofix.AllowFallthrough, "emit suggested fixes for switches with fallthrough")

	for i := range cfg.Comparators {
		if cfg.Comparators[i].Type == "alphabetical" {
			a.Flags.BoolVar(&cfg.Comparators[i].IgnoreCase, "ignore-case", cfg.Comparators[i].IgnoreCase, "compare string cases case-insensitively")
			break
		}
	}

	singlechecker.Main(a)
}
