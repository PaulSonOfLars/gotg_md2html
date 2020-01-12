package tg_md2html

func findLinkSections(in []rune) (int, int) {
	var textEnd, linkEnd int
	var offset int
	foundTextEnd := false
	for offset < len(in) {
		idx := stringIndex(in[offset:], "](")
		if idx < 0 {
			return -1, -1
		}
		textEnd = offset + idx
		if !IsEscaped(in, textEnd) {
			foundTextEnd = true
			break
		}
		offset = idx + 1
	}
	if !foundTextEnd {
		return -1, -1
	}

	offset = textEnd
	for offset < len(in) {
		idx := stringIndex(in[offset:], ")")
		if idx < 0 {
			return -1, -1
		}
		linkEnd = offset + idx
		if !IsEscaped(in, linkEnd) {
			return textEnd, linkEnd
		}
		offset = idx + 1
	}
	return -1, -1

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

func getTagClose(in []rune) int {
	for ix, c := range in {
		if c == '>' {
			return ix
		}
	}
	return -1
}

func getClosingTagOpen(in []rune) int {
	for ix, c := range in {
		if c == '<' && ix+1 < len(in) && in[ix+1] == '/' {
			return ix
		}
	}
	return -1
}

func getClosingTag(in []rune, tag string) (int, int) {
	offset := 0
	for offset < len(in) {
		o := getClosingTagOpen(in[offset:])
		if o < 0 {
			return -1, -1
		}
		open := offset + o
		c := getTagClose(in[open+2:])
		if c < 0 {
			return -1, -1
		}
		close := open + 2 + c
		if string(in[open+2:close]) == tag {
			return open, close
		}
		offset = open + 1
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
