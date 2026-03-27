package multi

func multi(x int) {
	switch x {
	case 3, 1, 2: // want "case value 1 should come before 3"
	case 0: // want "case 0 should come before 1"
	}
}
