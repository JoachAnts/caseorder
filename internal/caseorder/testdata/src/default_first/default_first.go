package default_first

func defaultFirst(x int) {
	switch x {
	case 1:
	case 2:
	default: // want "case default should come before 2"
		println("default")
	}
}
