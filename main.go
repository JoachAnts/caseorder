package main

import (
	"github.com/JoachAnts/switch-order/switchorder"
	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(switchorder.Analyzer)
}
