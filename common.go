package tg_md2html

import (
	"regexp"
	"unicode"
)

func validStart(pos int, input []rune) bool {
	return (pos == 0 || !(unicode.IsLetter(input[pos-1]) || unicode.IsDigit(input[pos-1]))) && !(pos == len(input)-1 || unicode.IsSpace(input[pos+1]))
}

func validEnd(pos int, input []rune) bool {
	return !(pos == 0 || unicode.IsSpace(input[pos-1])) && (pos == len(input)-1 || !(unicode.IsLetter(input[pos+1]) || unicode.IsDigit(input[pos+1])))
}

func contains(r rune, rr []rune) bool {
	for _, x := range rr {
		if r == x {
			return true
		}
	}

	return false
}

// todo: remove regexp dep
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
