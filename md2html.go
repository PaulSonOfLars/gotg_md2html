package tg_md2html// Copyright (c) Improbable Worlds Ltd, All Rights Reserved

import (
"strings"
)

var open = map[rune][]rune{
	'_': []rune("<i>"),
	'*': []rune("<b>"),
	'`': []rune("<code>"),
}

var close = map[rune][]rune{
	'_': []rune("</i>"),
	'*': []rune("</b>"),
	'`': []rune("</code>"),
}

// todo: ``` support? -> add \n char to md chars and hence on \n, skip
func Md2html(input []rune) ([]rune) {
	var output []rune
	v := map[rune][]int{}
	var containedMDChars []rune
	// todo: check why removing the escape characters doesnt require changing the offset
	escaped := false
	for pos, char := range input {
		if escaped {
			escaped = false
			continue
		}
		switch char {
		case '_', '*', '`', '[', ']', '(', ')':
			v[char] = append(v[char], pos)
			containedMDChars = append(containedMDChars, char)
		case '\\':
			escaped = true
			input = append(input[:pos], input[pos+1:]...)
		}
	}

	prev := 0
	for i := 0; i < len(containedMDChars) && prev < len(input); i++ {
		currChar := containedMDChars[i]
		switch currChar {
		case '_', '*', '`':
			posArr := v[currChar]
			if len(posArr) < 2 { // if less than two, skip
				continue
			}

			// skip i to next same char (hence jumping all inbetween) (could be done with a normal range and continues?)
			// todo: OOB check on +1?
			for _, val := range containedMDChars[i+1:] {
				i++
				if val == currChar {
					break
				}
				if len(v[val]) > 1 {
					v[val] = v[val][1:] // pop from map when skipped
				}
			}
			// pop currChar
			fstPos, sndPos, rest := posArr[0], posArr[1], posArr[2:]
			v[currChar] = rest
			output = append(output, input[prev:fstPos]...)
			output = append(output, open[currChar]...)
			output = append(output, input[fstPos+1:sndPos]...)
			output = append(output, close[currChar]...)
			prev = sndPos + 1

		case '[':
			openNameArr := v['[']
			nameOpen, rest := openNameArr[0], openNameArr[1:]
			v['['] = rest
			if len(v[']']) < 1 || len(v['(']) < 1 || len(v[')']) < 1 {
				continue
			}

			var nextNameClose int
			var nextLinkOpen int
			var nextLinkClose int

			wastedLinkClose := v[')']
		LinkLabel:
			for _, closeLinkPos := range v[')'] {
				if closeLinkPos > nameOpen {
					wastedLinkOpen := v['(']
					for _, openLinkpos := range v['('] {
						if openLinkpos > nameOpen && openLinkpos < closeLinkPos {
							wastedNameClose := v[']']
							for _, closeNamePos := range v[']'] {
								if closeNamePos == openLinkpos-1 {
									nextNameClose = closeNamePos
									nextLinkOpen = openLinkpos
									nextLinkClose = closeLinkPos
									v[']'] = wastedNameClose
									v['('] = wastedLinkOpen
									v[')'] = wastedLinkClose

									break LinkLabel
								}
								wastedNameClose = wastedNameClose[1:]
							}
						}
						wastedLinkOpen = wastedLinkOpen[1:]
					}
				}
				wastedLinkClose = wastedLinkClose[1:]
			}
			if nextLinkClose == 0 {
				// invalid
				continue
			}
			link := string(input[nextLinkOpen+1 : nextLinkClose])
			name := string(input[nameOpen+1 : nextNameClose])
			if strings.HasPrefix(link, "buttonurl:") {
				// is a button
				// todo: return correctly formatted values -> boolean to enable/disable? 2 exported funcs?
			} else {
				output = append(output, input[prev:nameOpen]...)
				output = append(output, []rune(`<a href="`+link+`">`+name+`</a>`)...)
			}

			prev = nextLinkClose + 1
		}

	}
	output = append(output, input[prev:]...)

	return output
}
