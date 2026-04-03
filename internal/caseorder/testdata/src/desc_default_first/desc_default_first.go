package desc_default_first

func descDefaultFirst(x int) {
	switch x {
	case 1:
	case 2: // want "case 2 should come before 1"
	default: // want "case default should come before 2"
		println("default")
	}
}
