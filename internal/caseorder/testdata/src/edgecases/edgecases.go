package edgecases

func negative(x int) {
	switch x {
	case 1:
	case -1: // want "case -1 should come before 1"
	case 0:
	}
}

func comments(x int) {
	switch x {
	case 2:
		// Comment for 2
		println(2)
	case 1: // want "case 1 should come before 2"
		/* Comment for 1 */
		println(1)
	}
}

func defaultPlacement(x int) {
	switch x {
	case 2:
	case 1: // want "case 1 should come before 2"
	default:
		println("default")
	}
}

func defaultMiddle(x int) {
	switch x {
	case 2:
	default:
		println("default")
	case 1: // want "case 1 should come before default"
	}
}
