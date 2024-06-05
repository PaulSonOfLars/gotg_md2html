package tg_md2html

import (
	"regexp"
	"unicode"
)

func validStart(pos int, input []rune) bool {
	// Last char is not a valid start char.
	// If the next char is a space, it isn't a valid start either.
	if pos == len(input)-1 || unicode.IsSpace(input[pos+1]) {
		return false
	}

	// First char is always a valid start.
	if pos == 0 {
		return true
	}

	// If the previous char is alphanumeric, it is an invalid start char.
	return !(unicode.IsLetter(input[pos-1]) || unicode.IsDigit(input[pos-1]))
}

func validEnd(pos int, input []rune) bool {
	// First char is not a valid end char; we do NOT allow empty entities.
	// If the end char has a space before it, its not valid either.
	if pos == 0 || unicode.IsSpace(input[pos-1]) {
		return false
	}

	// Last char is always a valid end char;
	if pos == len(input)-1 {
		return true
	}

	// If the next char is alphanumeric, it is an invalid end char.
	return !(unicode.IsLetter(input[pos+1]) || unicode.IsDigit(input[pos+1]))
}

func contains(r rune, rr []rune) bool {
	for _, x := range rr {
		if r == x {
			return true
		}
	}

	return false
}

var link = regexp.MustCompile(`a href="(.*)"`)
var customEmoji = regexp.MustCompile(`tg-emoji emoji-id="(.*)"`)

func IsEscaped(input []rune, pos int) bool {
	if pos == 0 {
		return false
	}

	i := pos - 1
	for ; i >= 0; i-- {
		if input[i] == '\\' {
			continue
		}
		break
	}

	return (pos-i)%2 == 0
}
