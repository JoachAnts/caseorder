package large

func large(x int) {
	switch x {
	case 42:
	case 41: // want "case 41 should come before 42"
	case 91:
	case 9: // want "case 9 should come before 91"
	case 65:
	case 50: // want "case 50 should come before 65"
	case 1: // want "case 1 should come before 50"
	case 70:
	case 15: // want "case 15 should come before 70"
	case 78:
	case 73: // want "case 73 should come before 78"
	case 10: // want "case 10 should come before 73"
	case 55:
	case 56:
	case 72:
	case 45: // want "case 45 should come before 72"
	case 48:
	case 92:
	case 76: // want "case 76 should come before 92"
	case 37: // want "case 37 should come before 76"
	case 30: // want "case 30 should come before 37"
	case 21: // want "case 21 should come before 30"
	case 32:
	case 96:
	case 80: // want "case 80 should come before 96"
	case 49: // want "case 49 should come before 80"
	case 83:
	case 26: // want "case 26 should come before 83"
	case 87:
	case 33: // want "case 33 should come before 87"
	case 8: // want "case 8 should come before 33"
	case 47:
	case 59:
	case 63:
	case 74:
	case 44: // want "case 44 should come before 74"
	case 98:
	case 52: // want "case 52 should come before 98"
	case 85:
	case 12: // want "case 12 should come before 85"
	case 36:
	case 23: // want "case 23 should come before 36"
	case 39:
	case 40:
	case 18: // want "case 18 should come before 40"
	case 66:
	case 61: // want "case 61 should come before 66"
	case 60: // want "case 60 should come before 61"
	case 7: // want "case 7 should come before 60"
	case 34:
	case 99:
	case 46: // want "case 46 should come before 99"
	case 2: // want "case 2 should come before 46"
	case 51:
	case 16: // want "case 16 should come before 51"
	case 38:
	case 58:
	case 68:
	case 22: // want "case 22 should come before 68"
	case 62:
	case 24: // want "case 24 should come before 62"
	case 5: // want "case 5 should come before 24"
	case 6:
	case 67:
	case 82:
	case 19: // want "case 19 should come before 82"
	case 79:
	case 43: // want "case 43 should come before 79"
	case 90:
	case 20: // want "case 20 should come before 90"
	case 0: // want "case 0 should come before 20"
	case 95:
	case 57: // want "case 57 should come before 95"
	case 93:
	case 53: // want "case 53 should come before 93"
	case 89:
	case 25: // want "case 25 should come before 89"
	case 71:
	case 84:
	case 77: // want "case 77 should come before 84"
	case 64: // want "case 64 should come before 77"
	case 29: // want "case 29 should come before 64"
	case 27: // want "case 27 should come before 29"
	case 88:
	case 97:
	case 4: // want "case 4 should come before 97"
	case 54:
	case 75:
	case 11: // want "case 11 should come before 75"
	case 69:
	case 86:
	case 13: // want "case 13 should come before 86"
	case 17:
	case 28:
	case 31:
	case 35:
	case 94:
	case 3: // want "case 3 should come before 94"
	case 14:
	case 81:
	}
}
