package tg_md2html

func findLinkSectionsIdx(in []rune) (int, int) {
	var textEnd, linkEnd int
	var offset int
	for offset < len(in) {
		idx := stringIndex(in[offset:], "](")
		if idx < 0 {
			return -1, -1
		}
		textEnd = offset + idx
		if !IsEscaped(in, textEnd) {
			break
		}
		offset = textEnd + 1
	}
	if offset >= len(in) {
		return -1, -1
	}

	offset = textEnd
	for offset < len(in) {
		idx := getValidLinkEnd(in[offset:])
		if idx < 0 {
			return -1, -1
		}
		linkEnd = offset + idx
		if !IsEscaped(in, linkEnd) {
			return textEnd, linkEnd
		}
		offset = linkEnd + 1
	}
	return -1, -1
}

func getLinkContents(in []rune) (bool, []rune, string, int) {
	// find ]( and then )
	linkText, linkURL := findLinkSectionsIdx(in)
	if linkText < 0 || linkURL < 0 {
		return false, nil, "", 0
	}

	content := string(in[linkText+2 : linkURL])
	text := in[1:linkText]
	return true, text, content, linkURL + 1
}

func getValidEnd(in []rune, s string) int {
	offset := 0
	for offset < len(in) {
		idx := stringIndex(in[offset:], s)
		if idx < 0 {
			return -1
		}

		end := offset + idx
		// validEnd check has double logic to account for multi char strings
		if validEnd(end, in) && validEnd(end+len(s)-1, in) && !IsEscaped(in, end) {
			idx = stringIndex(in[end+1:], s)
			for idx == 0 {
				end++
				idx = stringIndex(in[end+1:], s)
			}
			return end
		}
		offset = end + 1
	}
	return -1
}

func getValidLinkEnd(in []rune) int {
	offset := 0
	for offset < len(in) {
		idx := stringIndex(in[offset:], ")")
		if idx < 0 {
			return -1
		}

		end := offset + idx
		// validEnd check has double logic to account for multi char strings
		if validEnd(end, in) && !IsEscaped(in, end) {
			return end
		}
		offset = end + 1
	}
	return -1
}

func getHTMLTagOpenIndex(in []rune) int {
	for idx, c := range in {
		if c == '<' {
			return idx
		}
	}
	return -1
}

func getHTMLTagCloseIndex(in []rune) int {
	for idx, c := range in {
		if c == '>' {
			return idx
		}
	}
	return -1
}

func isClosingTag(in []rune, pos int) bool {
	if in[pos] == '<' && pos+1 < len(in) && in[pos+1] == '/' {
		return true
	}
	return false
}

func getClosingTag(in []rune, tag string) (int, int) {
	offset := 0
	subtags := 0
	for offset < len(in) {
		o := getHTMLTagOpenIndex(in[offset:])
		if o < 0 {
			return -1, -1
		}
		openingTagIdx := offset + o

		c := getHTMLTagCloseIndex(in[openingTagIdx+2:])
		if c < 0 {
			return -1, -1
		}

		closingTagIdx := openingTagIdx + 2 + c
		if string(in[openingTagIdx+1:closingTagIdx]) == tag { // found a nested tag, this is annoying
			subtags++
		} else if isClosingTag(in, openingTagIdx) && string(in[openingTagIdx+2:closingTagIdx]) == tag {
			if subtags == 0 {
				return openingTagIdx, closingTagIdx
			}
			subtags--
		}
		offset = openingTagIdx + 1
	}
	return -1, -1
}

func stringIndex(in []rune, s string) int {
	r := []rune(s)
	for idx := range in {
		if startsWith(in[idx:], r) {
			return idx
		}
	}
	return -1
}

func startsWith(i []rune, p []rune) bool {
	for idx, x := range p {
		if idx >= len(i) || i[idx] != x {
			return false
		}
	}
	return true
}
