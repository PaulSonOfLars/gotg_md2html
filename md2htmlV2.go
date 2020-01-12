package tg_md2html

import (
	"fmt"
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
	"`":   "code",
	"```": "pre",
	"_":   "i",
	"*":   "b",
	"~":   "s",
	"__":  "u",
	"[":   "", // for links
	"]":   "", // for links
	"(":   "", // for links
	")":   "", // for links
	"\\":  "", // for escapes
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
		case '`', '*', '~', '_': // '__' and '```' are included here too
			item := string(c)
			if c == '_' && i+1 < len(in) && in[i+1] == '_' { // support __
				item = "__"
				i++
			} else if c == '`' && i+2 < len(in) && in[i+1] == '`' && in[i+2] == '`' { // support ```
				item = "```"
				i += 2
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

func (v ConverterV2) Reverse(in string, bs []ButtonV2) (string, error) {
	return v.reverse([]rune(in), bs)
}

func (v ConverterV2) reverse(in []rune, buttons []ButtonV2) (string, error) {
	prev := 0
	out := strings.Builder{}
	for i := 0; i < len(in); i++ {
		switch in[i] {
		case '<':
			c := getTagClose(in[i+1:])
			if c < 0 {
				// "no close tag"
				return "", fmt.Errorf("no closing '>' for opening bracket at %d", i)
			}
			closeTag := i + c + 1
			tagContent := string(in[i+1 : closeTag])
			tagFields := strings.Fields(tagContent)
			if len(tagFields) < 1 {
				return "", fmt.Errorf("no tag name for HTML tag started at %d", i)
			}
			tag := tagFields[0]

			co, cc := getClosingTag(in[closeTag+1:], tag)
			if co < 0 || cc < 0 {
				// "no closing open"
				return "", fmt.Errorf("no closing tag for HTML tag %s started at %d", tag, i)
			}
			closingOpen, closingClose := closeTag+1+co, closeTag+1+cc
			out.WriteString(html.UnescapeString(string(in[prev:i])))

			nested, err := ReverseV2(string(in[closeTag+1:closingOpen]), nil)
			if err != nil {
				return "", err
			}

			switch tag {
			case "b", "strong":
				out.WriteString("*" + nested + "*")
			case "i", "em":
				out.WriteString("_" + nested + "_")
			case "u", "ins":
				out.WriteString("__" + nested + "_")
			case "s", "strike", "del":
				out.WriteString("~" + nested + "~")
			case "code":
				out.WriteString("`" + nested + "`")
			case "pre":
				out.WriteString("```" + nested + "```")
			case "a":
				if link.MatchString(tagContent) {
					matches := link.FindStringSubmatch(tagContent)
					out.WriteString("[" + nested + "](" + matches[1] + ")")
				} else {
					return "", fmt.Errorf("badly formatted anchor tag %q", tagContent)
				}
			default:
				return "", fmt.Errorf("unknown tag %q", tag)
			}

			prev = closingClose + 1
			i = closingClose

		case '\\', '_', '*', '~', '`', '[', ']', '(', ')': // these all need to be escaped to ensure we retain the same message
			out.WriteString(html.UnescapeString(string(in[prev:i])))
			out.WriteRune('\\')
			out.WriteRune(in[i])
			prev = i + 1
		}
	}
	out.WriteString(html.UnescapeString(string(in[prev:])))

	for _, btn := range buttons {
		out.WriteString("\n[" + btn.Name + "](" + v.BtnPrefix + "://" + html.UnescapeString(btn.Content))
		if btn.SameLine {
			out.WriteString(v.SameLineSuffix)
		}
		out.WriteString(")")
	}

	return out.String(), nil
}

func ReverseV2(in string, bs []ButtonV2) (string, error) {
	return defaultConverterV2.Reverse(in, bs)
}

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

		end := offset + idx + len(s) - 1 // to account for __
		if validEnd(end, in) && !IsEscaped(in, end) {
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
