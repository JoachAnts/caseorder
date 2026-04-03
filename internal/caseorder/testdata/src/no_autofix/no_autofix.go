package no_autofix

func noAutofix(x string) {
	switch x {
	case "banana":
	case "apple": // want `case "apple" should come before "banana"`
	}
}
