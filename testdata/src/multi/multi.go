package multi

func multi(x int) {
	switch x {
	case 3, 1, 2: // want "case value 1 should come before 3"
	case 0: // want "case 0 should come before 1"
	}
}

func multi2(x int) {
	switch x {
	case 5:
	case 10, 20, 1: // want "case value 1 should come before 20" "case 1 should come before 5"
	}
}
