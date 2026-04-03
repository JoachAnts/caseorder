package comments

func comments(x int) {
	// Switch comment
	switch x {
	// Outer comment for 2
	case 2:
		// Comment for 2
		println(2)
		// End comment for 2
	case 5:
		return
	/*
	 * Multi line indented comment for 3
	 * This is another line
	 */
	case 3: // want "case 3 should come before 5"
		println(1)
	/**
	Multi line comment for 1
	*/
	case 1: // want "case 1 should come before 3"
		/* Comment for 1 */
		println(1)
		/* End comment for 1 */
	}
}
