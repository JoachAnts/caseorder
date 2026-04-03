package negative

func negative(x int) {
	switch x {
	case 1:
	case -1: // want "case -1 should come before 1"
	case 0:
	}
}
