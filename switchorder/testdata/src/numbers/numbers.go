package numbers

func numbers() {
	switch 1 {
	case 2:
	case 1: // want "case 1 should come before 2"
	case 0: // want "case 0 should come before 1"
	}

	switch 1.0 {
	case 1.1:
	case 1.0: // want "case 1.0 should come before 1.1"
	case 0.9: // want "case 0.9 should come before 1.0"
	}

	switch 'a' {
	case 'c':
	case 'b': // want "case 'b' should come before 'c'"
	case 'a': // want "case 'a' should come before 'b'"
	}
}
