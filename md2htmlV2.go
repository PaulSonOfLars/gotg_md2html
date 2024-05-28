package tg_md2html

import (
	"html"
	"sort"
	"strings"
	"unicode"
)

var defaultConverterV2 = ConverterV2{
	Prefixes: map[string]string{
		"url": "buttonurl:",
	},
	SameLineSuffix: sameLineSuffix,
}

// ButtonV2 identifies a button. It can contain either a URL, or Text, depending on whether it is a buttonURL: or a buttonText:
type ButtonV2 struct {
	Name     string
	Type     string
	Content  string
	SameLine bool
}

type ConverterV2 struct {
	Prefixes       map[string]string
	SameLineSuffix string
}

func NewV2(prefixes map[string]string) *ConverterV2 {
	return &ConverterV2{
		Prefixes:       prefixes,
		SameLineSuffix: sameLineSuffix,
	}
}

func MD2HTMLV2(in string) string {
	return defaultConverterV2.MD2HTML(in)
}

func MD2HTMLButtonsV2(in string) (string, []ButtonV2) {
	return defaultConverterV2.MD2HTMLButtons(in)
}

var chars = map[string]string{
	"`":      "code",
	"```":    "pre",
	"_":      "i",
	"*":      "b",
	"~":      "s",
	"__":     "u",
	"|":      "", // this is a placeholder for || to work
	"||":     "span class=\"tg-spoiler\"",
	"!":      "", // for emoji
	"[":      "", // for links
	"]":      "", // for links/emoji
	"(":      "", // for links/emoji
	")":      "", // for links/emoji
	"\\":     "", // for escapes
	"&":      "", // for blockquotes
	"&gt;":   "blockquote",
	"**&gt;": "blockquote", // expandable blockquotes
}

var AllMarkdownV2Chars = func() []rune {
	var outString []string
	for k := range chars {
		outString = append(outString, k)
	}
	sort.Strings(outString)
	var out []rune
	for _, x := range outString {
		out = append(out, []rune(x)[0])
	}
	return out
}()

func (cv ConverterV2) MD2HTML(in string) string {
	text, _ := cv.md2html([]rune(html.EscapeString(in)), false)
	return strings.TrimSpace(text)
}

func (cv ConverterV2) MD2HTMLButtons(in string) (string, []ButtonV2) {
	text, btns := cv.md2html([]rune(html.EscapeString(in)), true)
	return strings.TrimSpace(text), btns
}

var skipStarts = map[rune]bool{
	'!': true, // premium emoji
	'[': true, // links
}

func getItem(in []rune, i int) (string, int, bool) {
	c := in[i]
	if _, ok := chars[string(c)]; !ok {
		return "", 0, false
	}

	if !validStart(i, in) && !skipStarts[c] {
		if c == '\\' && i+1 < len(in) {
			escaped := string(in[i+1])
			if _, ok := chars[escaped]; ok {
				return escaped, 1, false
			}
		}
		return "", 0, false
	}

	item := string(c)
	if c == '|' { // support ||
		if !(i+1 < len(in) && in[i+1] == '|') {
			return "", 0, false
		}
		return "||", 1, true

	} else if c == '_' && i+1 < len(in) && in[i+1] == '_' { // support __
		return "__", 1, true

	} else if c == '`' && i+2 < len(in) && in[i+1] == '`' && in[i+2] == '`' { // support ```
		return "```", 2, true

	} else if c == '&' {
		if !(i+3 < len(in) && in[i+1] == 'g' && in[i+2] == 't' && in[i+3] == ';') {
			return "", 0, false
		}
		if !validBlockQuoteStart(in, i) {
			return "", 0, false
		}
		return "&gt;", 3, true

	} else if c == '*' && i+5 < len(in) && in[i+1] == '*' && in[i+2] == '&' && in[i+3] == 'g' && in[i+4] == 't' && in[i+5] == ';' {
		if !validBlockQuoteStart(in, i) {
			return "", 0, false
		}
		return "**&gt;", 5, true
	}

	return item, 0, true
}

