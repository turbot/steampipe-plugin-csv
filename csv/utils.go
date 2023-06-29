package csv

// Convert column index number to corresponding letter
// For example, 1:a, 2:b, 27:aa, 55:bc
func intToLetters(colIndex int) (letter string) {
	colIndex--
	if firstLetter := colIndex / 26; firstLetter > 0 {
		letter += intToLetters(firstLetter)
		letter += string(rune('a' + colIndex%26))
	} else {
		letter += string(rune('a' + colIndex))
	}

	return
}
