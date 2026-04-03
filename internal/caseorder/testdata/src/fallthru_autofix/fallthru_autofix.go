package fallthru_autofix

func fallthru(x string) {
	switch x {
	case "banana":
	case "apple": // want "case \"apple\" should come before \"banana\""
	case "zebra":
		fallthrough
	case "yacht":
		fallthrough
	case "venus":
	case "tribe": // want "case \"tribe\" should come before \"zebra\""
	}
}
