package comments

func comments(x int) {
	// Switch comment
	switch x {
	// Outer comment for 2
	case 2:
		// Comment for 2
		println(2)
	/**
	Multi line comment for 1
	*/
	case 1: // want "case 1 should come before 2"
		/* Comment for 1 */
		println(1)
	}
}
