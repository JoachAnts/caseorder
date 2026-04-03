package _default

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
