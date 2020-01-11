package tg_md2html

import (
	"html"
	"strings"
)

var defaultConverterV2 = ConverterV2{
	BtnPrefix:      btnPrefix,
	SameLineSuffix: sameLineSuffix,
}

type ButtonV2 struct {
	Name     string
	Content  string
	SameLine bool
}

type ConverterV2 struct {
	BtnPrefix      string
	SameLineSuffix string
}

func NewV2() *ConverterV2 {
	return &ConverterV2{
		BtnPrefix:      btnPrefix,
		SameLineSuffix: sameLineSuffix,
	}
}

func MD2HTMLV2(in string) string {
	return defaultConverterV2.MD2HTML(in)
}

func MD2HTMLButtonsV2(in string) (string, []ButtonV2) {
	return defaultConverterV2.md2html([]rune(html.EscapeString(in)), true)
}

var chars = map[string]string{
	"`":  "code",
	"_":  "i",
	"*":  "b",
	"~":  "s",
	"__": "u",
	"[":  "", // for links
	"]":  "", // for links
	"(":  "", // for links
	")":  "", // for links
	"\\": "", // for escapes
}

func (v ConverterV2) MD2HTML(in string) string {
	text, _ := v.md2html([]rune(html.EscapeString(in)), false)
	return text
}

func (v ConverterV2) MD2HTMLButtons(in string) (string, []ButtonV2) {
	return v.md2html([]rune(html.EscapeString(in)), true)
}

func (v ConverterV2) md2html(in []rune, b bool) (string, []ButtonV2) {
	out := strings.Builder{}

	for i := 0; i < len(in); i++ {
		c := in[i]
		if _, ok := chars[string(c)]; !ok {
			out.WriteRune(c)
			continue
		}

		if !validStart(i, in) {
			if c == '\\' && i+1 < len(in) {
				if _, ok := chars[string(in[i+1])]; ok {
					out.WriteRune(in[i+1])
					i++
					continue
				}
			}
			out.WriteRune(c)
			continue
		}

		switch c {
		case '`', '*', '~', '_': // '__' is included here too
			item := string(c)
			if c == '_' && i+1 < len(in) && in[i+1] == '_' {
				item = "__"
				i++
			}

			if i+1 >= len(in) {
				out.WriteString(item)
				continue
			}
			idx := getValidEnd(in[i+1:], item)
			if idx < 0 {
				// not found; write and move on.
				out.WriteString(item)
				continue
			}
			nStart, nEnd := i+1, i+idx+2-len(item) // +2 because start is at +1 already

			// internal is guaranteed not to have any valid item closings, since we've greedily taken them.
			nestedT, nestedB := v.md2html(in[nStart:nEnd], b)
			followT, followB := v.md2html(in[nEnd+len(item):], b) // offset?
			return out.String() + "<" + chars[item] + ">" + nestedT + "</" + chars[item] + ">" + followT, append(nestedB, followB...)

		case '[':
			// find ]( and then )
			linkText, linkURL := findLinkSections(in[i:])
			if linkText < 0 || linkURL < 0 {
				out.WriteRune(c)
				continue
			}

			content := string(in[i+linkText+2 : i+linkURL])
			text := in[i+1 : i+linkText]
			end := i + linkURL + 1
			followT, followB := v.md2html(in[end:], b)

			if b && strings.HasPrefix(content, v.BtnPrefix) {
				content = strings.TrimLeft(content[len(v.BtnPrefix):], "/")
				sameline := false
				if strings.HasSuffix(content, v.SameLineSuffix) {
					sameline = true
					content = content[:len(content)-len(v.SameLineSuffix)]
				}
				return out.String() + followT, append([]ButtonV2{{
					Name:     string(text),
					Content:  content,
					SameLine: sameline,
				}}, followB...)
			}
			nestedT, nestedB := v.md2html(text, true)
			return out.String() + `<a href="` + content + `">` + nestedT + "</a>" + followT, append(nestedB, followB...)

		case ']', '(', ')':
			out.WriteRune(c)

		case '\\':
			if i+1 < len(in) {
				if _, ok := chars[string(in[i+1])]; ok {
					out.WriteRune(in[i+1])
					i++
					continue
				}
			}
			out.WriteRune(c)
		}
	}

	return out.String(), nil
}

func findLinkSections(in []rune) (int, int) {
	var textEnd, linkEnd int
	var offset int
	foundTextEnd := false
	for offset < len(in) {
		idx := stringIndex(string(in[offset:]), "](")
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
		idx := stringIndex(string(in[offset:]), ")")
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

func getNext(in []rune, s string) int {
	offset := 0
	for offset < len(in) {
		idx := stringIndex(string(in[offset:]), s)
		if idx < 0 {
			return -1
		}

		end := offset + idx + len(s) - 1 // to account for __
		if !IsEscaped(in, end) {
			return end
		}
		offset = end + 1
	}
	return -1
}

func getValidEnd(in []rune, s string) int {
	offset := 0
	for offset < len(in) {
		idx := stringIndex(string(in[offset:]), s)
		if idx < 0 {
			return -1
		}

		end := offset + idx + len(s) - 1 // to account for __
		if validEnd(end, in) && !IsEscaped(in, end) {
			return end
		}
		offset = end + 1
	}
	return -1
}

func runeIndex(in string, r rune) int {
	i := strings.IndexRune(in, r)
	if i < 0 {
		return i
	}
	return len([]rune(in[:i]))
}

func stringIndex(in string, s string) int {
	i := strings.Index(in, s)
	if i < 0 {
		return i
	}
	return len([]rune(in[:i]))
}
