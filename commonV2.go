package tg_md2html

import (
	"strings"
)

// finds the middle '](' section of in a link markdown
func findLinkMidSectionIdx(in []rune, emoji bool) int {
	var textEnd int
	var offset int
	for offset < len(in) {
		idx := stringIndex(in[offset:], "](")
		if idx < 0 {
			return -1
		}
		textEnd = offset + idx
		if !IsEscaped(in, textEnd) {
			isEmoji := strings.HasPrefix(string(in[textEnd+2:]), "tg://emoji?id=")
			if (isEmoji && emoji) || (!isEmoji && !emoji) {
				return textEnd
			}
		}
		offset = textEnd + 1
	}
	return -1
}

// finds the closing ')' section of in a link markdown
func findLinkEndSectionIdx(in []rune) int {
	var linkEnd int
	var offset int
	for offset < len(in) {
		idx := getLinkEnd(in[offset:])
		if idx < 0 {
			return -1
		}
		linkEnd = offset + idx
		if !IsEscaped(in, linkEnd) {
			return linkEnd
		}
		offset = linkEnd + 1
	}
	return -1
}

// finds the middle and closing sections of in a link markdown
func findLinkSectionsIdx(in []rune, isEmojiLink bool) (int, int) {
	textEnd := findLinkMidSectionIdx(in, isEmojiLink)
	if textEnd < 0 {
		return -1, -1
	}

	linkEnd := findLinkEndSectionIdx(in[textEnd:])
	if linkEnd < 0 {
		return -1, -1
	}
	offsetLinkEnd := textEnd + linkEnd

	// We've found the first valid "mid" section above; and we've found the "end" section too.
	// Now, we iterate over the text in between the mid and end sections to see if any other mid sections exist.
	// If yes, we choose those instead - it would be invalid in a URL anyway.
	for textEnd < offsetLinkEnd {
		newTextEnd := findLinkMidSectionIdx(in[textEnd+1:offsetLinkEnd], isEmojiLink)
		if newTextEnd == -1 {
			break
		}
		textEnd = textEnd + newTextEnd + 1
	}

	return textEnd, offsetLinkEnd
}

func getLinkContents(in []rune, emoji bool) (bool, []rune, string, int) {
	// find ]( and then )
	textEndIdx, urlEndIdx := findLinkSectionsIdx(in, emoji)
	if textEndIdx < 0 || urlEndIdx < 0 {
		return false, nil, "", 0
	}

	content := string(in[textEndIdx+2 : urlEndIdx])
	text := in[1:textEndIdx]
	return true, text, content, urlEndIdx + 1
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

func getLinkEnd(in []rune) int {
	offset := 0
	for offset < len(in) {
		idx := stringIndex(in[offset:], ")")
		if idx < 0 {
			return -1
		}

		end := offset + idx
		// we don't check validEnd, since links can be inlined
		if !IsEscaped(in, end) {
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

func getClosingTag(in []rune, openingTag string, closingTag string) (int, int) {
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
		if string(in[openingTagIdx+1:closingTagIdx]) == openingTag { // found a nested tag, this is annoying
			subtags++
		} else if isClosingTag(in, openingTagIdx) && string(in[openingTagIdx+2:closingTagIdx]) == closingTag {
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
