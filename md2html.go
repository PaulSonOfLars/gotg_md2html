package tg_md2html

import (
	"strings"
	"unicode"
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

func MD2HTML(input string) string {
	text, _, _ := md2html([]rune(input), false)
	return text
}

func MD2HTMLButtons(input string) (string, []string, []string) {
	return md2html([]rune(input), true)
}

// todo: ``` support? -> add \n char to md chars and hence on \n, skip
func md2html(input []rune, buttons bool) (string, []string, []string) {
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
	var btnNames []string
	var btnLinks []string
	for i := 0; i < len(containedMDChars) && prev < len(input); i++ {
		currChar := containedMDChars[i]
		switch currChar {
		case '_', '*', '`':
			posArr := v[currChar]
			// if fewer than 2 chars left, pass
			if len(posArr) < 2 {
				continue
			}
			// if we're past the currChar position, pass and update
			if posArr[0] < prev {
				v[currChar] = posArr[1:]
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
			fstPos, rest := posArr[0], posArr[1:]
			v[currChar] = rest

			if !((fstPos == 0 || !(unicode.IsLetter(input[fstPos-1]) || unicode.IsDigit(input[fstPos-1]))) && !(fstPos == len(input)-1 || unicode.IsSpace(input[fstPos+1]))) {
				continue
			}
			ok := false
			var sndPos int
			for _, sndPos = range rest {
				rest = rest[1:]
				if !(sndPos == 0 || unicode.IsSpace(input[sndPos-1])) && (sndPos == len(input)-1 || !(unicode.IsLetter(input[sndPos+1]) || unicode.IsDigit(input[sndPos+1]))) {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}

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
			if len(v[']']) < 1 || len(v['(']) < 1 || len(v[')']) < 1 || nameOpen < prev {
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

			output = append(output, input[prev:nameOpen]...)
			link := string(input[nextLinkOpen+1 : nextLinkClose])
			name := string(input[nameOpen+1 : nextNameClose])
			if buttons && strings.HasPrefix(link, "buttonurl:") {
				// is a button
				btnNames = append(btnNames, name)
				btnLinks = append(btnLinks, link)
			} else {
				output = append(output, []rune(`<a href="`+link+`">`+name+`</a>`)...)
			}

			prev = nextLinkClose + 1
		}
	}
	output = append(output, input[prev:]...)

	return string(output), btnNames, btnLinks
}
