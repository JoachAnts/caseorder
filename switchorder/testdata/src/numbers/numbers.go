package numbers

func numbers() {
	switch 1 {
	case 1:
	case 0: // want "case 0 should come before 1"
	}
}
