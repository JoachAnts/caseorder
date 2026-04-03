package expressions

// Expression switches have cases that are boolean conditions, not constant
// values. Multiple cases can evaluate to true for the same input, so
// reordering them would change behaviour. caseorder must not flag these.

func overlap(x int) {
	switch {
	case x > 5:
		println("big")
	case x > 3:
		println("medium")
	case x > 1:
		println("small")
	}
}

func mixedConditions(x, y int) {
	switch {
	case x > 10 && y < 0:
		println("a")
	case x > 5:
		println("b")
	case y == 0:
		println("c")
	}
}
