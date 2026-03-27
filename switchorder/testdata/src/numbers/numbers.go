package numbers

func numbers() {
	switch 1 {
	case 2:
	case 1: // want "case 1 should come before 2"
	case 0: // want "case 0 should come before 1"
	}
}
