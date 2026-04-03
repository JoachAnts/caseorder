package desc_numbers

func descNumbers() {
	switch 0 {
	case 1:
	case 2: // want "case 2 should come before 1"
	case 3: // want "case 3 should come before 2"
	}
}
