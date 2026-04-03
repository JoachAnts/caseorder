package desc_strings

func descStrings(x string) {
	switch x {
	case "apple":
	case "orange": // want `case "orange" should come before "apple"`
	case "zebra": // want `case "zebra" should come before "orange"`
	}
}