// TODO: add support for a map-like check of which items cannot be included.
//
//	Eg: `code` cannot be italic/bold/underline/strikethrough
//	however... this is currently implemented by server side by telegram, so not my problem :runs:
//
// (see notes on: https://core.telegram.org/bots/api#markdownv2-style)
func (cv ConverterV2) md2html(in []rune, enableButtons bool) (string, []ButtonV2) {
	out := strings.Builder{}

	for i := 0; i < len(in); i++ {
		// TODO: Add blockquote handling here too
		item, offset, ok := getItem(in, i)
		if !ok {
			if item == "" {
				item = string(in[i])
			}
			out.WriteString(item)
			i += offset
			continue
		}
		i += offset

		switch item {
		// All cases where start and closing tags are the same.
		case "`", "*", "~", "_", "|", "```", "||", "__":
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

			nStart, nEnd := i+1, i+idx+1
			followT, followB := cv.md2html(in[nEnd+len(item):], enableButtons)

			if item == "`" {
				// ` doesn't support nested items, so don't parse children.
				return out.String() + "<code>" + string(in[nStart:nEnd]) + "</code>" + followT, followB

			} else if item == "```" {
				// ``` doesn't support nested items, so don't parse children.
				nestedT := string(in[nStart:nEnd])

				// Attempt to extract language details; should only be first line
				splitLines := strings.Split(nestedT, "\n")
				if len(splitLines) > 1 {
					// TODO: How do we decide the language; first word? first line?
					firstLine := strings.TrimSpace(splitLines[0])
					if len(firstLine) > 0 && strings.HasPrefix(nestedT, firstLine) {
						content := strings.TrimPrefix(nestedT, firstLine+"\n")
						return out.String() + "<pre><code class=\"language-" + firstLine + "\">" + content + "</code></pre>" + followT, followB
					}
				}
				return out.String() + "<pre>" + strings.TrimPrefix(nestedT, "\n") + "</pre>" + followT, followB
			}

			// internal won't have any interesting item closings
			nestedT, nestedB := cv.md2html(in[nStart:nEnd], enableButtons)
			return out.String() + "<" + chars[item] + ">" + nestedT + "</" + closeSpans(chars[item]) + ">" + followT, append(nestedB, followB...)

		case "&gt;", "**&gt;":
			nStart := i + 1
			for nStart < len(in) && unicode.IsSpace(in[nStart]) {
				nStart++
			}

			if nStart >= len(in) {
				out.WriteString(item)
				continue
			}

			nEnd, contents := getBlockQuoteEnd(in, nStart, item == "**&gt;")
			nestedT, nestedB := cv.md2html(contents, enableButtons)
			followT, followB := cv.md2html(in[nEnd:], enableButtons)

			if item == "**&gt;" {
				return out.String() + "<blockquote expandable>" + strings.TrimSpace(nestedT) + "</blockquote>" + followT, append(nestedB, followB...)
			}
			return out.String() + "<blockquote>" + strings.TrimSpace(nestedT) + "</blockquote>" + followT, append(nestedB, followB...)

		case "!":
			// TODO: Handle this in item selection
			if len(in) <= i+1 || in[i+1] != '[' {
				out.WriteString(item)
				continue
			}

			ok, text, content, newEnd := getLinkContents(in[i+1:], true)
			if !ok {
				out.WriteString(item)
				continue
			}
			end := i + 1 + newEnd

			content = strings.TrimPrefix(content, "tg://emoji?id=")

			nestedT, nestedB := cv.md2html(text, enableButtons)
			followT, followB := cv.md2html(in[end:], enableButtons)
			return out.String() + `<tg-emoji emoji-id="` + content + `">` + nestedT + "</tg-emoji>" + followT, append(nestedB, followB...)

		case "[":
			ok, text, content, newEnd := getLinkContents(in[i:], false)
			if !ok {
				out.WriteString(item)
				continue
			}
			end := i + newEnd

			followT, followB := cv.md2html(in[end:], enableButtons)

			if enableButtons {
				for buttonType, prefix := range cv.Prefixes {
					if !strings.HasPrefix(content, prefix) {
						continue
					}

					content := strings.TrimLeft(strings.TrimPrefix(content, prefix), "/")
					sameline := strings.HasSuffix(content, cv.SameLineSuffix)
					if sameline {
						content = strings.TrimSuffix(content, cv.SameLineSuffix)
					}
					cleanedName := cv.StripMDV2(string(text))
					return out.String() + followT, append([]ButtonV2{{
						Name:     html.UnescapeString(cleanedName),
						Type:     buttonType,
						Content:  content,
						SameLine: sameline,
					}}, followB...)
				}
			}

			nestedT, nestedB := cv.md2html(text, enableButtons)
			return out.String() + `<a href="` + content + `">` + nestedT + "</a>" + followT, append(nestedB, followB...)

		case "]", "(", ")":
			out.WriteString(item)

		case "\\":
			if i+1 < len(in) {
				if _, ok := chars[string(in[i+1])]; ok {
					out.WriteRune(in[i+1])
					i++
					continue
				}
			}
			out.WriteString(item)
		}
	}

	return out.String(), nil
}

func getBlockQuoteEnd(in []rune, nStart int, expandable bool) (int, []rune) {
	var contents []rune // We store all the contents, minus the > characters, so we avoid double-html tags
	lineStart := true
	for j := nStart; j < len(in); j++ {
		if lineStart && in[j] == ' ' {
			// Skip space chars at start of lines
			continue
		}

		lineStart = in[j] == '\n'
		contents = append(contents, in[j])

		// Keep skipping until we get a newline
		if in[j] != '\n' {
			continue
		}

		if isExpandableEnd(in, expandable, j) {
			return j, contents[:len(contents)-3]
		}

		if j+4 < len(in) && in[j+1] == '&' && in[j+2] == 'g' && in[j+3] == 't' && in[j+4] == ';' {
			j = j + 4 // skip '>' symbol for the next blockquote start
			continue
		}
		return j, contents
	}

	if isExpandableEnd(in, expandable, len(in)) {
		return len(in), contents[:len(contents)-2]
	}

	return len(in), contents
}

func isExpandableEnd(in []rune, expandable bool, j int) bool {
	return expandable && j-2 >= 0 && in[j-1] == '|' && in[j-2] == '|'
}

func validBlockQuoteStart(in []rune, i int) bool {
	for j := i - 1; j >= 0; j-- {
		if !unicode.IsSpace(in[j]) {
			return false
		}
		if in[j] == '\n' {
			return true
		}
	}

	// Start of message; must be valid.
	return true
}

func EscapeMarkdownV2(r []rune) string {
	out := strings.Builder{}
	for i, x := range r {
		if contains(x, AllMarkdownV2Chars) {
			if i == 0 || i == len(r)-1 || validEnd(i, r) || validStart(i, r) {
				out.WriteRune('\\')
			}
		}
		out.WriteRune(x)
	}
	return out.String()
}

// closeSpans gets the correct closing tags for spans.
// eg:
// - closeSpans("span class=\"tg-spoiler\"") should return just "span"
// - closeSpans("pre") -> returns "pre"
func closeSpans(s string) string {
	if !strings.HasPrefix(s, "span") {
		return s
	}

	return "span"
}
